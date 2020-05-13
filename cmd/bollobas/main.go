package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/taxibeat/bollobas/internal/config"
	"github.com/taxibeat/bollobas/internal/ingestion"
	"github.com/taxibeat/bollobas/internal/ingestion/driver"
	"github.com/taxibeat/bollobas/internal/ingestion/passenger"
	"github.com/taxibeat/bollobas/internal/ingestion/ride"
	"github.com/taxibeat/bollobas/internal/mixpanel"
	"github.com/taxibeat/bollobas/internal/mixpanel/configuration"
	"github.com/taxibeat/bollobas/internal/mixpanel/identity"
	"github.com/taxibeat/bollobas/internal/mixpanel/riderequest"
	"github.com/taxibeat/bollobas/internal/mixpanel/riderequest/cancellation"
	"github.com/taxibeat/bollobas/internal/mixpanel/riderequest/confirmation"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"os"

	"context"

	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"

	phttp "github.com/beatlabs/patron/component/http"

	"time"

	semhttp "github.com/taxibeat/bollobas/internal/semovi/rest/http"

	mpsdk "github.com/dukex/mixpanel"
	"github.com/joho/godotenv"
)

var (
	version = "dev"
)

func main() {

	err := godotenv.Load("../../config/.env")
	if err != nil {
		log.Debugf("no .env file exists: %v", err)
	}

	ctx := context.Background()
	name := "bollobas"

	cfg := &config.Configuration{}
	// Setupbollobas config
	h, err := config.NewConfig(cfg)
	if err != nil {
		fmt.Printf("failed to set up configuration: %v", err)
		os.Exit(1)
	}

	err = h.Harvest(ctx)
	if err != nil {
		fmt.Printf("failed harvesting configuration: %v", err)
		os.Exit(1)
	}

	err = patron.SetupLogging(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	store, err := sql.New(cfg)
	if err != nil {
		fmt.Printf("failed to connect to the database: %v", err)
		os.Exit(1)
	}
	defer store.Close()

	// Setup HTTP route builder with a singe GET route
	rrb := phttp.NewRoutesBuilder()
	appendRoutes(rrb, semhttp.Routes(ctx, store))

	//Set up Kafka
	rt := uint(10)
	rtw := time.Duration(5 * time.Second)

	durl := "inproc://driver-publisher"
	drAccProc, err := driver.NewAccountProcessor(durl, cfg.Location.Get())
	drKfkCmp, err := ingestion.NewKafkaComponent(
		"driver-identity",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KafkaDriverIdentityTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		drAccProc,
		rt,
		rtw,
	)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	purl := "inproc://passenger-publisher"
	paAccProc, err := passenger.NewAccountProcessor(
		purl,
		"kafka",
		cfg.KafkaPassengerIdentityTopic.Get(),
		cfg.Location.Get(),
	)
	paAccProc.Activate(true)
	paKfkCmp, err := ingestion.NewKafkaComponent(
		"passenger-identity",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KafkaPassengerIdentityTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		paAccProc,
		rt,
		rtw,
	)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	prurl := "inproc://riderequest-publisher"
	paRRProc, err := passenger.NewRequestProcessor(prurl, "kafka", cfg.KkPRRTopic.Get())
	paRRProc.Activate(true)
	paRRKfkCmp, err := ingestion.NewKafkaComponent(
		"passenger-request",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KkPRRTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		paRRProc,
		rt,
		rtw,
	)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	pclurl := "inproc://riderequestcancel-publisher"
	paRCProc, err := passenger.NewCancellationProcessor(pclurl, "kafka", cfg.KkPRCTopic.Get())
	paRCProc.Activate(true)
	paRCKfkCmp, err := ingestion.NewKafkaComponent(
		"passenger-cancel",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KkPRCTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		paRCProc,
		rt,
		rtw,
	)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	pakurl := "inproc://ride-publisher"
	paRAKProc, err := ride.NewRideProcessor(pakurl, "kafka", cfg.KkRTopic.Get())
	paRAKProc.Activate(true)
	paRAKKfkCmp, err := ingestion.NewKafkaComponent(
		"ride",
		cfg.KkRTopic.Get(),
		[]string{cfg.KkPRCTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		paRAKProc,
		rt,
		rtw,
	)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	sig := func() {
		fmt.Println("exit gracefully...")
		os.Exit(0)
	}

	// Append components
	ccmp := make([]patron.Component, 0)
	ccmp = append(ccmp, drKfkCmp, paKfkCmp, paRRKfkCmp, paRCKfkCmp, paRAKKfkCmp)

	defaultConf := map[string]interface{}{}
	err = json.Unmarshal([]byte(cfg.BConf.Get()), &defaultConf)
	if err != nil {
		log.Fatalf("Wrong configuration provided %v", err)
	}

	// MIXPANEL handlers
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}

	//Conf Manager MixPanel
	cfm := mixpanel.NewConfigurator()
	cfm.Configure(defaultConf)

	plCfg := configuration.RestPoller{
		Manager:       cfm,
		RestURL:       cfg.RestURL.Get(),
		PollingPeriod: cfg.SettingsPeriod.Get(),
		DefaultConf:   defaultConf,
		RestKey:       cfg.RestKey.Get(),
		Path:          cfg.RestMixpanelPath.Get(),
	}

	plCfg.UpdateSettings()

	mpCl := mpsdk.NewFromClient(c, cfg.MpToken.Get(), "")
	// Handler for any identity change from Passegers and Drivers
	ipr := &identity.Processor{Mixpanel: mpCl}
	mph := mixpanel.NewHandler(ipr, []string{purl, durl}, cfm)
	mph.Run()

	// Handler for ride requests
	rrp := &riderequest.Processor{Mixpanel: mpCl}
	rrh := mixpanel.NewHandler(rrp, []string{prurl}, cfm)
	rrh.Run()

	// Handler for Request cancellations
	rrcp := &cancellation.Processor{Mixpanel: mpCl}
	rrch := mixpanel.NewHandler(rrcp, []string{pclurl}, cfm)
	rrch.Run()

	// Handler for Confirmations of rides
	rrakp := &confirmation.Processor{Mixpanel: mpCl}
	rrakh := mixpanel.NewHandler(rrakp, []string{pakurl}, cfm)
	rrakh.Run()

	// Initiate patron
	err = patron.New(name, version).
		WithSIGHUP(sig).
		WithRoutesBuilder(rrb).
		WithComponents(ccmp...).
		Run(ctx)

	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

}

func appendRoutes(rrb *phttp.RoutesBuilder, rb []*phttp.RouteBuilder) {
	for _, v := range rb {
		rrb.Append(v)
	}
}

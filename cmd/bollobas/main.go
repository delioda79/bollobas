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
	"github.com/taxibeat/bollobas/internal/ingestion/semovi"
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

const (
	version = "dev"
	name    = "bollobas"
)

// @Title Bollobas
// @Version 1.0.0
// @Description Bollobas microservice is responsible for any analytics that go through Beat's backend platform.
// @tag.name bollobas
// @contact.name RDXP3
// @contact.url https://confluence.taxibeat.com/display/TEAM404/RDXP3+-+TechnoMules
// @schemes http
// @license.name BEAT Mobility Services
func main() {

	err := godotenv.Load("../../config/.env")
	if err != nil {
		log.Debugf("no .env file exists: %v", err)
	}

	ctx := context.Background()

	// Setup bollobas config
	cfg := &config.Configuration{}
	h, err := config.NewConfig(cfg)
	if err != nil {
		fmt.Printf("failed to set up configuration: %v", err)
		os.Exit(1)
	}
	if err = h.Harvest(ctx); err != nil {
		fmt.Printf("failed harvesting configuration: %v", err)
		os.Exit(1)
	}
	defaultConf := map[string]interface{}{}
	if err = json.Unmarshal([]byte(cfg.BConf.Get()), &defaultConf); err != nil {
		log.Fatalf("Wrong configuration provided %v", err)
	}

	// Setup logging
	if err = patron.SetupLogging(name, version); err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	// Setup SQL
	store, err := sql.New(cfg)
	if err != nil {
		fmt.Printf("failed to connect to the database: %v", err)
		os.Exit(1)
	}
	defer store.Close()

	// Initialize HTTP routes
	rrb := phttp.NewRoutesBuilder()

	// Initialize components
	ccmp := make([]patron.Component, 0)

	// Setup mixpanel components
	if err := setupMixpanelComponents(cfg, defaultConf, &ccmp); err != nil {
		log.Fatalf("Failed to set up processor %v", err)
	}

	// Setup semovi routes and components
	if err = setupSemoviComponents(cfg, store, rrb, &ccmp); err != nil {
		log.Fatalf("Failed to set up processor %v", err)
	}

	sig := func() {
		fmt.Println("exit gracefully...")
		os.Exit(0)
	}

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

func setupMixpanelComponents(cfg *config.Configuration, defaultConf map[string]interface{}, ccmp *[]patron.Component) (err error) {
	// MIXPANEL handlers
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}

	//Conf Manager MixPanel
	cfm := mixpanel.NewConfigurator()
	cfm.Configure(defaultConf)

	rt := uint(10)
	rtw := time.Duration(5 * time.Second)
	mpCl := mpsdk.NewFromClient(c, cfg.MpToken.Get(), "")

	durl := "inproc://driver-publisher"
	drAccProc, err := driver.NewAccountProcessor(durl, cfg.Location.Get())
	if err != nil {
		return err
	}
	drKfkCmp, err := ingestion.NewKafkaComponent(
		"driver-identity",
		"driver-kafka-cmp",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KafkaDriverIdentityTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		drAccProc,
		rt,
		rtw,
	)
	if err != nil {
		return err
	}

	purl := "inproc://passenger-publisher"
	paAccProc, err := passenger.NewAccountProcessor(
		purl,
		"kafka",
		cfg.KafkaPassengerIdentityTopic.Get(),
		cfg.Location.Get(),
	)
	if err != nil {
		return err
	}
	paAccProc.Activate(true)
	paKfkCmp, err := ingestion.NewKafkaComponent(
		"passenger-identity",
		"driver-kafka-cmp",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KafkaPassengerIdentityTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		paAccProc,
		rt,
		rtw,
	)
	if err != nil {
		return err
	}

	// Handler for any identity change from Passegers and Drivers
	ipr := &identity.Processor{Mixpanel: mpCl}
	mph := mixpanel.NewHandler(ipr, []string{purl, durl}, cfm)
	mph.Run()

	prurl := "inproc://riderequest-publisher"
	paRRProc, err := passenger.NewRequestProcessor(prurl, "kafka", cfg.KkPRRTopic.Get())
	if err != nil {
		return err
	}
	paRRProc.Activate(true)
	paRRKfkCmp, err := ingestion.NewKafkaComponent(
		"passenger-request",
		"driver-kafka-cmp",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KkPRRTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		paRRProc,
		rt,
		rtw,
	)
	if err != nil {
		return err
	}
	// Handler for ride requests
	rrp := &riderequest.Processor{Mixpanel: mpCl}
	rrh := mixpanel.NewHandler(rrp, []string{prurl}, cfm)
	rrh.Run()

	pclurl := "inproc://riderequestcancel-publisher"
	paRCProc, err := passenger.NewCancellationProcessor(pclurl, "kafka", cfg.KkPRCTopic.Get())
	if err != nil {
		return err
	}
	paRCProc.Activate(true)
	paRCKfkCmp, err := ingestion.NewKafkaComponent(
		"passenger-cancel",
		"driver-kafka-cmp",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KkPRCTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		paRCProc,
		rt,
		rtw,
	)
	if err != nil {
		return err
	}

	// Handler for Request cancellations
	rrcp := &cancellation.Processor{Mixpanel: mpCl}
	rrch := mixpanel.NewHandler(rrcp, []string{pclurl}, cfm)
	rrch.Run()

	pakurl := "inproc://ride-publisher"
	paRAKProc, err := ride.NewRideProcessor(pakurl, "kafka", cfg.KkRTopic.Get())
	if err != nil {
		return err
	}
	paRAKProc.Activate(true)
	paRAKKfkCmp, err := ingestion.NewKafkaComponent(
		"ride",
		"driver-kafka-cmp",
		cfg.KkRTopic.Get(),
		[]string{cfg.KkPRCTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		paRAKProc,
		rt,
		rtw,
	)
	if err != nil {
		return err
	}
	// Handler for Confirmations of rides
	rrakp := &confirmation.Processor{Mixpanel: mpCl}
	rrakh := mixpanel.NewHandler(rrakp, []string{pakurl}, cfm)
	rrakh.Run()

	*ccmp = append(*ccmp, drKfkCmp, paKfkCmp, paRRKfkCmp, paRCKfkCmp, paRAKKfkCmp)

	plCfg := configuration.RestPoller{
		Manager:       cfm,
		RestURL:       cfg.RestURL.Get(),
		PollingPeriod: cfg.SettingsPeriod.Get(),
		DefaultConf:   defaultConf,
		RestKey:       cfg.RestKey.Get(),
		Path:          cfg.RestMixpanelPath.Get(),
	}
	plCfg.UpdateSettings()

	return nil
}

func setupSemoviComponents(cfg *config.Configuration, store *sql.Store, rrb *phttp.RoutesBuilder, ccmp *[]patron.Component) error {
	// Route handlers
	routes := semhttp.Routes(store)
	for _, r := range routes {
		rrb.Append(r)
	}

	// Kafka components
	rt := uint(10)
	rtw := 5 * time.Second

	osp := semovi.NewOperatorStatsProcessor(sql.NewOperatorStatsRepository(store))
	osp.Activate(true)
	ospKafka, err := ingestion.NewKafkaComponent(
		"stats_operador",
		"semovi-kafka-cmp",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KkSOTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		osp,
		rt,
		rtw,
	)
	if err != nil {
		return err
	}

	tip := semovi.NewTrafficIncidentsProcessor(sql.NewTrafficIncidentsRepository(store))
	tip.Activate(true)
	tipKafka, err := ingestion.NewKafkaComponent(
		"hecho_transito",
		"semovi-kafka-cmp",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KkHTTopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		tip,
		rt,
		rtw,
	)
	if err != nil {
		return err
	}

	vap := semovi.NewAggregatedTripsProcessor(sql.NewAggregatedTripsRepository(store))
	vap.Activate(true)
	vapKafka, err := ingestion.NewKafkaComponent(
		"viajes_agregados",
		"semovi-kafka-cmp",
		cfg.KafkaGroup.Get(),
		[]string{cfg.KkVATopic.Get()},
		[]string{cfg.KafkaBroker.Get()},
		vap,
		rt,
		rtw,
	)
	if err != nil {
		return err
	}

	*ccmp = append(*ccmp, ospKafka, tipKafka, vapKafka)

	return nil
}

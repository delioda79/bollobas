package main

import (
	"bollobas/configuration"
	"bollobas/ingestion"
	"bollobas/ingestion/driver"
	"bollobas/ingestion/passenger"
	"bollobas/ingestion/ride"
	"bollobas/mixpanel"
	"bollobas/mixpanel/identity"
	"bollobas/mixpanel/riderequest"
	"bollobas/mixpanel/riderequest/cancellation"
	"bollobas/mixpanel/riderequest/confirmation"
	"bollobas/pkg/ciphrest"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/beatlabs/patron/async"

	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"
	mpsdk "github.com/dukex/mixpanel"
	"github.com/joho/godotenv"
)

var (
	version = "0.0.1"
	kafkaBroker, kafkaDriverIdentityTopic, kafkaGroup,
	kafkaPassengerIdentityTopic, mpToken,
	kkPRRTopic, kkPRCTopic, kkRTopic string
	kafkaTimeout                            time.Duration
	defaultConf                             map[string]interface{}
	settingsPeriod                          time.Duration
	restKey, restURL, restMixpanelPath      string
	cipherKey, cipherInitVec, mixpanelToken string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Debugf("no .env file exists: %v", err)
	}

	kafkaBroker = mustGetEnv("BOLLOBAS_KAFKA_CONNECTION_STRING")
	kafkaDriverIdentityTopic = mustGetEnvWithDefault("BOLLOBAS_KAFKA_DRIVER_TOPIC", "driver_analytics")
	kafkaPassengerIdentityTopic = mustGetEnvWithDefault("BOLLOBAS_KAFKA_PASSENGER_TOPIC", "passenger_analytics")
	kafkaTimeout = mustGetEnvDurationWithDefault("BOLLOBAS_KAFKA_TIMEOUT", "2s")
	kafkaGroup = mustGetEnv("BOLLOBAS_KAFKA_GROUP")
	mpToken = mustGetEnv("BOLLOBAS_MIXPANEL_TOKEN")
	kkPRRTopic = mustGetEnvWithDefault("BOLLOBAS_KAFKA_REQUEST_TOPIC", "request")
	kkPRCTopic = mustGetEnvWithDefault("BOLLOBAS_KAFKA_REQUEST_CANCEL_TOPIC", "request_cancel")
	kkRTopic = mustGetEnvWithDefault("BOLLOBAS_KAFKA_RIDE_TOPIC", "ride")
	bConf := mustGetEnvWithDefault("BOLLOBAS_BASE_CONF", "{}")
	restKey = mustGetEnvWithDefault("REST_KEY", "")
	restURL = getEnvWithDefault("REST_CONNECTION_STRING", "")
	restMixpanelPath = mustGetEnvWithDefault("REST_MIXPANEL_PATH", "/taxidmin/bollobas/mixpanel-passenger-settings")
	cipherKey = mustGetEnvWithDefault("BOLLOBAS_CIPHER_KEY", "")
	cipherInitVec = mustGetEnvWithDefault("BOLLOBAS_INIT_VECTOR", "")

	defaultConf = map[string]interface{}{}
	err = json.Unmarshal([]byte(bConf), &defaultConf)
	if err != nil {
		panic(fmt.Sprintf("Wrong configuratiopn provided %v", err))
	}

	settingsPeriod = mustGetEnvDurationWithDefault("BOLLOBAS_SETTINGS_DURATION", "10s")

	ciphrest.InitCipher(cipherKey, cipherInitVec)
}

func main() {
	name := "bollobas"

	failure := async.ConsumerRetry(10, 5*time.Second)

	err := patron.Setup(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}
	log.Debugf("Starting %s v%s", name, version)

	durl := "inproc://driver-publisher"
	drAccProc, err := driver.NewAccountProcessor(durl)
	drKfkCmp, err := ingestion.NewKafkaComponent("driver-identity", kafkaBroker, kafkaDriverIdentityTopic, kafkaGroup, drAccProc, failure)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	purl := "inproc://passenger-publisher"
	paAccProc, err := passenger.NewAccountProcessor(purl, "kafka", kafkaPassengerIdentityTopic)
	paAccProc.Activate(true)
	paKfkCmp, err := ingestion.NewKafkaComponent("passenger-identity", kafkaBroker, kafkaPassengerIdentityTopic, kafkaGroup, paAccProc, failure)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	prurl := "inproc://riderequest-publisher"
	paRRProc, err := passenger.NewRequestProcessor(prurl, "kafka", kkPRRTopic)
	paRRProc.Activate(true)
	paRRKfkCmp, err := ingestion.NewKafkaComponent("passenger-request", kafkaBroker, kkPRRTopic, kafkaGroup, paRRProc, failure)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	pclurl := "inproc://riderequestcancel-publisher"
	paRCProc, err := passenger.NewCancellationProcessor(pclurl, "kafka", kkPRCTopic)
	paRCProc.Activate(true)
	paRCKfkCmp, err := ingestion.NewKafkaComponent("passenger-cancel", kafkaBroker, kkPRCTopic, kafkaGroup, paRCProc, failure)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	pakurl := "inproc://ride-publisher"
	paRAKProc, err := ride.NewRideProcessor(pakurl, "kafka", kkRTopic)
	paRAKProc.Activate(true)
	paRAKKfkCmp, err := ingestion.NewKafkaComponent("ride", kafkaBroker, kkRTopic, kafkaGroup, paRAKProc, failure)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	// MIXPANEL handlers
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}

	//Conf Manager MixPanel
	cfm := &mixpanel.Configurator{}
	cfm.Configure(defaultConf)

	plCfg := configuration.RestPoller{
		Manager:cfm,
		RestURL: restURL,
		PollingPeriod:settingsPeriod,
		DefaultConf:defaultConf,
		RestKey:restKey,
		Path:restMixpanelPath,
	}

	plCfg.UpdateSettings()

	mpCl := mpsdk.NewFromClient(c, mpToken, "")
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

	srv, err := patron.New(
		name,
		version,
		patron.Components(drKfkCmp, paKfkCmp, paRRKfkCmp, paRAKKfkCmp, paRCKfkCmp),
	)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	fmt.Println(v, ok, key)
	if !ok {
		fmt.Println("Exactly!", key)
		log.Fatalf("Missing configuration %s", key)
	}
	return v
}

func mustGetEnvWithDefault(key, def string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		if def == "" {
			log.Fatalf("Missing configuration %s", key)
		} else {
			return def
		}
	}
	return v
}

func mustGetEnvDurationWithDefault(key, def string) time.Duration {
	dur, err := time.ParseDuration(mustGetEnvWithDefault(key, def))
	if err != nil {
		log.Fatalf("env %s is not a duration: %v", key, err)
	}

	return dur
}

func getEnvWithDefault(key, def string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return v
}
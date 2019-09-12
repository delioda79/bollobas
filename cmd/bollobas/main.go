package main

import (
	"bollobas/driver"
	"fmt"
	"os"
	"time"

	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/log"
	"github.com/joho/godotenv"
)

var (
	version                                            = "dev"
	kafkaBroker, kafkaDriverIdentityTopic, kafkaGroup,
	kafkaPassengerIdentityTopic 					  string
	kafkaTimeout                                      time.Duration
)



func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Debugf("no .env file exists: %v", err)
	}

	kafkaBroker = mustGetEnv("BOLLOBAS_KAFKA_CONNECTION_STRING")
	kafkaDriverIdentityTopic = mustGetEnvWithDefault("BOLLOBAS_KAFKA_DRIVER_TOPIC", "driver_account")
	kafkaPassengerIdentityTopic = mustGetEnvWithDefault("BOLLOBAS_KAFKA_PASSENGER_TOPIC", "passenger_account")
	kafkaTimeout = mustGetEnvDurationWithDefault("BOLLOBAS_KAFKA_TIMEOUT", "2s")
	kafkaGroup = mustGetEnv("BOLLOBAS_KAFKA_GROUP")
}

func main() {
	name := "bollobas"

	err := patron.Setup(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	drKfkCmp, err := driver.NewKafkaComponent("driver-identity", kafkaBroker, kafkaDriverIdentityTopic, kafkaGroup)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	paKfkCmp, err := driver.NewKafkaComponent("passenger-identity", kafkaBroker, kafkaDriverIdentityTopic, kafkaGroup)
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	srv, err := patron.New(
		name,
		version,
		patron.Components(drKfkCmp, paKfkCmp),
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
		fmt.Println("Exactly!")
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
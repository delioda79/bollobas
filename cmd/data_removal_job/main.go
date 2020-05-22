package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/taxibeat/bollobas/internal/config"
	"github.com/taxibeat/bollobas/internal/storage/sql"
)

var tables = []string{
	"aggregated_trips",
	"operator_stats",
	"traffic_incidents",
}

func main() {
	log.Println("initializing clear job...")

	if len(os.Args) > 1 {
		if err := godotenv.Load(os.Args[1]); err != nil {
			log.Printf("cannot open given .env file: %v", err)
		}
	}

	cfg, err := setupConfig()
	if err != nil {
		log.Fatalf("failed to fetch configuration: %v\n", err)
	}
	if !cfg.DataRemovalEnabled.Get() {
		log.Println("functionality is disabled, exiting...")
		os.Exit(0)
	}

	// Setup SQL
	store, err := sql.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer store.Close()

	// boom
	if err := store.RemoveDataInTable(context.TODO(), tables...); err != nil {
		log.Fatal(err)
	}
	log.Println("all data are successfully removed, exiting...")
}

func setupConfig() (*config.Configuration, error) {
	cfg := &config.Configuration{}
	h, err := config.NewConfig(cfg)
	if err == nil {
		err = h.Harvest(context.Background())
	}

	return cfg, err
}

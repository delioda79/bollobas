package storagetest

import (
	"context"
	gsql "database/sql"
	"github.com/beatlabs/harvester"
	"github.com/taxibeat/bollobas/internal/config"
	"github.com/taxibeat/bollobas/internal/storage/sql"
	"log"
)

var (
	// Used to initialize repos
	store *sql.Store
	// Used to handle db in the background
	db *gsql.DB
)

// SetConfig returns a store
func SetConfig() (*sql.Store, error) {
	var err error
	cfg := &config.Configuration{}
	var h harvester.Harvester
	if h, err = config.NewConfig(cfg); err != nil {
		log.Fatal(err)
	}
	if err = h.Harvest(context.Background()); err != nil {
		log.Printf("failed to harvest configuration: %v", err)
	}
	// black box handle
	return sql.New(cfg)

}

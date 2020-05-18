package storagetest

import (
	"context"
	gsql "database/sql"
	"fmt"
	"github.com/beatlabs/harvester"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/config"
	"github.com/taxibeat/bollobas/internal/storage/sql"
	"log"
	"testing"
	"time"
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

// TestFilteredQuery is a general test to check if filtering works
func TestFilteredQuery(t *testing.T, f func(ctx context.Context, filter internal.DateFilter) (interface{}, error)) {

	nn := []struct {
		from  time.Time
		to    time.Time
		count int
	}{
		{
			from:  time.Now().Add(time.Hour * (-2)),
			to:    time.Now().Add(time.Hour * (-1)),
			count: 0,
		},
		{
			from:  time.Now().Add(time.Hour * (-2)),
			to:    time.Now().Add(time.Minute),
			count: 1,
		},
		{
			from:  time.Now().Add(time.Minute * 2),
			to:    time.Now().Add(time.Minute * 4),
			count: 0,
		},

		{
			from:  time.Now().Add(time.Minute * 20),
			to:    time.Now().Add(time.Hour + (time.Minute * 20)),
			count: 1,
		},
		{
			from:  time.Now().Add(time.Minute * (-4)),
			to:    time.Now().Add(time.Hour * 2),
			count: 2,
		},
	}

	dd := make([]struct {
		from  *time.Time
		to    *time.Time
		count int
	}, len(nn)+2)

	for i, n := range nn {
		ln := n
		dd[i] = struct {
			from  *time.Time
			to    *time.Time
			count int
		}{
			from:  &ln.from,
			to:    &ln.to,
			count: n.count,
		}
	}

	to := time.Now().Add(time.Minute * 30)
	from := time.Now().Add(time.Hour * (-2))

	dd[len(nn)] = struct {
		from  *time.Time
		to    *time.Time
		count int
	}{
		to:    &to,
		count: 1,
	}

	dd[len(nn)+1] = struct {
		from  *time.Time
		to    *time.Time
		count int
	}{
		from:  &from,
		count: 2,
	}

	for i, d := range dd {
		t.Run(fmt.Sprintf("Data %d", i), func(t *testing.T) {
			rr, err := f(context.Background(), internal.DateFilter{From: d.from, To: d.to})
			assert.Nil(t, err)

			assert.Len(t, rr, d.count)
		})

	}
}

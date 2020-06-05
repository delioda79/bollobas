package storagetest

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/config"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"github.com/beatlabs/harvester"
	"github.com/stretchr/testify/assert"
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
func TestFilteredQuery(t *testing.T, f func(ctx context.Context, filter internal.DateFilter) (interface{}, int, error)) {
	now := time.Now().AddDate(0, -1, 0)

	nn := []struct {
		from  time.Time
		to    time.Time
		count int
	}{
		{
			from:  now.Add(time.Hour * (-2)),
			to:    now.Add(time.Hour * (-1)),
			count: 0,
		},
		{
			from:  now.Add(time.Hour * (-2)),
			to:    now.Add(time.Minute),
			count: 1,
		},
		{
			from:  now.Add(time.Minute * 2),
			to:    now.Add(time.Minute * 4),
			count: 0,
		},

		{
			from:  now.Add(time.Minute * 20),
			to:    now.Add(time.Hour + (time.Minute * 20)),
			count: 1,
		},
		{
			from:  now.Add(time.Minute * (-4)),
			to:    now.Add(time.Hour * 2),
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

	to := now.Add(time.Minute * 30)
	from := now.Add(time.Hour * (-2))

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
			rr, pi, err := f(context.Background(), internal.DateFilter{From: d.from, To: d.to})
			assert.Nil(t, err)

			assert.NotNil(t, pi)
			assert.Len(t, rr, d.count)
		})

	}
}

// +build integration

package sql_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"github.com/taxibeat/bollobas/internal/storagetest"
	"testing"
	"time"
)

func TestGetAllAggregatedTrips(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewAggregatedTripsRepository(context.Background(), st)
	err = populateAggregatedTripsTable(at)
	assert.Nil(t, err)

	rr, err := at.GetAll(context.Background())
	assert.Nil(t, err)

	assert.Len(t, rr, 2)
	assert.Equal(t, 2, rr[0].ID)
	assert.Equal(t, "Test2", rr[0].SupplierID)
	assert.Equal(t, 1, rr[1].ID)
	assert.Equal(t, "Test1", rr[1].SupplierID)
}

func populateAggregatedTripsTable(r *sql.AggregatedTripsRepo) error {
	ctx := context.Background()
	r.DB().Exec(ctx, "TRUNCATE aggregated_trips")
	a := &internal.AggregatedTrips{
		Date:              time.Now(),
		SupplierID:        "Test1",
		TotalDistTraveled: 10.45,
	}

	if err := r.Add(ctx, a); err != nil {
		return err
	}
	a = &internal.AggregatedTrips{
		Date:              time.Now(),
		SupplierID:        "Test2",
		TotalDistTraveled: 10.46,
	}
	return r.Add(context.Background(), a)
}

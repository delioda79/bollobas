// +build integration

package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"fmt"
	"github.com/taxibeat/bollobas/internal/storagetest"
)

func TestGetAllAggregatedTrips(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewAggregatedTripsRepository(st)
	err = populateAggregatedTripsTable(at)
	assert.Nil(t, err)

	rr, pi, err := at.GetAll(context.Background(), internal.DateFilter{}, internal.Pagination{})
	assert.Nil(t, err)

	assert.NotNil(t, pi)
	assert.Len(t, rr, 2)
	assert.Equal(t, int64(2), rr[0].ID)
	assert.Equal(t, "Test2", *rr[0].SupplierID)
	assert.Equal(t, int64(1), rr[1].ID)
	assert.Equal(t, "Test1", *rr[1].SupplierID)
}

func TestFilteredTripsQuery(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewAggregatedTripsRepository(st)
	err = populateAggregatedTripsTable(at)
	assert.Nil(t, err)

	f := func(ctx context.Context, filter internal.DateFilter) (interface{}, int, error) {

		return at.GetAll(ctx, filter, internal.Pagination{})
	}

	storagetest.TestFilteredQuery(t, f)
}

func populateAggregatedTripsTable(r *sql.AggregatedTripsRepo) error {
	ctx := context.Background()
	r.DB().Exec(ctx, "TRUNCATE aggregated_trips")

	id1 := "Test1"
	id2 := "Test2"
	dist1 := 10.45
	dist2 := 10.46
	a := &internal.AggregatedTrips{
		Date:              time.Now().AddDate(0, -1, 0),
		SupplierID:        &id1,
		TotalDistTraveled: &dist1,
	}

	if err := r.Add(ctx, a); err != nil {
		return err
	}
	a = &internal.AggregatedTrips{
		Date:              time.Now().AddDate(0, -1, 0).Add(time.Hour),
		SupplierID:        &id2,
		TotalDistTraveled: &dist2,
	}

	fmt.Println(time.Now())
	return r.Add(context.Background(), a)
}

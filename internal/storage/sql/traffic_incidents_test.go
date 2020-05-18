// +build integration

package sql_test

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"testing"
	"time"

	"github.com/taxibeat/bollobas/internal/storagetest"
)

func TestGetAllTrafficIncidents(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewTrafficIncidentsRepository(st)
	err = populateTrafficIncidentsTable(at)
	assert.Nil(t, err)

	rr, err := at.GetAll(context.Background(), internal.DateFilter{})
	assert.Nil(t, err)

	assert.Len(t, rr, 2)

	assert.Equal(t, int64(2), rr[0].ID)
	assert.Equal(t, "222-BBB", rr[0].Plates)
	assert.Equal(t, "C456456456", rr[0].Licence)

	assert.Equal(t, int64(1), rr[1].ID)
	assert.Equal(t, "111-AAA", rr[1].Plates)
	assert.Equal(t, "C12312312", rr[1].Licence)
}

func TestFilteredIncidentsQuery(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewTrafficIncidentsRepository(st)
	err = populateTrafficIncidentsTable(at)
	assert.Nil(t, err)

	f := func(ctx context.Context, filter internal.DateFilter) (interface{}, error) {
		return at.GetAll(ctx, filter)
	}

	storagetest.TestFilteredQuery(t, f)
}

func populateTrafficIncidentsTable(r *sql.TrafficIncidentsRepo) error {
	ctx := context.Background()
	r.DB().Exec(ctx, "TRUNCATE traffic_incidents")
	a := &internal.TrafficIncident{
		Date:           time.Now(),
		Type:           3,
		Plates:         "111-AAA",
		Licence:        "C12312312",
		TravelTime:     "15-19",
		TravelDistance: "6000-8999",
		Coordinates:    "-99.829343, 19.716384",
	}

	if err := r.Add(ctx, a); err != nil {
		return err
	}
	a = &internal.TrafficIncident{
		Date:           time.Now().Add(time.Hour),
		Type:           3,
		Plates:         "222-BBB",
		Licence:        "C456456456",
		TravelTime:     "12-15",
		TravelDistance: "3000-5999",
		Coordinates:    "-99.43255, 19.6473825",
	}
	return r.Add(context.Background(), a)
}

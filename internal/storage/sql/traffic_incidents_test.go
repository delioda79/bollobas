// +build integration

package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/storage/sql"

	"github.com/taxibeat/bollobas/internal/storagetest"
)

func TestGetAllTrafficIncidents(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewTrafficIncidentsRepository(st)
	err = populateTrafficIncidentsTable(at)
	assert.Nil(t, err)

	rr, pi, err := at.GetAll(context.Background(), internal.DateFilter{}, internal.Pagination{})
	assert.Nil(t, err)

	assert.NotNil(t, pi)
	assert.Len(t, rr, 2)

	assert.Equal(t, int64(2), rr[0].ID)
	assert.Equal(t, "222-BBB", *rr[0].Plates)
	assert.Equal(t, "C456456456", *rr[0].Licence)

	assert.Equal(t, int64(1), rr[1].ID)
	assert.Equal(t, "111-AAA", *rr[1].Plates)
	assert.Equal(t, "C12312312", *rr[1].Licence)
}

func TestFilteredIncidentsQuery(t *testing.T) {
	st, err := storagetest.SetConfig()
	assert.Nil(t, err)
	at := sql.NewTrafficIncidentsRepository(st)
	err = populateTrafficIncidentsTable(at)
	assert.Nil(t, err)

	f := func(ctx context.Context, filter internal.DateFilter) (interface{}, int, error) {
		return at.GetAll(ctx, filter, internal.Pagination{})
	}

	storagetest.TestFilteredQuery(t, f)
}

func populateTrafficIncidentsTable(r *sql.TrafficIncidentsRepo) error {
	ctx := context.Background()
	r.DB().Exec(ctx, "TRUNCATE traffic_incidents")
	tp := 3
	p1 := "111-AAA"
	p2 := "222-BBB"
	l1 := "C12312312"
	l2 := "C456456456"
	tt1 := "15-19"
	tt2 := "12-15"
	td1 := "6000-8999"
	td2 := "3000-5999"
	c1 := "-99.829343, 19.716384"
	c2 := "-99.43255, 19.6473825"

	a := &internal.TrafficIncident{
		Date:           time.Now().AddDate(0, -1, 0),
		Type:           &tp,
		Plates:         &p1,
		Licence:        &l1,
		TravelTime:     &tt1,
		TravelDistance: &td1,
		Coordinates:    &c1,
		ProducedAt:     time.Now(),
	}

	if err := r.Add(ctx, a); err != nil {
		return err
	}
	a = &internal.TrafficIncident{
		Date:           time.Now().AddDate(0, -1, 0).Add(time.Hour),
		Type:           &tp,
		Plates:         &p2,
		Licence:        &l2,
		TravelTime:     &tt2,
		TravelDistance: &td2,
		Coordinates:    &c2,
		ProducedAt:     time.Now(),
	}
	return r.Add(context.Background(), a)
}

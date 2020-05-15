package http

import (
	"context"
	"errors"

	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/internalfakes"

	"testing"

	phttp "github.com/beatlabs/patron/component/http"
	"github.com/stretchr/testify/assert"
)

func defaultInput() (context.Context, *phttp.Request) {
	return context.Background(), &phttp.Request{}
}

func TestGetAggregatedRides(t *testing.T) {
	ctx, req := defaultInput()
	rp := &internalfakes.FakeAggregatedTripsRepository{}

	dd := []struct {
		trips []internal.AggregatedTrips
		err   error
	}{
		{nil, nil},
		{[]internal.AggregatedTrips{{ID: 1}}, nil},
		{nil, errors.New("an error")},
	}

	for i, d := range dd {
		rp.GetAllReturnsOnCall(i, d.trips, d.err)
		rsp, err := (&AggregatedRidesHandler{Rp: rp}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			assert.Equal(t, &phttp.Response{Payload: []internal.AggregatedTrips(d.trips)}, rsp)
		} else {
			assert.Equal(t, phttp.NewErrorWithCodeAndPayload(500, d.err), err)
			var r *phttp.Response
			assert.EqualValues(t, r, rsp)
		}

	}
}

func TestGetOperatorStats(t *testing.T) {
	ctx, req := defaultInput()
	rp := &internalfakes.FakeOperatorStatsRepository{}

	dd := []struct {
		trips []internal.OperatorStats
		err   error
	}{
		{nil, nil},
		{[]internal.OperatorStats{{ID: 1}}, nil},
		{nil, errors.New("an error")},
	}

	for i, d := range dd {
		rp.GetAllReturnsOnCall(i, d.trips, d.err)
		rsp, err := (&OperatorStatsHandler{Rp: rp}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			assert.Equal(t, &phttp.Response{Payload: []internal.OperatorStats(d.trips)}, rsp)
		} else {
			assert.Equal(t, phttp.NewErrorWithCodeAndPayload(500, d.err), err)
			var r *phttp.Response
			assert.EqualValues(t, r, rsp)
		}
	}
}

func TestGetTransitsMade(t *testing.T) {
	ctx, req := defaultInput()
	rp := &internalfakes.FakeTrafficIncidentsRepository{}

	dd := []struct {
		trips []internal.TrafficIncident
		err   error
	}{
		{nil, nil},
		{[]internal.TrafficIncident{{ID: 1}}, nil},
		{nil, errors.New("an error")},
	}

	for i, d := range dd {
		rp.GetAllReturnsOnCall(i, d.trips, d.err)
		rsp, err := (&TrafficIncidentsHandler{Rp: rp}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			assert.Equal(t, &phttp.Response{Payload: []internal.TrafficIncident(d.trips)}, rsp)
		} else {
			assert.Equal(t, phttp.NewErrorWithCodeAndPayload(500, d.err), err)
			var r *phttp.Response
			assert.EqualValues(t, r, rsp)
		}
	}
}

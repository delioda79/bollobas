package http

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/internalfakes"
	"github.com/taxibeat/bollobas/internal/semovi/rest/http/view"

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
		rsp, err := (&RouteHandler{Handler: &AggregatedRidesHandler{Rp: rp}}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			var v []view.AggregatedTrips
			if len(d.trips) > 0 {
				v = []view.AggregatedTrips{
					{
						ID:   d.trips[0].ID,
						Date: time.Time{}.Format("2006-01-02T15:04:05"),
					},
				}
			}
			assert.Equal(t, &phttp.Response{Payload: v}, rsp)
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
		rsp, err := (&RouteHandler{Handler: &OperatorStatsHandler{Rp: rp}}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			var vv []interface{}
			if len(d.trips) > 0 {
				v := view.OperatorStats{
					ID:   d.trips[0].ID,
					Date: d.trips[0].Date.Format("2006-01-02T15:04:05"),
				}
				vv = append(vv, v)
			}
			assert.Equal(t, &phttp.Response{Payload: vv}, rsp)
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
		rsp, err := (&RouteHandler{Handler: &TrafficIncidentsHandler{Rp: rp}}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			var vv []interface{}
			if len(d.trips) > 0 {
				v := view.TrafficIncident{
					ID:   d.trips[0].ID,
					Date: d.trips[0].Date.Format("2006-01-02T15:04:05"),
				}
				vv = append(vv, v)
			}
			assert.Equal(t, &phttp.Response{Payload: vv}, rsp)
		} else {
			assert.Equal(t, phttp.NewErrorWithCodeAndPayload(500, d.err), err)
			var r *phttp.Response
			assert.EqualValues(t, r, rsp)
		}
	}
}

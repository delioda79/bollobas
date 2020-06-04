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
		Meta  Metadata
		err   error
	}{
		{nil, Metadata{TotalCount: 0, TotalPages: 0, CurrentPage: 1, PageSize: 10}, nil},
		{[]internal.AggregatedTrips{{ID: 1}}, Metadata{TotalCount: 1, TotalPages: 1, CurrentPage: 1, PageSize: 10}, nil},
		{nil, Metadata{TotalCount: 0, TotalPages: 0, CurrentPage: 0, PageSize: 10}, errors.New("failed to fetch data")},
	}

	for i, d := range dd {
		rp.GetAllReturnsOnCall(i, d.trips, d.Meta.TotalCount, d.err)
		rsp, err := (&RouteHandler{Handler: &AggregatedRidesHandler{Rp: rp}}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			var vv []interface{}
			if len(d.trips) > 0 {
				for _, d := range d.trips {
					v := view.AggregatedTrips{
						ID:   d.ID,
						Date: time.Time{}.Format("2006-01-02T15:04:05"),
					}
					vv = append(vv, v)
				}
			}
			assert.Equal(t, &phttp.Response{
				Payload: PaginatedResponse{
					Meta: Metadata{
						TotalCount:  d.Meta.TotalCount,
						TotalPages:  d.Meta.TotalPages,
						PageSize:    d.Meta.PageSize,
						CurrentPage: d.Meta.CurrentPage,
					},
					Data: vv,
				},
			}, rsp)
		} else {
			assert.Equal(t, phttp.NewServiceUnavailableErrorWithPayload(d.err.Error()), err)
			var r *phttp.Response
			assert.EqualValues(t, r, rsp)
		}

	}
}

func TestGetOperatorStats(t *testing.T) {
	ctx, req := defaultInput()
	rp := &internalfakes.FakeOperatorStatsRepository{}

	dd := []struct {
		opStats []internal.OperatorStats
		Meta    Metadata
		err     error
	}{
		{nil, Metadata{TotalCount: 0, TotalPages: 0, CurrentPage: 1, PageSize: 10}, nil},
		{[]internal.OperatorStats{{ID: 1}}, Metadata{TotalCount: 1, TotalPages: 1, CurrentPage: 1, PageSize: 10}, nil},
		{nil, Metadata{TotalCount: 0, TotalPages: 0, CurrentPage: 1, PageSize: 10}, errors.New("failed to fetch data")},
	}

	for i, d := range dd {
		rp.GetAllReturnsOnCall(i, d.opStats, d.Meta.TotalCount, d.err)
		rsp, err := (&RouteHandler{Handler: &OperatorStatsHandler{Rp: rp}}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			var vv []interface{}
			if len(d.opStats) > 0 {
				for _, d := range d.opStats {
					v := view.OperatorStats{
						ID:   d.ID,
						Date: d.Date.Format("2006-01-02T15:04:05"),
					}
					vv = append(vv, v)
				}
			}
			assert.Equal(t, &phttp.Response{
				Payload: PaginatedResponse{
					Meta: Metadata{
						TotalCount:  d.Meta.TotalCount,
						TotalPages:  d.Meta.TotalPages,
						PageSize:    d.Meta.PageSize,
						CurrentPage: d.Meta.CurrentPage,
					},
					Data: vv,
				},
			}, rsp)
		} else {
			assert.Equal(t, phttp.NewServiceUnavailableErrorWithPayload(d.err.Error()), err)
			var r *phttp.Response
			assert.EqualValues(t, r, rsp)
		}
	}
}

func TestGetTransitsMade(t *testing.T) {
	ctx, req := defaultInput()
	rp := &internalfakes.FakeTrafficIncidentsRepository{}

	dd := []struct {
		ti   []internal.TrafficIncident
		Meta Metadata
		err  error
	}{
		{nil, Metadata{TotalCount: 0, TotalPages: 0, CurrentPage: 1, PageSize: 10}, nil},
		{[]internal.TrafficIncident{{ID: 1}}, Metadata{TotalCount: 1, TotalPages: 1, CurrentPage: 1, PageSize: 10}, nil},
		{nil, Metadata{TotalCount: 0, TotalPages: 0, CurrentPage: 1, PageSize: 10}, errors.New("failed to fetch data")},
	}

	for i, d := range dd {
		rp.GetAllReturnsOnCall(i, d.ti, d.Meta.TotalCount, d.err)
		rsp, err := (&RouteHandler{Handler: &TrafficIncidentsHandler{Rp: rp}}).Handle(ctx, req)

		if d.err == nil {
			assert.Equal(t, d.err, err)
			var vv []interface{}
			if len(d.ti) > 0 {
				for _, d := range d.ti {
					v := view.TrafficIncident{
						ID:   d.ID,
						Date: d.Date.Format("2006-01-02T15:04:05"),
					}
					vv = append(vv, v)
				}
			}
			assert.Equal(t, &phttp.Response{
				Payload: PaginatedResponse{
					Meta: Metadata{
						TotalCount:  d.Meta.TotalCount,
						TotalPages:  d.Meta.TotalPages,
						PageSize:    d.Meta.PageSize,
						CurrentPage: d.Meta.CurrentPage,
					},
					Data: vv,
				},
			}, rsp)
		} else {
			assert.Equal(t, phttp.NewServiceUnavailableErrorWithPayload(d.err.Error()), err)
			var r *phttp.Response
			assert.EqualValues(t, r, rsp)
		}
	}
}

package http

import (
	"context"
	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/semovi/rest/http/view"
	"github.com/taxibeat/bollobas/internal/storage/sql"
	"net/http"
	"strconv"
	"time"

	phttp "github.com/beatlabs/patron/component/http"
)

const (
	basePath = "/semovi"

	version = "/1.0.0" // the only true version
)

// Routes returns an array of all served routes
func Routes(st *sql.Store) []*phttp.RouteBuilder {
	or := sql.NewOperatorStatsRepository(st)
	ar := sql.NewAggregatedTripsRepository(st)
	tr := sql.NewTrafficIncidentsRepository(st)
	routes := [...]route{
		{http.MethodGet, "/viajes_agregados", &RouteHandler{Handler: &AggregatedRidesHandler{Rp: ar}}},
		{http.MethodGet, "/stats_operador", &RouteHandler{Handler: &OperatorStatsHandler{Rp: or}}},
		{http.MethodGet, "/hecho_transito", &RouteHandler{Handler: &TrafficIncidentsHandler{Rp: tr}}},
	}
	rb := make([]*phttp.RouteBuilder, len(routes))
	for i, r := range routes {
		rb[i] = r.ToPatronBuilder()
	}

	return rb
}

type handler interface {
	Handle(context.Context, *phttp.Request) (*phttp.Response, error)
}

type route struct {
	method   string
	endpoint string
	handler  handler
}

func (r *route) ToPatronBuilder() *phttp.RouteBuilder {
	uri := basePath + version + r.endpoint

	rb := phttp.NewRouteBuilder(uri, r.handler.Handle).WithTrace()

	switch r.method {
	case http.MethodGet:
		rb.MethodGet()
	}

	return rb
}

// Dates returns the date filter details
func getDateFilter(req *phttp.Request) (internal.DateFilter, error) {
	var fromP, toP *time.Time
	f := internal.DateFilter{}

	fromS, ok := req.Fields["from"]
	if ok {
		fromI, err := strconv.ParseInt(fromS, 10, 64)
		if err != nil {
			return f, err
		}

		from := time.Unix(fromI, 0)
		fromP = &from
	}

	toS, ok := req.Fields["to"]
	if ok {
		toI, err := strconv.ParseInt(toS, 10, 64)
		if err != nil {
			return f, err
		}

		to := time.Unix(toI, 0)
		toP = &to
	}

	return internal.DateFilter{From: fromP, To: toP}, nil
}

// AggregatedRidesHandler is the controller for the related route
type AggregatedRidesHandler struct {
	Rp internal.AggregatedTripsRepository
}

// GetAll returns all the items
func (a *AggregatedRidesHandler) GetAll(ctx context.Context, f internal.DateFilter) (interface{}, error) {
	ats, err := a.Rp.GetAll(ctx, f)
	if err != nil {
		return nil, err
	}

	var vats []view.AggregatedTrips
	for _, at := range ats {
		v := view.AggregatedTrips{
			ID:                     at.ID,
			Date:                   at.Date.Format("2006-01-02T15:04:05"),
			SupplierID:             at.SupplierID,
			TotalRides:             at.TotalRides,
			TotalVehicleRides:      at.TotalVehicleRides,
			TotalAvailableVehicles: at.TotalAvailableVehicles,
			TotalDistTraveled:      at.TotalDistTraveled,
			PassingTime:            at.PassingTime,
			RequestTime:            at.RequestTime,
			EmptyTime:              at.EmptyTime,
			EodMultiplier:          at.EodMultiplier,
			Accessibility:          at.Accessibility,
			FemaleOperator:         at.FemaleOperator,
			EodStart:               at.EodStart,
			EodEnd:                 at.EodEnd,
			EodPassDist:            at.EodPassDist,
			EodPassTime:            at.EodPassTime,
			RequestDist:            at.RequestDist,
			EmptyDist:              at.EmptyDist,
			EodRequestDist:         at.EodRequestDist,
			EodRequestTime:         at.EodRequestTime,
			EodEmptyDist:           at.EodEmptyDist,
			EodEmptyTime:           at.EodEmptyTime,
		}

		vats = append(vats, v)
	}

	return vats, nil
}

// OperatorStatsHandler is the controller for the related route
type OperatorStatsHandler struct {
	Rp internal.OperatorStatsRepository
}

// GetAll returns all the items
func (o *OperatorStatsHandler) GetAll(ctx context.Context, f internal.DateFilter) (interface{}, error) {
	return o.Rp.GetAll(ctx, f)
}

// TrafficIncidentsHandler is the controller for the related route
type TrafficIncidentsHandler struct {
	Rp internal.TrafficIncidentsRepository
}

// GetAll returns all the items
func (t *TrafficIncidentsHandler) GetAll(ctx context.Context, f internal.DateFilter) (interface{}, error) {
	return t.Rp.GetAll(ctx, f)
}

// DataHandler is a generic data handler which returns interfaces
type DataHandler interface {
	GetAll(ctx context.Context, f internal.DateFilter) (interface{}, error)
}

// RouteHandler is the controller for the related route
type RouteHandler struct {
	Handler DataHandler
}

// Handle handles the request
func (t *RouteHandler) Handle(ctx context.Context, req *phttp.Request) (*phttp.Response, error) {
	df, e := getDateFilter(req)
	if e != nil {
		return nil, phttp.NewErrorWithCodeAndPayload(500, e)
	}

	r, e := t.Handler.GetAll(ctx, df)
	if e != nil {
		return nil, phttp.NewErrorWithCodeAndPayload(500, e)
	}
	rsp := phttp.NewResponse(r)
	return rsp, nil
}

package http

import (
	"context"
	"net/http"

	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/storage/sql"

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
		route{http.MethodGet, "/viajes_agregados", &AggregatedRidesHandler{Rp: ar}},
		route{http.MethodGet, "/stats_operador", &OperatorStatsHandler{Rp: or}},
		route{http.MethodGet, "/hecho_transito", &TrafficIncidentsHandler{Rp: tr}},
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

// AggregatedRidesHandler is the controller for the related route
type AggregatedRidesHandler struct {
	Rp internal.AggregatedTripsRepository
}

// Handle handles the request
func (a AggregatedRidesHandler) Handle(ctx context.Context, req *phttp.Request) (*phttp.Response, error) {
	r, e := a.Rp.GetAll(ctx)
	if e != nil {
		return nil, phttp.NewErrorWithCodeAndPayload(500, e)
	}
	rsp := phttp.NewResponse(r)
	return rsp, nil
}

// OperatorStatsHandler is the controller for the related route
type OperatorStatsHandler struct {
	Rp internal.OperatorStatsRepository
}

// Handle handles the request
func (o *OperatorStatsHandler) Handle(ctx context.Context, req *phttp.Request) (*phttp.Response, error) {
	r, e := o.Rp.GetAll(ctx)
	if e != nil {
		return nil, phttp.NewErrorWithCodeAndPayload(500, e)
	}
	rsp := phttp.NewResponse(r)
	return rsp, nil
}

// TrafficIncidentsHandler is the controller for the related route
type TrafficIncidentsHandler struct {
	Rp internal.TrafficIncidentsRepository
}

// Handle handles the request
func (t *TrafficIncidentsHandler) Handle(ctx context.Context, req *phttp.Request) (*phttp.Response, error) {
	r, e := t.Rp.GetAll(ctx)
	if e != nil {
		return nil, phttp.NewErrorWithCodeAndPayload(500, e)
	}
	rsp := phttp.NewResponse(r)
	return rsp, nil
}

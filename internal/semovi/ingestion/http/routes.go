package http

import (
	"context"

	"net/http"

	phttp "github.com/beatlabs/patron/component/http"
)

const (
	basePath = "/semovi"

	version = "/1.0.0" // the only true version
)

var routes = [...]route{
	route{http.MethodGet, "/viajes_agregados", getAggregatedRides},
	route{http.MethodGet, "/stats_operador", getOperatorStats},
	route{http.MethodGet, "/hecho_transito", getTransitsMade},
}

// Routes returns an array of all served routes
func Routes() []*phttp.RouteBuilder {
	rb := make([]*phttp.RouteBuilder, len(routes))
	for i, r := range routes {
		rb[i] = r.ToPatronBuilder()
	}

	return rb
}

type route struct {
	method   string
	endpoint string
	handler  func(context.Context, *phttp.Request) (*phttp.Response, error)
}

func (r *route) ToPatronBuilder() *phttp.RouteBuilder {
	uri := basePath + version + r.endpoint

	rb := phttp.NewRouteBuilder(uri, r.handler).WithTrace()

	switch r.method {
	case http.MethodGet:
		rb.MethodGet()
	}

	return rb
}

func getAggregatedRides(ctx context.Context, req *phttp.Request) (*phttp.Response, error) {
	return phttp.NewResponse(nil), nil
}

func getOperatorStats(ctx context.Context, req *phttp.Request) (*phttp.Response, error) {
	return phttp.NewResponse(nil), nil
}

func getTransitsMade(ctx context.Context, req *phttp.Request) (*phttp.Response, error) {
	return phttp.NewResponse(nil), nil
}

package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/taxibeat/bollobas/internal"
	"github.com/taxibeat/bollobas/internal/semovi/rest/http/view"
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

// Dates returns the date filter details
func getPagination(req *phttp.Request) (internal.Pagination, error) {

	f := internal.Pagination{}

	first, err := getIntField(req, "first", 0)
	if err != nil {
		return f, err
	}

	count, err := getIntField(req, "count", 10)
	if err != nil {
		return f, err
	}

	f.First = int(first)
	f.Count = int(count)

	return f, nil
}

// AggregatedRidesHandler is the controller for the related route
type AggregatedRidesHandler struct {
	Rp internal.AggregatedTripsRepository
}

// GetAll returns all the items
func (a *AggregatedRidesHandler) GetAll(ctx context.Context, f internal.DateFilter, pg internal.Pagination) (interface{}, error) {
	ats, err := a.Rp.GetAll(ctx, f, pg)
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
func (o *OperatorStatsHandler) GetAll(ctx context.Context, f internal.DateFilter, pg internal.Pagination) (interface{}, error) {
	ops, err := o.Rp.GetAll(ctx, f, pg)
	if err != nil {
		return nil, err
	}

	var opsIntf []interface{}
	for _, op := range ops {
		v := view.OperatorStats{
			ID:             op.ID,
			Date:           op.Date.Format("2006-01-02T15:04:05"),
			OperatorID:     op.OperatorID,
			Gender:         op.Gender,
			CompletedTrips: op.CompletedTrips,
			DaysSince:      op.DaysSince,
			AgeRange:       op.AgeRange,
			HoursConnected: op.HoursConnected,
			TripHours:      op.TripHours,
			TotRevenue:     op.TotRevenue,
		}
		opsIntf = append(opsIntf, v)
	}

	return opsIntf, nil
}

// TrafficIncidentsHandler is the controller for the related route
type TrafficIncidentsHandler struct {
	Rp internal.TrafficIncidentsRepository
}

// GetAll returns all the items
func (t *TrafficIncidentsHandler) GetAll(ctx context.Context, f internal.DateFilter, pg internal.Pagination) (interface{}, error) {
	tis, err := t.Rp.GetAll(ctx, f, pg)
	if err != nil {
		return nil, err
	}

	var tisIntf []interface{}
	for _, ti := range tis {
		v := view.TrafficIncident{
			ID:             ti.ID,
			Date:           ti.Date.Format("2006-01-02T15:04:05"),
			Type:           ti.Type,
			Plates:         ti.Plates,
			Licence:        ti.Licence,
			TravelDistance: ti.TravelDistance,
			TravelTime:     ti.TravelTime,
			Coordinates:    ti.Coordinates,
		}
		tisIntf = append(tisIntf, v)
	}

	return tisIntf, nil
}

// DataHandler is a generic data handler which returns interfaces
type DataHandler interface {
	GetAll(ctx context.Context, f internal.DateFilter, pg internal.Pagination) (interface{}, error)
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

	pn, e := getPagination(req)
	if e != nil {
		return nil, phttp.NewErrorWithCodeAndPayload(400, e.Error())
	}

	r, e := t.Handler.GetAll(ctx, df, pn)
	if e != nil {
		return nil, phttp.NewErrorWithCodeAndPayload(500, e)
	}
	var nxt *int
	if len(r.([]interface{})) == pn.Count {
		nv := pn.First + pn.Count
		nxt = &nv
	}

	mdr := Metadata{
		First: pn.First,
		Next:  nxt,
	}
	rsp := phttp.NewResponse(PaginatedResponse{
		Meta: mdr,
		Data: r,
	})
	return rsp, nil
}

// Metadata are the metadata for the response
type Metadata struct {
	First int
	Next  *int
}

// PaginatedResponse is the response with paginated data
type PaginatedResponse struct {
	Meta Metadata
	Data interface{}
}

func getIntField(req *phttp.Request, param string, dv int64) (int64, error) {
	intText, ok := req.Fields[param]
	if !ok {
		return dv, nil
	}
	intVal, err := strconv.ParseInt(intText, 10, 64)
	if err != nil {
		return dv, fmt.Errorf("%s is not valid integer", param)
	}
	return intVal, nil
}

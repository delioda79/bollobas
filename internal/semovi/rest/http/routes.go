package http

import (
	"context"
	"fmt"
	"github.com/beatlabs/patron/log"
	"math"
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

	f.From = fromP
	f.To = toP
	err := f.Validate()
	if err != nil {
		return f, err
	}

	return f, nil
}

// Dates returns the date filter details
func getPagination(req *phttp.Request) (internal.Pagination, error) {

	f := internal.Pagination{}

	offset, err := getIntField(req, "offset", 0)
	if err != nil {
		return f, err
	}

	limit, err := getIntField(req, "limit", 10)
	if err != nil {
		return f, err
	}

	f.Offset = int(offset)
	f.Limit = int(limit)

	return f, nil
}

// AggregatedRidesHandler is the controller for the related route
type AggregatedRidesHandler struct {
	Rp internal.AggregatedTripsRepository
}

// GetAll godoc
// @Summary Get All the aggregated rides
// @Tags bollobas
// @Accept json
// @Produce json
// @Param from query int false "start date (epoch time)"
// @Param to query int false "end date (epoch time)"
// @Param limit query int false "Limit Value"
// @Param offset query int false "Offset Value"
// @Success 200 {array} view.AggregatedTrips
// @Failure 400 {object} view.ErrorSwagger
// @Failure 500 {object} view.ErrorSwagger
// @Router /viajes_agregados [get]
func (a *AggregatedRidesHandler) GetAll(ctx context.Context, f internal.DateFilter, pg internal.Pagination) ([]interface{}, int, error) {
	ats, pi, err := a.Rp.GetAll(ctx, f, pg)
	if err != nil {
		return nil, 0, err
	}

	vats := make([]interface{}, 0)
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

	return vats, pi, nil
}

// OperatorStatsHandler is the controller for the related route
type OperatorStatsHandler struct {
	Rp internal.OperatorStatsRepository
}

// GetAll godoc
// @Summary Get return all the operator stats
// @Tags bollobas
// @Accept json
// @Produce json
// @Param limit query int false "Limit Value"
// @Param offset query int false "Offset Value"
// @Success 200 {array} view.OperatorStats
// @Failure 400 {object} view.ErrorSwagger
// @Failure 500 {object} view.ErrorSwagger
// @Router /stats_operador [get]
func (o *OperatorStatsHandler) GetAll(ctx context.Context, f internal.DateFilter, pg internal.Pagination) ([]interface{}, int, error) {
	ops, pi, err := o.Rp.GetAll(ctx, f, pg)
	if err != nil {
		return nil, 0, err
	}

	opsIntf := make([]interface{}, 0)
	for _, op := range ops {
		v := view.OperatorStats{
			ID:             op.ID,
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

	return opsIntf, pi, nil
}

// TrafficIncidentsHandler is the controller for the related route
type TrafficIncidentsHandler struct {
	Rp internal.TrafficIncidentsRepository
}

// GetAll godoc
// @Summary Get return all the traffic incidents
// @Tags bollobas
// @Accept json
// @Produce json
// @Param from query int false "start date (epoch time)"
// @Param to query int false "end date (epoch time)"
// @Param limit query int false "Limit Value"
// @Param offset query int false "Offset Value"
// @Success 200 {array} view.TrafficIncident
// @Failure 400 {object} view.ErrorSwagger
// @Failure 500 {object} view.ErrorSwagger
// @Router /hecho_transito [get]
func (t *TrafficIncidentsHandler) GetAll(ctx context.Context, f internal.DateFilter, pg internal.Pagination) ([]interface{}, int, error) {
	tis, pi, err := t.Rp.GetAll(ctx, f, pg)
	if err != nil {
		return nil, 0, err
	}

	tisIntf := make([]interface{}, 0)
	for _, ti := range tis {
		v := view.TrafficIncident{
			ID:             ti.ID,
			Type:           ti.Type,
			Plates:         ti.Plates,
			Licence:        ti.Licence,
			TravelDistance: ti.TravelDistance,
			TravelTime:     ti.TravelTime,
			Coordinates:    ti.Coordinates,
			Date:           ti.Date.Format("2006-01-02T15:04:05"),
		}
		tisIntf = append(tisIntf, v)
	}

	return tisIntf, pi, nil
}

// DataHandler is a generic data handler which returns interfaces
type DataHandler interface {
	GetAll(ctx context.Context, f internal.DateFilter, pg internal.Pagination) ([]interface{}, int, error)
}

// RouteHandler is the controller for the related route
type RouteHandler struct {
	Handler DataHandler
}

// Handle handles the request
func (t *RouteHandler) Handle(ctx context.Context, req *phttp.Request) (*phttp.Response, error) {
	df, e := getDateFilter(req)
	if e != nil {
		return nil, phttp.NewValidationErrorWithPayload(e.Error())
	}

	pn, e := getPagination(req)
	if e != nil {
		return nil, phttp.NewValidationErrorWithPayload(e.Error())
	}

	r, td, e := t.Handler.GetAll(ctx, df, pn)
	if e != nil {
		log.Error(e.Error())
		return nil, phttp.NewServiceUnavailableErrorWithPayload("failed to fetch data")
	}

	totalPages := int(math.Ceil(float64(td) / float64(pn.CalcLimit())))
	mdr := Metadata{
		TotalCount:  td,
		TotalPages:  totalPages,
		PageSize:    pn.CalcLimit(),
		CurrentPage: (pn.CalcOffset() / pn.CalcLimit()) + 1,
	}
	rsp := phttp.NewResponse(PaginatedResponse{
		Data: r,
		Meta: mdr,
	})
	return rsp, nil
}

// Metadata are the metadata for the response
type Metadata struct {
	TotalCount  int
	TotalPages  int
	PageSize    int
	CurrentPage int
}

// PaginatedResponse is the response with paginated data
type PaginatedResponse struct {
	Data []interface{}
	Meta Metadata
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

package sql

import (
	"context"
	"github.com/taxibeat/bollobas/internal"
)

// AggregatedTripsRepo implements the interface for MySQL
type AggregatedTripsRepo struct {
	*Store
}

// GetAll returns the city with the respective id or an error if it does not exist
func (va *AggregatedTripsRepo) GetAll(ctx context.Context, df internal.DateFilter, pg internal.Pagination) (data []internal.AggregatedTrips, err error) {
	f := AllFilter{
		DateFilter: df,
		Pagination: pg,
	}
	var args []interface{}

	query := `SELECT
			id,
			date,
			supplier_id,
			total_rides,
			total_vehicle_rides,
			total_available_vehicles,
			total_dist_traveled,
			passing_time,
			request_time,
			empty_time,
			eod_multiplier,
			accessibility,
			female_operator,
			eod_start,
			eod_end,
			eod_pass_dist,
			eod_pass_time,
			request_dist,
			empty_dist,
			eod_request_dist,
			eod_request_time,
			eod_empty_dist,
			eod_empty_time
		FROM aggregated_trips
		WHERE 1=1 %s
			AND YEAR(date) = YEAR(CURRENT_DATE - INTERVAL 1 MONTH)
			AND MONTH(date) = MONTH(CURRENT_DATE - INTERVAL 1 MONTH)
			AND deleted_at is null
		ORDER BY date DESC
		LIMIT ?,?`

	query, a := f.FilterDate(query)
	args = append(args, a...)
	a = f.Paginate()
	args = append(args, a...)

	rr, err := va.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rr.Close()

	var res []internal.AggregatedTrips
	for rr.Next() {
		r := &internal.AggregatedTrips{}
		err := rr.Scan(
			&r.ID,
			&r.Date,
			&r.SupplierID,
			&r.TotalRides,
			&r.TotalVehicleRides,
			&r.TotalAvailableVehicles,
			&r.TotalDistTraveled,
			&r.PassingTime,
			&r.RequestTime,
			&r.EmptyTime,
			&r.EodMultiplier,
			&r.Accessibility,
			&r.FemaleOperator,
			&r.EodStart,
			&r.EodEnd,
			&r.EodPassDist,
			&r.EodPassTime,
			&r.RequestDist,
			&r.EmptyDist,
			&r.EodRequestDist,
			&r.EodRequestTime,
			&r.EodEmptyDist,
			&r.EodEmptyTime,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, *r)
	}
	return res, nil
}

// Add inserts a new record
func (va *AggregatedTripsRepo) Add(ctx context.Context, r *internal.AggregatedTrips) error {
	q := "INSERT INTO aggregated_trips  " +
		"(" +
		"date, " +
		"supplier_id, " +
		"total_rides, " +
		"total_vehicle_rides, " +
		"total_available_vehicles, " +
		"total_dist_traveled, " +
		"passing_time, " +
		"request_time, " +
		"empty_time, " +
		"eod_multiplier, " +
		"accessibility, " +
		"female_operator, " +
		"eod_start, " +
		"eod_end, " +
		"eod_pass_dist, " +
		"eod_pass_time, " +
		"request_dist, " +
		"empty_dist, " +
		"eod_request_dist, " +
		"eod_request_time, " +
		"eod_empty_dist, " +
		"eod_empty_time " +
		") " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	rr, err := va.db.Exec(ctx, q,
		r.Date,
		r.SupplierID,
		r.TotalRides,
		r.TotalVehicleRides,
		r.TotalAvailableVehicles,
		r.TotalDistTraveled,
		r.PassingTime,
		r.RequestTime,
		r.EmptyTime,
		r.EodMultiplier,
		r.Accessibility,
		r.FemaleOperator,
		r.EodStart,
		r.EodEnd,
		r.EodPassDist,
		r.EodPassTime,
		r.RequestDist,
		r.EmptyDist,
		r.EodRequestDist,
		r.EodRequestTime,
		r.EodEmptyDist,
		r.EodEmptyTime,
	)
	if err != nil {
		return err
	}
	r.ID, err = rr.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// NewAggregatedTripsRepository creates a new repo
func NewAggregatedTripsRepository(store *Store) *AggregatedTripsRepo {
	return &AggregatedTripsRepo{store}
}

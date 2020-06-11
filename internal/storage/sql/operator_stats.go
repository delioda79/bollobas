package sql

import (
	"context"

	"github.com/taxibeat/bollobas/internal"
)

// GetOperatorStatsQuery query
const GetOperatorStatsQuery = `SELECT
			id,
			operator_id,
			gender,
			completed_trips,
			days_since,
			age_range,
			hours_connected,
			trip_hours,
			tot_revenue
		FROM operator_stats
		WHERE 1=1 %s
			AND deleted_at is null
		ORDER BY id ASC
		LIMIT ?,?`

// OperatorStatsRepo implements the interface for MySQL
type OperatorStatsRepo struct {
	*Store
}

// GetAll returns the city with the respective id or an error if it does not exist
func (va *OperatorStatsRepo) GetAll(ctx context.Context, df internal.DateFilter, pg internal.Pagination) (data []internal.OperatorStats, totalCount int, err error) {
	f := AllFilter{
		DateFilter: df,
		Pagination: pg,
	}
	var args []interface{}

	query := GetOperatorStatsQuery

	query, a := f.FilterDate(query)
	args = append(args, a...)
	a = f.Paginate()
	args = append(args, a...)

	rr, err := va.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rr.Close()

	var res []internal.OperatorStats
	for rr.Next() {
		r := &internal.OperatorStats{}
		err := rr.Scan(
			&r.ID,
			&r.OperatorID,
			&r.Gender,
			&r.CompletedTrips,
			&r.DaysSince,
			&r.AgeRange,
			&r.HoursConnected,
			&r.TripHours,
			&r.TotRevenue,
		)
		if err != nil {
			return nil, 0, err
		}

		res = append(res, *r)
	}
	if err = rr.Err(); err != nil {
		return nil, 0, err
	}

	totalCount, err = va.getTotalCount(ctx, df)
	if err != nil {
		return nil, 0, err
	}

	return res, totalCount, nil
}

func (va *OperatorStatsRepo) getTotalCount(ctx context.Context, df internal.DateFilter) (int, error) {

	f := AllFilter{
		DateFilter: df,
	}
	var args []interface{}

	query := GetOperatorStatsQuery

	sqlCount := getSQLCountStmt(query)

	query, a := f.FilterDate(sqlCount)
	args = append(args, a...)

	var n int
	err := va.db.QueryRow(ctx, query, args...).Scan(&n)
	if err != nil {
		return n, err
	}

	return n, nil
}

// Add inserts a new record
func (va *OperatorStatsRepo) Add(ctx context.Context, r *internal.OperatorStats) error {
	q := `INSERT INTO operator_stats (
			operator_id, 
			gender, 
			completed_trips, 
			days_since, 
			age_range, 
			hours_connected, 
			trip_hours, 
			tot_revenue,
			produced_at
		) 
		VALUES (?,?,?,?,?,?,?,?,?)`

	rr, err := va.db.Exec(ctx, q,
		&r.OperatorID,
		&r.Gender,
		&r.CompletedTrips,
		&r.DaysSince,
		&r.AgeRange,
		&r.HoursConnected,
		&r.TripHours,
		&r.TotRevenue,
		&r.ProducedAt,
	)
	if err == nil {
		r.ID, err = rr.LastInsertId()
	}

	return err
}

// NewOperatorStatsRepository creates a new repo
func NewOperatorStatsRepository(store *Store) *OperatorStatsRepo {
	return &OperatorStatsRepo{store}
}

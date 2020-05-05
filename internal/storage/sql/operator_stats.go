package sql

import (
	"context"
	"github.com/taxibeat/bollobas/internal"
)

// OperatorStatsRepo implements the interface for MySQL
type OperatorStatsRepo struct {
	context.Context
	*Store
	table string
}

// GetAll returns the city with the respective id or an error if it does not exist
func (va *OperatorStatsRepo) GetAll(ctx context.Context) (data []internal.OperatorStats, err error) {
	rr, err := va.db.Query(ctx, "SELECT * from operator_stats  WHERE date <= ? ORDER BY date DESC", "NOW()")
	if err != nil {
		return nil, err
	}
	defer rr.Close()

	res := []internal.OperatorStats{}

	for rr.Next() {
		r := &internal.OperatorStats{}
		err := rr.Scan(
			&r.ID,
			&r.Date,
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
			return nil, err
		}

		res = append(res, *r)
	}
	return res, nil
}

// Add inserts a new record
func (va *OperatorStatsRepo) Add(ctx context.Context, r *internal.OperatorStats) error {
	q := "INSERT INTO operator_stats  " +
		"(" +
		"date, " +
		"operator_id, " +
		"gender, " +
	"completed_trips, " +
	"days_since, " +
	"age_range, " +
	"hours_connected, " +
	"trip_hours, " +
	"tot_revenue" +
		") " +
		"VALUES (?,?,?,?,?,?,?,?,?)"

	rr, err := va.db.Exec(ctx, q,
		&r.Date,
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
		return err
	}
	r.ID, err = rr.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// NewOperatorStatsRepository creates a new repo
func NewOperatorStatsRepository(ctx context.Context, store *Store) *OperatorStatsRepo {
	return &OperatorStatsRepo{ctx, store, "operator_stats"}
}

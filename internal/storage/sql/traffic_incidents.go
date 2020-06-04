package sql

import (
	"context"

	"github.com/taxibeat/bollobas/internal"
)

// GetTrafficIncidentsQuery query
const GetTrafficIncidentsQuery = `SELECT
			id,
			date,
			type,
			plates,
			licence,
			travel_distance,
			travel_time,
			coordinates
		FROM traffic_incidents
		WHERE 1=1 %s
			AND YEAR(date) = YEAR(CURRENT_DATE - INTERVAL 1 MONTH)
			AND MONTH(date) = MONTH(CURRENT_DATE - INTERVAL 1 MONTH)
			AND deleted_at is null
		ORDER BY date DESC, id ASC
		LIMIT ?,?`

// TrafficIncidentsRepo implements the interface for MySQL
type TrafficIncidentsRepo struct {
	*Store
}

// GetAll returns all the traffic incidents
func (ti *TrafficIncidentsRepo) GetAll(ctx context.Context, df internal.DateFilter, pg internal.Pagination) (data []internal.TrafficIncident, totalCount int, err error) {
	f := AllFilter{
		DateFilter: df,
		Pagination: pg,
	}
	var args []interface{}

	query := GetTrafficIncidentsQuery

	query, a := f.FilterDate(query)
	args = append(args, a...)
	a = f.Paginate()
	args = append(args, a...)

	ii, err := ti.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer ii.Close()

	var res []internal.TrafficIncident
	for ii.Next() {
		i := &internal.TrafficIncident{}
		err := ii.Scan(
			&i.ID,
			&i.Date,
			&i.Type,
			&i.Plates,
			&i.Licence,
			&i.TravelDistance,
			&i.TravelTime,
			&i.Coordinates,
		)
		if err != nil {
			return nil, 0, err
		}

		res = append(res, *i)
	}
	if err = ii.Err(); err != nil {
		return nil, 0, err
	}

	totalCount, err = ti.getTotalCount(ctx, df)
	if err != nil {
		return nil, 0, err
	}

	return res, totalCount, nil
}

func (ti *TrafficIncidentsRepo) getTotalCount(ctx context.Context, df internal.DateFilter) (int, error) {

	f := AllFilter{
		DateFilter: df,
	}
	var args []interface{}

	query := GetTrafficIncidentsQuery

	sqlCount := getSQLCountStmt(query)

	query, a := f.FilterDate(sqlCount)
	args = append(args, a...)

	var n int
	err := ti.db.QueryRow(ctx, query, args...).Scan(&n)
	if err != nil {
		return n, err
	}

	return n, nil
}

// Add inserts a new record
func (ti *TrafficIncidentsRepo) Add(ctx context.Context, i *internal.TrafficIncident) error {
	q := "INSERT INTO traffic_incidents  " +
		"(" +
		"date, " +
		"type, " +
		"plates, " +
		"licence, " +
		"travel_distance, " +
		"travel_time, " +
		"coordinates " +
		") " +
		"VALUES (?,?,?,?,?,?,?)"

	rr, err := ti.db.Exec(ctx, q,
		i.Date,
		i.Type,
		i.Plates,
		i.Licence,
		i.TravelDistance,
		i.TravelTime,
		i.Coordinates,
	)
	if err != nil {
		return err
	}
	i.ID, err = rr.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

// NewTrafficIncidentsRepository creates a new repo
func NewTrafficIncidentsRepository(store *Store) *TrafficIncidentsRepo {
	return &TrafficIncidentsRepo{store}
}

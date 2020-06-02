package sql

import (
	"context"

	"github.com/taxibeat/bollobas/internal"
)

// TrafficIncidentsRepo implements the interface for MySQL
type TrafficIncidentsRepo struct {
	*Store
	table string
}

// GetAll returns all the traffic incidents
func (ti *TrafficIncidentsRepo) GetAll(ctx context.Context, df internal.DateFilter) (data []internal.TrafficIncident, err error) {
	f := DateFilter{&df}

	query := `SELECT
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
		ORDER BY date DESC`

	query, args := f.Filter(query)

	ii, err := ti.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		res = append(res, *i)
	}
	return res, nil
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
	return &TrafficIncidentsRepo{store, "traffic_incidents"}
}

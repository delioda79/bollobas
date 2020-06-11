package internal

import (
	"context"
	"time"
)

// TrafficIncidentsRepository allows to manipulate the traffic incidents
type TrafficIncidentsRepository interface {
	GetAll(ctx context.Context, df DateFilter, pg Pagination) ([]TrafficIncident, int, error)
	Add(ctx context.Context, trips *TrafficIncident) error
}

// TrafficIncident represents the stats for traffic incident report
type TrafficIncident struct {
	ID             int64     `json:"id"`
	Type           *int      `transl:"type" json:"hecho_trans"`
	Plates         *string   `transl:"plates" json:"placa"`
	Licence        *string   `transl:"licence" json:"licencia"`
	TravelDistance *string   `transl:"travel_distance" json:"distancia_viaje"`
	TravelTime     *string   `transl:"travel_time" json:"tiempo_viaje"`
	Coordinates    *string   `transl:"coordinates" json:"ubicacion"`
	Date           time.Time `transl:"date" json:"tiempo_hecho"`
	ProducedAt     time.Time `json:"fecha_produccion"`
	DeletedAt      time.Time `transl:"deleted_at"`
}

package internal

import (
	"context"
	"time"
)

// TrafficIncidentsRepository allows to manipulate the traffic incidents
type TrafficIncidentsRepository interface {
	GetAll(ctx context.Context) ([]TrafficIncident, error)
	Add(ctx context.Context, trips *TrafficIncident) error
}

// TrafficIncident represents the stats for traffic incident report
type TrafficIncident struct {
	ID             int64     `json:"id"`
	Date           time.Time `json:"date" transl:"tiempo_hecho"`
	Type           int       `json:"type" transl:"hecho_trans"`
	Plates         string    `json:"plates" transl:"placa"`
	Licence        string    `json:"licence" transl:"licencia"`
	TravelDistance string    `json:"travel_distance" transl:"distancia_viaje"`
	TravelTime     string    `json:"travel_time" transl:"tiempo_viaje"`
	Coordinates    string    `json:"coordinates" transl:"ubicación"`
}

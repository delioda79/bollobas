package internal

import "time"

// OperatorStatsRepository allows to manipulate the operator stats
type OperatorStatsRepository interface {
	GetAll() ([]AggregatedTrips, error)
	Add(trips *AggregatedTrips)
}

// OperatorStats represents the stats for report 2
type OperatorStats struct {
	ID int64 `json:"id"`
	Date time.Time `json:"date" transl:"fecha_produccion"`
	OperatorID int `json:"operator_id" transl:"id_operador"`
	Gender int `json:"gender" transl:"genero"`
	CompletedTrips int `json:"completed_trips" transl:"cant_viajes"`
	DaysSince int `json:"days_since" transl:"tiempo_registro"`
	AgeRange string `json:"age_range" transl:"edad"`
	HoursConnected string `json:"hours_connected" transl:"horas_conectado"`
	TripHours string `json:"trip_hours" transl:"horas_viaje"`
	TotRevenue string `json:"tot_revenue" transl:"ingreso_totales"`
}

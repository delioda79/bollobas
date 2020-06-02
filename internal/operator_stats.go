package internal

import (
	"context"
	"time"
)

// OperatorStatsRepository allows to manipulate the operator stats
type OperatorStatsRepository interface {
	GetAll(ctx context.Context, df DateFilter, pg Pagination) ([]OperatorStats, error)
	Add(ctx context.Context, trips *OperatorStats) error
}

// OperatorStats represents the stats for report 2
type OperatorStats struct {
	ID             int64     `json:"id"`
	Date           time.Time `transl:"date" json:"fecha_produccion"`
	OperatorID     *string   `transl:"operator_id" json:"id_operador"`
	Gender         *int      `transl:"gender" json:"genero"`
	CompletedTrips *int      `transl:"completed_trips" json:"cant_viajes"`
	DaysSince      *int      `transl:"days_since" json:"tiempo_registro"`
	AgeRange       *string   `transl:"age_range" json:"edad"`
	HoursConnected *string   `transl:"hours_connected" json:"horas_conectado"`
	TripHours      *string   `transl:"trip_hours" json:"horas_viaje"`
	TotRevenue     *string   `transl:"tot_revenue" json:"ingreso_totales"`
	DeletedAt      time.Time `transl:"deleted_at"`
}

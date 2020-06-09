package semovi

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/beatlabs/patron/component/async"
	"github.com/taxibeat/bollobas/internal"
)

type operatorStatsPayload struct {
	Date           int64   `json:"fecha_produccion"`
	OperatorID     *string `json:"id_operador"`
	Gender         *int    `json:"genero"`
	CompletedTrips *int    `json:"cant_viajes"`
	DaysSince      *int    `json:"tiempo_registro"`
	AgeRange       *string `json:"edad"`
	HoursConnected *string `json:"horas_conectado"`
	TripHours      *string `json:"horas_viaje"`
	TotRevenue     *string `json:"ingreso_totales"`
}

func (p operatorStatsPayload) toDomainModel() *internal.OperatorStats {
	return &internal.OperatorStats{
		Date:           time.Unix(p.Date/1000, 0),
		OperatorID:     p.OperatorID,
		Gender:         p.Gender,
		CompletedTrips: p.CompletedTrips,
		DaysSince:      p.DaysSince,
		AgeRange:       p.AgeRange,
		HoursConnected: p.HoursConnected,
		TripHours:      p.TripHours,
		TotRevenue:     p.TotRevenue,
	}
}

// OperatorStatsProcessor is responsible for consuming operator stats through kafka and inserting them to database
type OperatorStatsProcessor struct {
	store  internal.OperatorStatsRepository
	active bool
}

// NewOperatorStatsProcessor returns a newly created processor for operator stats events
func NewOperatorStatsProcessor(store internal.OperatorStatsRepository) *OperatorStatsProcessor {
	return &OperatorStatsProcessor{store, false}
}

// Process handles a given message
func (osp *OperatorStatsProcessor) Process(msg async.Message) error {
	if !osp.active {
		msg.Nack()
		return nil
	}

	// We need to decode the message again to reject if there are unknown fields.
	var payload operatorStatsPayload
	decoder := json.NewDecoder(bytes.NewReader(msg.Payload()))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil || (payload == operatorStatsPayload{}) {
		// If a key in the message is unknown we don't want to consume it again
		msg.Ack()
		return err
	}

	if err := osp.store.Add(msg.Context(), payload.toDomainModel()); err != nil {
		msg.Nack()
		return err
	}

	msg.Ack()
	return nil
}

// Activate turns the processing on or off depending on the switch given
func (osp *OperatorStatsProcessor) Activate(v bool) {
	osp.active = v
}

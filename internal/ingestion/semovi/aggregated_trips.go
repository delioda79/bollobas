package semovi

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/taxibeat/bollobas/internal"

	"github.com/beatlabs/patron/component/async"
)

type aggregatedTripsPayload struct {
	ID                     int64    `json:"id"`
	Date                   int64    `json:"fecha"`
	SupplierID             *string  `json:"id_proveedor"`
	TotalRides             *int     `json:"tot_viajes"`
	TotalVehicleRides      *int     `json:"tot_veh_viaje"`
	TotalAvailableVehicles *int     `json:"tot_veh_disp"`
	TotalDistTraveled      *float64 `json:"dist_pasajero"`
	PassingTime            *float64 `json:"tiempo_pasajero"`
	RequestTime            *float64 `json:"tiempo_solicitud"`
	EmptyTime              *float64 `json:"tiempo_vacio"`
	EodMultiplier          *float64 `json:"multiplicador_eod"`
	Accessibility          *float64 `json:"accesibilidad"`
	FemaleOperator         *float64 `json:"operador_mujer"`
	EodStart               *int     `json:"inicio_eod"`
	EodEnd                 *int     `json:"fin_eod"`
	EodPassDist            *float64 `json:"dist_pasajero_eod"`
	EodPassTime            *int     `json:"tiempo_pasajero_eod"`
	RequestDist            *float64 `json:"dist_solicitud"`
	EmptyDist              *float64 `json:"dist_vacío"`
	EodRequestDist         *float64 `json:"dist_solicitud_eod"`
	EodRequestTime         *float64 `json:"tiempo_solicitud_eod"`
	EodEmptyDist           *float64 `json:"dist_vacio_eod"`
	EodEmptyTime           *float64 `json:"tiempo_vacio_eod"`
}

func (p aggregatedTripsPayload) toDomainModel() *internal.AggregatedTrips {
	return &internal.AggregatedTrips{
		ID:                     p.ID,
		Date:                   time.Unix(p.Date/1000, 0),
		SupplierID:             p.SupplierID,
		TotalRides:             p.TotalRides,
		TotalVehicleRides:      p.TotalVehicleRides,
		TotalAvailableVehicles: p.TotalAvailableVehicles,
		TotalDistTraveled:      p.TotalDistTraveled,
		PassingTime:            p.PassingTime,
		RequestTime:            p.RequestTime,
		EmptyTime:              p.EmptyTime,
		EodMultiplier:          p.EodMultiplier,
		Accessibility:          p.Accessibility,
		FemaleOperator:         p.FemaleOperator,
		EodStart:               p.EodStart,
		EodEnd:                 p.EodEnd,
		EodPassDist:            p.EodPassDist,
		EodPassTime:            p.EodPassTime,
		RequestDist:            p.RequestDist,
		EmptyDist:              p.EmptyDist,
		EodRequestDist:         p.EodRequestDist,
		EodRequestTime:         p.EodRequestTime,
		EodEmptyDist:           p.EodEmptyDist,
		EodEmptyTime:           p.EodEmptyTime,
	}
}

// AggregatedTripsProcessor is responsible for consuming aggregated trips through kafka and inserting them to database
type AggregatedTripsProcessor struct {
	store  internal.AggregatedTripsRepository
	active bool
}

// NewAggregatedTripsProcessor returns a newly created processor for aggregated trips events
func NewAggregatedTripsProcessor(store internal.AggregatedTripsRepository) *AggregatedTripsProcessor {
	return &AggregatedTripsProcessor{store, false}
}

// Process handles a given message
func (atp *AggregatedTripsProcessor) Process(msg async.Message) error {
	if !atp.active {
		msg.Nack()
		return nil
	}

	// We need to decode the message again to reject if there are unknown fields.
	var payload aggregatedTripsPayload
	decoder := json.NewDecoder(bytes.NewReader(msg.Payload()))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil || (payload == aggregatedTripsPayload{}) {
		// If a key in the message is unknown we don't want to consume it again
		msg.Ack()
		return err
	}

	if err := atp.store.Add(msg.Context(), payload.toDomainModel()); err != nil {
		msg.Nack()
		return err
	}

	msg.Ack()
	return nil
}

// Activate turns the processing on or off depending on the switch given
func (atp *AggregatedTripsProcessor) Activate(v bool) {
	atp.active = v
}

package semovi

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/taxibeat/bollobas/internal"

	"github.com/beatlabs/patron/component/async"
)

type trafficIncidentsPayload struct {
	Date           int64   `json:"tiempo_hecho"`
	Type           *int    `json:"hecho_trans"`
	Plates         *string `json:"placa"`
	Licence        *string `json:"licencia"`
	TravelDistance *string `json:"distancia_viaje"`
	TravelTime     *string `json:"tiempo_viaje"`
	Coordinates    *string `json:"ubicación"`
}

func (p trafficIncidentsPayload) toDomainModel() *internal.TrafficIncident {
	return &internal.TrafficIncident{
		Date:           time.Unix(p.Date/1000, 0),
		Type:           p.Type,
		Plates:         p.Plates,
		Licence:        p.Licence,
		TravelDistance: p.TravelDistance,
		TravelTime:     p.TravelTime,
		Coordinates:    p.Coordinates,
	}
}

// TrafficIncidentsProcessor is responsible for consuming traffic incidents through kafka and inserting them to database
type TrafficIncidentsProcessor struct {
	store  internal.TrafficIncidentsRepository
	active bool
}

// NewTrafficIncidentsProcessor returns a newly created processor for traffic incidents events
func NewTrafficIncidentsProcessor(store internal.TrafficIncidentsRepository) *TrafficIncidentsProcessor {
	return &TrafficIncidentsProcessor{store, false}
}

// Process handles a given message
func (tip *TrafficIncidentsProcessor) Process(msg async.Message) error {
	if !tip.active {
		msg.Nack()
		return nil
	}

	kfkPl := make(map[string]interface{})
	if err := msg.Decode(&kfkPl); err != nil {
		// If the message is faulty we don't want to consume it again
		msg.Ack()
		return err
	}

	// We need to decode the message again to reject if there are unknown fields.
	// TODO: Refactor this to avoid decoding again
	var payload trafficIncidentsPayload
	vp, err := json.Marshal(kfkPl)
	if err != nil {
		msg.Ack()
		return err
	}
	decoder := json.NewDecoder(strings.NewReader(string(vp)))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&payload)
	if err != nil {
		// If a key in the message is unknown we don't want to consume it again
		msg.Ack()
		return err
	}

	if err := tip.store.Add(msg.Context(), payload.toDomainModel()); err != nil {
		msg.Nack()
		return err
	}

	msg.Ack()
	return nil
}

// Activate turns the processing on or off depending on the switch given
func (tip *TrafficIncidentsProcessor) Activate(v bool) {
	tip.active = v
}

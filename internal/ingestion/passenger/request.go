package passenger

import (
	"github.com/taxibeat/bollobas/internal/mixpanel"
	"nanomsg.org/go/mangos/v2"
	"time"

	"github.com/beatlabs/patron/log"
	"github.com/taxibeat/bollobas/internal/ingestion"
	"github.com/taxibeat/bollobas/internal/mixpanel/pkg/parseid"

	"github.com/beatlabs/patron/component/async"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/pkg/errors"
)

//RequestProcessor processes the messages from request topics and forwards a request message
type RequestProcessor struct {
	mangos.Socket
	active   bool
	provider string
	topic    string
}

// Process is part of the patron interface and processes incoming messages
func (kc *RequestProcessor) Process(msg async.Message) error {
	start := time.Now()

	if !kc.active {
		ingestion.ObserveCount(kc.provider, kc.topic, false, true)
		return nil
	}
	rr := RequestRide{}

	err := msg.Decode(&rr)
	if err != nil {
		ingestion.ObserveCount(kc.provider, kc.topic, true, false)
		return errors.Errorf("failed to unmarshal ride cancellation %v", err)
	}

	ingestion.ObserveCount(kc.provider, kc.topic, true, true)
	return kc.publish(rr, start)
}

func (kc *RequestProcessor) publish(cr RequestRide, start time.Time) error {
	passengerID, err := parseid.EncryptString(cr.Passenger.ID, "pa")
	if err != nil {
		log.Errorf("error encrypting passenger ID: %v", err)
		return err
	}

	idt := mixpanel.RideRequest{
		UserID:   passengerID,
		RquestID: cr.RequestID,
	}

	bts, err := json.Encode(idt)
	if err != nil {
		return err
	}

	log.Debugf("Sending ride request %+v", idt)

	err = kc.Send(bts)
	if err != nil {
		return err
	}

	ingestion.ObserveRepublishedCount("request", "ride_request")
	ingestion.ObserveLatency(kc.provider, kc.topic, time.Since(start))

	return nil
}

// NewRequestProcessor instantiates a new component
func NewRequestProcessor(url, provider, topic string) (*RequestProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	return &RequestProcessor{Socket: sock, active: false, provider: provider, topic: topic}, nil
}

// Activate will activate the processor
func (kc *RequestProcessor) Activate(v bool) {
	kc.active = v
}

// RequestRide represents a ride request message coming from kafka
type RequestRide struct {
	Passenger RequestPassenger `json:"passenger"`
	RequestID int              `json:"id_request"`
}

//RequestPassenger Struct
type RequestPassenger struct {
	ID int `json:"id_passenger"`
}

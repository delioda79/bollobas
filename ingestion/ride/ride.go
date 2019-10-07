package ride

import (
	"bollobas"
	"bollobas/ingestion"
	"bollobas/pkg/parseid"
	"time"

	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/log"
	"github.com/pkg/errors"
	"nanomsg.org/go/mangos/v2"
	_ "nanomsg.org/go/mangos/v2/transport/all" // import transports (important)
)

//Processor processes the messages from ride topics and forwards a ride confirmation message
type Processor struct {
	mangos.Socket
	active bool
	provider string
	topic string
}

// Process is part of the patron interface and processes incoming messages
func (kc *Processor) Process(msg async.Message) error {
	start := time.Now()

	if !kc.active {
		ingestion.ObserveCount(kc.provider, kc.topic, false, true)
		return nil
	}

	cr := Ride{}

	err := msg.Decode(&cr)
	if err != nil {
		ingestion.ObserveCount(kc.provider, kc.topic, true, false)
		return errors.Errorf("failed to unmarshal ride cancellation %v", err)
	}

	ingestion.ObserveCount(kc.provider, kc.topic, true, true)
	return kc.publish(cr, start)
}

func (kc *Processor) publish(cr Ride, start time.Time) error {

	log.Debugf("Ride received: %+v | events: %+v", cr, cr.Events)
	if len(cr.Events) == 0 {
		log.Debugf("Ready to send")
		idt := bollobas.RideRequestConfirmed{
			UserID:   parseid.EncryptString(cr.Passenger.PassengerID, "pa"),
			RquestID: cr.RequestID,
			Date:     time.Unix(cr.Created, 0),
		}

		bts, err := json.Encode(idt)
		if err != nil {
			return errors.Errorf("Error when decoding the event %v", err)
		}
		log.Debugf("Sending ride confirmation: %+v", idt)

		err = kc.Send(bts)
		if err != nil {
			return errors.Errorf("Error when sending the event: %v", err)
		}

		ingestion.ObserveLatency(kc.provider, kc.topic, time.Since(start))
		ingestion.ObserveRepublishedCount("confirmation", "ride_request")
	}

	return nil
}

// NewRideProcessor instantiates a new processor
func NewRideProcessor(url string, provider, topic string) (*Processor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	rp := &Processor{Socket: sock, active: false, provider:provider, topic:topic}

	return rp, nil
}

// Activate will activate the processor
func (kc *Processor) Activate(v bool) {
	kc.active = v
}

// Event represents a message for a ride coming from kafka
type Event struct {
	Who  string
	What string
	When int64
	Key  string
}

//Passenger represents passenger infos nested inside the ride message from kafka
type Passenger struct {
	PassengerID int `json:"id_passenger"`
}

//Ride ...
type Ride struct {
	Passenger Passenger
	Events    []Event
	RequestID int `json:"id_request"`
	Duration  interface{}
	Created   int64 `json:"ride_created_at"`
}

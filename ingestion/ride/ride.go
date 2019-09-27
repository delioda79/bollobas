package ride

import (
	"bollobas"
	"bollobas/ingestion"
	"bollobas/pkg/parseid"
	"fmt"
	"time"

	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/log"
	"github.com/pkg/errors"
	"nanomsg.org/go/mangos/v2"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

type RideProcessor struct {
	mangos.Socket
	active bool
}

// Process is part of the patron interface and processes incoming messages
func (kc *RideProcessor) Process(msg async.Message) error {
	if !kc.active {
		return nil
	}

	fmt.Println("RIDE ACTIVE!")
	cr := Ride{}

	err := msg.Decode(&cr)
	if err != nil {
		return errors.Errorf("failed to unmarshal ride cancellation %v", err)
	}

	return kc.publish(cr)
}

func (kc *RideProcessor) publish(cr Ride) error {

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
		log.Debugf("Sending: %+v", idt)
		err = kc.Send(bts)
		if err != nil {
			return errors.Errorf("Error when sending the event: %v", err)
		}
	}

	return nil
}

// NewCancellationProcessor instantiates a new component
func NewRideProcessor(url string) (*RideProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	rp := &RideProcessor{Socket: sock, active: false}

	return rp, nil
}

// Activate will activate the processor
func (kc *RideProcessor) Activate(v bool) {
	kc.active = v
}

// RideEvent represents a message for a ride coming from kafka
type RideEvent struct {
	Who  string
	What string
	When int64
	Key  string
}

type Passenger struct {
	PassengerID int `json:"id_passenger"`
}

type Ride struct {
	Passenger Passenger
	Events    []RideEvent
	RequestID int `json:"id_request"`
	Duration  interface{}
	Created   int64 `json:"ride_created_at"`
}

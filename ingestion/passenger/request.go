package passenger

import (
	"bollobas"
	"bollobas/ingestion"
	"bollobas/pkg/parseid"

	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/pkg/errors"
	"nanomsg.org/go/mangos/v2"
)

//RequestProcessor processes the messages from request topics and forwards a request message
type RequestProcessor struct {
	mangos.Socket
	active bool
}

// Process is part of the patron interface and processes incoming messages
func (kc *RequestProcessor) Process(msg async.Message) error {
	if !kc.active {
		return nil
	}
	rr := RequestRide{}

	err := msg.Decode(&rr)
	if err != nil {
		return errors.Errorf("failed to unmarshal ride cancellation %v", err)
	}

	return kc.publish(rr)
}

func (kc *RequestProcessor) publish(cr RequestRide) error {

	idt := bollobas.RideRequest{
		UserID:   parseid.EncryptString(cr.Passenger.ID, "pa"),
		RquestID: cr.RequestID,
	}

	bts, err := json.Encode(idt)
	if err != nil {
		return err
	}

	return kc.Send(bts)
}

// NewRequestProcessor instantiates a new component
func NewRequestProcessor(url string) (*RequestProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	return &RequestProcessor{Socket: sock, active: false}, nil
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

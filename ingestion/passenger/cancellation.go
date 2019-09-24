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

type CancellationProcessor struct {
	mangos.Socket
	active bool
}

// Process is part of the patron interface and processes incoming messages
func (kc *CancellationProcessor) Process(msg async.Message) error {
	if !kc.active {
		return nil
	}
	cr := CancelRideRequest{}

	err := msg.Decode(&cr)
	if err != nil {
		return errors.Errorf("failed to unmarshal ride cancellation %v", err)
	}

	return kc.publish(cr)
}

func (kc *CancellationProcessor) publish(cr CancelRideRequest) error {

	idt := bollobas.RideRequest{
		UserID:   parseid.EncryptString(cr.PassengerID, "pa"),
		RquestID: cr.RequestID,
	}

	bts, err := json.Encode(idt)
	if err != nil {
		return err
	}

	return kc.Send(bts)
}

// NewCancellationProcessor instantiates a new component
func NewCancellationProcessor(url string) (*CancellationProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	return &CancellationProcessor{Socket: sock, active: false}, nil
}

// Activate will activate the processor
func (kc *CancellationProcessor) Activate(v bool) {
	kc.active = v
}

// CancelRideRequest represents a message for a ride request cancellation coming from kafka
type CancelRideRequest struct {
	PassengerID int `json:"id_passenger"`
	RequestID   int `json:"id_request"`
}

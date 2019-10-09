package passenger

import (
	"bollobas"
	"bollobas/ingestion"
	"bollobas/pkg/parseid"
	"github.com/beatlabs/patron/log"
	"time"

	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/pkg/errors"
	"nanomsg.org/go/mangos/v2"
)

//CancellationProcessor processes the messages from request_cancel topics and forwards a cancel ride message
type CancellationProcessor struct {
	mangos.Socket
	active   bool
	topic    string
	provider string
}

// Process is part of the patron interface and processes incoming messages
func (kc *CancellationProcessor) Process(msg async.Message) error {
	start := time.Now()

	if !kc.active {
		ingestion.ObserveCount(kc.provider, kc.topic, false, true)
		return nil
	}

	cr := CancelRideRequest{}

	err := msg.Decode(&cr)
	if err != nil {
		ingestion.ObserveCount(kc.provider, kc.topic, true, false)
		return errors.Errorf("failed to unmarshal ride cancellation %v", err)
	}

	ingestion.ObserveCount(kc.provider, kc.topic, true, true)
	return kc.publish(cr, start)
}

func (kc *CancellationProcessor) publish(cr CancelRideRequest, start time.Time) error {

	idt := bollobas.RideRequest{
		UserID:   parseid.EncryptString(cr.PassengerID, "pa"),
		RquestID: cr.RequestID,
	}

	bts, err := json.Encode(idt)
	if err != nil {
		return err
	}

	log.Debugf("Sending request cancellation %+v", idt)

	err = kc.Send(bts)
	if err != nil {
		return err
	}

	ingestion.ObserveRepublishedCount("cancellation", "ride_request")
	ingestion.ObserveLatency(kc.provider, kc.topic, time.Since(start))

	return nil
}

// NewCancellationProcessor instantiates a new component
func NewCancellationProcessor(url, provider, topic string) (*CancellationProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	return &CancellationProcessor{Socket: sock, active: false, provider: provider, topic: topic}, nil
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

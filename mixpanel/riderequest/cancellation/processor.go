package cancellation

import (
	"encoding/json"

	"github.com/taxibeat/bollobas"

	"github.com/beatlabs/patron/errors"
	"github.com/beatlabs/patron/log"
	"github.com/dukex/mixpanel"
	_ "nanomsg.org/go/mangos/v2/transport/all" //import
)

// Processor subscribes to messages sent by any registered publisher in the internal registry
type Processor struct {
	mixpanel.Mixpanel
}

// Process starts the go routine which will receive the messages
func (hdl *Processor) Process(msg []byte) error {

	idt := &bollobas.RideRequestCancellation{}
	err := json.Unmarshal(msg, idt)
	if err != nil {
		return errors.Errorf("error unmarshaling the data: %v", err)
	}
	log.Debugf("Ride canceled: %v", string(msg))
	return hdl.incrementRideCancellations(idt)
}

func (hdl *Processor) incrementRideCancellations(idt *bollobas.RideRequestCancellation) error {
	//id := idt.ID
	prps := &RideRequestCancellation{Cancel: 1}

	bts, err := json.Marshal(prps)
	if err != nil {
		return errors.Errorf("Impossible to unmarshal: %v", err)
	}

	mp := map[string]interface{}{}

	err = json.Unmarshal(bts, &mp)
	if err != nil {
		return errors.Errorf("error while unmarshaling the CancellationRequest: %v", err)
	}

	err = hdl.Update(idt.UserID, &mixpanel.Update{Properties: mp, Operation: "$add"})
	if err != nil {
		return errors.Errorf("error while updating the CancellationRequest: %v", err)
	}

	return nil
}

// Topic returns the topic
func (hdl *Processor) Topic() string {
	return "ride_request_cancellation"
}

// RideRequestCancellation represents a message notifying of a ride cancellation by a passenger
type RideRequestCancellation struct {
	Cancel int `json:"CancelledRequests"`
}

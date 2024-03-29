package confirmation

import (
	"encoding/json"
	"fmt"

	bollobas "github.com/taxibeat/bollobas/internal/mixpanel"

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

	idt := &bollobas.RideRequestConfirmed{}
	err := json.Unmarshal(msg, idt)
	if err != nil {
		return fmt.Errorf("error unmarshaling the data: %v", err)
	}

	log.Debugf("Request confirmed: %v", string(msg))

	return hdl.incrementRideCconfirmations(idt)
}

func (hdl *Processor) incrementRideCconfirmations(idt *bollobas.RideRequestConfirmed) error {
	log.Debugf("Sending Confirm Ride Request for UserID: %s", idt.UserID)

	err := hdl.Update(idt.UserID, &mixpanel.Update{Properties: map[string]interface{}{"ConfirmedRides": 1}, Operation: "$add"})
	if err != nil {
		return fmt.Errorf("error while updating the ConfirmationRequest: %v", err)
	}

	err = hdl.Update(idt.UserID, &mixpanel.Update{Properties: map[string]interface{}{"LastRide": idt.Date}, Operation: "$set"})
	if err != nil {
		return fmt.Errorf("error while updating the ConfirmationRequest: %v", err)
	}

	return nil
}

// Topic returns the topic
func (hdl *Processor) Topic() string {
	return "ride_request_confirmation"
}

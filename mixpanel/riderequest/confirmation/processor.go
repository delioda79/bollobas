package confirmation

import (
	"bollobas"
	"encoding/json"

	"github.com/beatlabs/patron/errors"
	"github.com/beatlabs/patron/log"
	"github.com/dukex/mixpanel"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
)

// Processor subscribes to messages sent by any registered publisher in the internal registry
type Processor struct {
	mixpanel.Mixpanel
}

// Run starts the go routine which will receive the messages
func (hdl *Processor) Process(msg []byte) error {

	idt := &bollobas.RideRequestConfirmed{}
	err := json.Unmarshal(msg, idt)
	if err != nil {
		return errors.Errorf("error unmarshaling the data: %v", err)
	}

	log.Debugf("Request confirmed: %v", string(msg))

	return hdl.incrementRideCconfirmations(idt)
}

func (hdl *Processor) incrementRideCconfirmations(idt *bollobas.RideRequestConfirmed) error {
	log.Debugf("Sending Confirm Ride Request for UserID: %s", idt.UserID)

	err := hdl.Update(idt.UserID, &mixpanel.Update{Properties: map[string]interface{}{"ConfirmedRides": 1}, Operation: "$add"})
	if err != nil {
		return errors.Errorf("error while updating the ConfirmationRequest: %v", err)
	}

	err = hdl.Update(idt.UserID, &mixpanel.Update{Properties: map[string]interface{}{"LastRide": idt.Date}, Operation: "$set"})
	if err != nil {
		return errors.Errorf("error while updating the ConfirmationRequest: %v", err)
	}

	return nil
}

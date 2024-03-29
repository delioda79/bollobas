package riderequest

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
func (p *Processor) Process(msg []byte) error {
	idt := &bollobas.RideRequest{}
	err := json.Unmarshal(msg, idt)
	if err != nil {
		return fmt.Errorf("error unmarshaling the data: %v", err)
	}

	return p.incrementRideRequests(idt)

}

func (p *Processor) incrementRideRequests(idt *bollobas.RideRequest) error {
	//id := idt.ID
	prps := &RideRequest{Request: 1}

	bts, err := json.Marshal(prps)
	if err != nil {
		return fmt.Errorf("Impossible to unmarshal: %v", err)
	}

	mp := map[string]interface{}{}

	err = json.Unmarshal(bts, &mp)
	if err != nil {
		return fmt.Errorf("error while unmarshaling the RideRequest: %v", err)
	}

	err = p.Update(idt.UserID, &mixpanel.Update{Properties: mp, Operation: "$add"})
	if err != nil {
		return fmt.Errorf("error while updating the RideRequest: %v", err)
	}
	log.Debugf("Sent to mixpanel %v %v", idt.UserID, mp)
	return nil
}

// Topic returns the topic
func (p *Processor) Topic() string {
	return "ride_request"
}

// RideRequest represents a ride request message
type RideRequest struct {
	Request int `json:"RequestedRides"`
}

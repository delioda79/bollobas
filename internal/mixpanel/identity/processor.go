package identity

import (
	"encoding/json"
	"fmt"
	"time"

	bollobas "github.com/taxibeat/bollobas/internal/mixpanel"

	"github.com/beatlabs/patron/log"

	"github.com/dukex/mixpanel"
	_ "nanomsg.org/go/mangos/v2/transport/all" //import
)

// Processor subscribes to messages sent by any registered publisher in the internal registry
type Processor struct {
	mixpanel.Mixpanel
}

// Process run starts the go routine which will receive the messages
func (p *Processor) Process(msg []byte) error {

	idt := &bollobas.Identity{}
	err := json.Unmarshal(msg, idt)
	if err != nil {
		return fmt.Errorf("error unmarshaling the data: %v", err)
	}
	log.Debugf("SENDING TO MIXPANEL %v", string(msg))
	return p.updateIdentity(idt)
}

func (p *Processor) updateIdentity(idt *bollobas.Identity) error {
	//id := idt.ID
	prps := &Identity{
		FirstName:        idt.FirstName,
		LastName:         idt.LastName,
		RegistrationDate: idt.RegistrationDate,
		ReferralCode:     idt.ReferralCode,
		Type:             idt.Type,
		Email:            idt.Email,
		Phone:            idt.Phone,
		Market:           idt.Market,
	}

	bts, err := json.Marshal(prps)
	if err != nil {
		return fmt.Errorf("Impossible to unmarshal: %v", err)
	}

	mp := map[string]interface{}{}

	err = json.Unmarshal(bts, &mp)
	if err != nil {
		return fmt.Errorf("error while unmarshaling the identity: %v", err)
	}

	err = p.Update(idt.ID, &mixpanel.Update{Properties: mp, Operation: "$set"})
	if err != nil {
		return fmt.Errorf("error while updating the identity: %v", err)
	}

	return nil
}

// Topic returns the topic
func (p *Processor) Topic() string {
	return "identity"
}

//Identity represents a mixpanel identity
type Identity struct {
	FirstName        string    `json:"$first_name,omitempty"`
	LastName         string    `json:"$last_name,omitempty"`
	RegistrationDate time.Time `json:"$created,omitempty"`
	Type             string    `json:"type,omitempty"`
	Email            string    `json:"$email,omitempty"`
	ReferralCode     string    `json:"referral_code,omitempty"`
	Phone            string    `json:"$phone,omitempty"`
	Market           string    `json:"market,omitempty"`
}

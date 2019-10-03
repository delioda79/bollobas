package passenger

import (
	"bollobas"
	"bollobas/ingestion"
	"fmt"
	"github.com/prometheus/common/log"
	"time"

	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/errors"
	"nanomsg.org/go/mangos/v2"
	_ "nanomsg.org/go/mangos/v2/transport/all" //import
)

//AccountProcessor processes the messages from passenger_analytics topics and forwards an account message
type AccountProcessor struct {
	mangos.Socket
	active bool
}

// Process is part of the patron interface and processes incoming messages
func (kc *AccountProcessor) Process(msg async.Message) error {
	if !kc.active {
		return nil
	}
	passenger := Passenger{}

	err := msg.Decode(&passenger)
	if err != nil {
		return errors.Errorf("failed to unmarshal passenger %v", err)
	}

	return kc.publish(passenger)
}

// Activate will activate the processor
func (kc *AccountProcessor) Activate(v bool) {
	kc.active = v
}

func (kc *AccountProcessor) publish(passenger Passenger) error {

	idt := bollobas.Identity{
		ID:               passenger.ID,
		FirstName:        passenger.FirstName,
		LastName:         passenger.LastName,
		RegistrationDate: passenger.RegistrationDate,
		Phone:            fmt.Sprintf("+%s%s", passenger.PhonePrefix, passenger.PhoneNo),
		Type:             "passenger",
		Email:            passenger.Email,
	}

	log.Debugf("Sending passenger %+v", idt)

	bts, err := json.Encode(idt)
	if err != nil {
		return err
	}

	return kc.Send(bts)
}

// NewAccountProcessor instantiates a new processor
func NewAccountProcessor(url string) (*AccountProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	return &AccountProcessor{Socket: sock, active: false}, nil
}

// Passenger represents a passenger message coming from kafka
type Passenger struct {
	ID               string `json:"passenger_id"`
	Email            string
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	PhoneNo          string    `json:"phone"`
	PhonePrefix      string    `json:"phone_prefix"`
	RegistrationDate time.Time `json:"registration_date"`
	Action           string    `json:"event_action"`
}

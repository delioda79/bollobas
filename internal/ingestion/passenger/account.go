package passenger

import (
	"fmt"
	"github.com/taxibeat/bollobas/internal/mixpanel"
	"nanomsg.org/go/mangos/v2"

	"github.com/beatlabs/patron/log"
	"github.com/taxibeat/bollobas/internal/ingestion"

	"time"

	"github.com/beatlabs/patron/component/async"
	"github.com/beatlabs/patron/encoding/json"
	_ "nanomsg.org/go/mangos/v2/transport/all" //import
)

//AccountProcessor processes the messages from passenger_analytics topics and forwards an account message
type AccountProcessor struct {
	mangos.Socket
	active   bool
	provider string
	topic    string
	market   string
}

// Process is part of the patron interface and processes incoming messages
func (kc *AccountProcessor) Process(msg async.Message) error {
	start := time.Now()

	if !kc.active {
		ingestion.ObserveCount(kc.provider, kc.topic, false, true)
		return nil
	}
	passenger := Passenger{}

	err := msg.Decode(&passenger)
	if err != nil {
		return fmt.Errorf("failed to unmarshal passenger %v", err)
	}

	ingestion.ObserveCount(kc.provider, kc.topic, true, true)
	return kc.publish(passenger, start)
}

// Activate will activate the processor
func (kc *AccountProcessor) Activate(v bool) {
	kc.active = v
}

func (kc *AccountProcessor) publish(passenger Passenger, start time.Time) error {

	idt := mixpanel.Identity{
		ID:               passenger.ID,
		FirstName:        passenger.FirstName,
		LastName:         passenger.LastName,
		RegistrationDate: passenger.RegistrationDate,
		Phone:            fmt.Sprintf("+%s%s", passenger.PhonePrefix, passenger.PhoneNo),
		Type:             "passenger",
		Email:            passenger.Email,
		Market:           kc.market,
	}

	log.Debugf("Sending passenger %+v", idt)

	bts, err := json.Encode(idt)
	if err != nil {
		return err
	}

	err = kc.Send(bts)
	if err != nil {
		return err
	}

	ingestion.ObserveRepublishedCount("identity", "passenger")
	ingestion.ObserveLatency(kc.provider, kc.topic, time.Since(start))

	return nil
}

// NewAccountProcessor instantiates a new processor
func NewAccountProcessor(url, provider, topic, market string) (*AccountProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	return &AccountProcessor{Socket: sock, active: false, provider: provider, topic: topic, market: market}, nil
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

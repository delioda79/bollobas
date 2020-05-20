package driver

import (
	"fmt"
	"github.com/taxibeat/bollobas/internal/mixpanel"
	"nanomsg.org/go/mangos/v2"
	"time"

	"github.com/taxibeat/bollobas/internal/ingestion"

	"github.com/beatlabs/patron/component/async"
	"github.com/beatlabs/patron/encoding/json"
	_ "nanomsg.org/go/mangos/v2/transport/all" //import
)

// AccountProcessor processes the messages from driver_analytics topics and forwards an account message
type AccountProcessor struct {
	mangos.Socket
	active bool
	market string
}

// Process is part of the patron interface and processes incoming messages
func (kc *AccountProcessor) Process(msg async.Message) error {
	if !kc.active {
		return nil
	}
	driver := Driver{}

	err := msg.Decode(&driver)
	if err != nil {
		return fmt.Errorf("failed to unmarshal driver %v", err)
	}

	return kc.publish(driver)
}

func (kc *AccountProcessor) publish(driver Driver) error {

	idt := mixpanel.Identity{
		ID:               driver.ID,
		FirstName:        driver.FirstName,
		LastName:         driver.LastName,
		RegistrationDate: driver.RegistrationDate,
		ReferralCode:     driver.ReferralCode,
		Phone:            fmt.Sprintf("+%s%s%s", driver.AreaPrefix, driver.PhonePrefix, driver.PhoneNo),
		Type:             "driver",
		Email:            driver.Email,
		Market:           kc.market,
	}

	bts, err := json.Encode(idt)
	if err != nil {
		return err
	}

	return kc.Send(bts)
}

// NewAccountProcessor instantiates a new processor
func NewAccountProcessor(url, market string) (*AccountProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	return &AccountProcessor{Socket: sock, active: false, market: market}, nil
}

// Activate will activate the processor
func (kc *AccountProcessor) Activate(v bool) {
	kc.active = v
}

// Driver represents a driver message coming from kafka
type Driver struct {
	ID               string `json:"driver_id"`
	Email            string
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	PhonePrefix      string    `json:"phone_prefix"`
	AreaPrefix       string    `json:"area_prefix"`
	PhoneNo          string    `json:"phone"`
	ReferralCode     string    `json:"registration_id_reference"`
	RegistrationDate time.Time `json:"registration_date"`
	Action           string    `json:"event_action"`
}

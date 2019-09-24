package driver

import (
	"bollobas"
	"bollobas/ingestion"
	"fmt"
	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/errors"
	"nanomsg.org/go/mangos/v2"
	_ "nanomsg.org/go/mangos/v2/transport/all"
	"time"
)

type AccountProcessor struct {
	mangos.Socket
	active bool
}

// Process is part of the patron interface and processes incoming messages
func (kc *AccountProcessor) Process(msg async.Message) error {
	if !kc.active {
		return nil
	}
	driver := Driver{}

	err := msg.Decode(&driver)
	if err != nil {
		return errors.Errorf("failed to unmarshal driver %v", err)
	}

	return kc.publish(driver)
}

func (kc *AccountProcessor) publish(driver Driver) error {

	idt := bollobas.Identity{
		ID:               driver.ID,
		FirstName:        driver.FirstName,
		LastName:         driver.LastName,
		RegistrationDate: driver.RegistrationDate,
		ReferralCode:     driver.ReferralCode,
		Phone:            fmt.Sprintf("%s %s %s", driver.AreaPrefix, driver.PhonePrefix, driver.PhoneNo),
		Type:             "driver",
		Email:            driver.Email,
	}

	bts, err := json.Encode(idt)
	if err != nil {
		return err
	}

	return kc.Send(bts)
}

// NewAccountProcessor instantiates a new component
func NewAccountProcessor(url string) (*AccountProcessor, error) {

	sock, err := ingestion.NewPublisher([]string{url})
	if err != nil {
		return nil, err
	}
	return &AccountProcessor{Socket: sock, active: false}, nil
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

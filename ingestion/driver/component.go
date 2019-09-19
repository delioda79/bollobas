package driver

import (
	"bollobas"
	"fmt"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/async/kafka"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/errors"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"time"

	"nanomsg.org/go/mangos/v2"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
)

// KafkaComponent is a receiver for a specific kafka topic which will then forward the message as an identity
type KafkaComponent struct {
	patron.Component
	mangos.Socket
}

// Process is part of the patron interface and processes incoming messages
func (kc *KafkaComponent) Process(msg async.Message) error {
	driver := Driver{}
	err := msg.Decode(&driver)
	if err != nil {
		return errors.Errorf("failed to unmarshal driver %v", err)
	}

	return kc.publish(driver)

}

func (kc *KafkaComponent) publish(driver Driver) error {

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

// NewKafkaComponent instantiates a new component
func NewKafkaComponent(name, broker, topic, group, url string) (*KafkaComponent, error) {
	var sock mangos.Socket
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		return nil, errors.Errorf("can't get new pub socket: %s", err)
	}
	if err = sock.Listen(url); err != nil {
		return nil, errors.Errorf("can't listen on pub socket: %s", err.Error())
	}

	kafkaCmp := KafkaComponent{}

	cf, err := kafka.New(name, json.Type, topic, group, []string{broker})
	if err != nil {
		return nil, err
	}

	cmp, err := async.New("driver-kafka-cmp", kafkaCmp.Process, cf, async.ConsumerRetry(10, 5*time.Second))
	if err != nil {
		return nil, err
	}

	kafkaCmp.Component = cmp
	kafkaCmp.Socket = sock

	return &kafkaCmp, nil
}

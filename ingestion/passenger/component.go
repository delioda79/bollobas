package passenger

import (
	"bollobas"
	"fmt"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/async/kafka"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/errors"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"time"

	_ "nanomsg.org/go/mangos/v2/transport/inproc"
)

// KafkaComponent is a receiver for a specific kafka topic which will then forward the message as an identity
type KafkaComponent struct {
	patron.Component
	mangos.Socket
}

// Process is part of the patron interface and processes incoming messages
func (kc *KafkaComponent) Process(msg async.Message) error {

	passenger := Passenger{}

	err := msg.Decode(&passenger)
	if err != nil {
		return errors.Errorf("failed to unmarshal passenger %v", err)
	}

	return kc.publish(passenger)
}

func (kc *KafkaComponent) publish(passenger Passenger) error {

	idt := bollobas.Identity{
		ID:               passenger.ID,
		FirstName:        passenger.FirstName,
		LastName:         passenger.LastName,
		RegistrationDate: passenger.RegistrationDate,
		Phone:            fmt.Sprintf("%s %s", passenger.PhonePrefix, passenger.PhoneNo),
		Type:             "passenger",
		Email:            passenger.Email,
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

package passenger

import (
	"fmt"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/async/kafka"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/errors"
	"time"
)

type KafkaComponent struct {
	patron.Component
}

func (kc *KafkaComponent) Process(msg async.Message) error {
	
	passenger := Passenger{}
	err := msg.Decode(&passenger)
	if err != nil {
		return errors.Errorf("failed to unmarshal passenger %v", err)
	}

	fmt.Printf("%+v\n", passenger)

	//dt := time.Unix(int64(passenger.RegistrationDate), 0)
	//fmt.Println("REGDATE", dt)
	return nil
}

func NewKafkaComponent(name, broker, topic, group string) (*KafkaComponent, error) {

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

	return &kafkaCmp, nil
}

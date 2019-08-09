package driver

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

	driver := Driver{}
	err := msg.Decode(&driver)
	if err != nil {
		return errors.Errorf("failed to unmarshal driver %v", err)
	}

	fmt.Printf("%+v\n",driver)
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

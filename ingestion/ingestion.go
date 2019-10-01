package ingestion

import (
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/async/kafka"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/errors"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
)

// Processor is an interface for structs able to process kafka messages
type Processor interface {
	mangos.Socket
	Process(msg async.Message) error
	Activate(v bool)
}

// NewPublisher returns a pub socket listening to url
func NewPublisher(urls []string) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		return nil, errors.Errorf("can't get new pub socket: %s", err)
	}

	for _, url := range urls {
		if err = sock.Listen(url); err != nil {
			return nil, errors.Errorf("can't listen on pub socket: %s", err.Error())
		}
	}
	return sock, nil
}

// KafkaComponent is a patron component for receiving and processing kafka messages
type KafkaComponent struct {
	patron.Component
}

// NewKafkaComponent instantiates a new component
func NewKafkaComponent(name, broker, topic, group string, processor Processor, oo ...async.OptionFunc) (*KafkaComponent, error) {

	kafkaCmp := KafkaComponent{}

	cf, err := kafka.New(name, json.Type, topic, group, []string{broker})
	if err != nil {
		return nil, err
	}

	cmp, err := async.New("driver-kafka-cmp", processor.Process, cf, oo...)
	if err != nil {
		return nil, err
	}

	kafkaCmp.Component = cmp
	return &kafkaCmp, nil
}

package ingestion

import (
	"fmt"
	"time"

	"github.com/beatlabs/patron"

	"github.com/beatlabs/patron/component/async"
	"github.com/beatlabs/patron/component/async/kafka"

	"github.com/beatlabs/patron/component/async/kafka/group"
	"github.com/beatlabs/patron/encoding/json"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
)

// Processor is an interface for structs able to process kafka messages
type Processor interface {
	Process(msg async.Message) error
	Activate(v bool)
}

// NewPublisher returns a pub socket listening to url
func NewPublisher(urls []string) (mangos.Socket, error) {
	var sock mangos.Socket
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		return nil, fmt.Errorf("can't get new pub socket: %s", err)
	}

	for _, url := range urls {
		if err = sock.Listen(url); err != nil {
			return nil, fmt.Errorf("can't listen on pub socket: %s", err.Error())
		}
	}
	return sock, nil
}

// KafkaComponent is a patron component for receiving and processing kafka messages
type KafkaComponent struct {
	patron.Component
}

// NewKafkaComponent instantiates a new component
func NewKafkaComponent(name, cname, cgroup string, topic, broker []string, processor Processor, rt uint, rtw time.Duration) (*KafkaComponent, error) {

	kafkaCmp := KafkaComponent{}

	cf, err := group.New(name, cgroup, topic, broker, kafka.Decoder(json.DecodeRaw))
	if err != nil {
		return nil, err
	}

	bld := async.New(cname, cf, processor.Process)
	cmp, err := bld.WithRetries(rt).WithRetryWait(rtw).Create()
	if err != nil {
		return nil, err
	}
	kafkaCmp.Component = cmp
	return &kafkaCmp, nil
}

package driver

import (
	"fmt"
	"github.com/beatlabs/patron"
	"github.com/beatlabs/patron/async"
	"github.com/beatlabs/patron/async/kafka"
	"github.com/beatlabs/patron/encoding/json"
	"github.com/beatlabs/patron/errors"
	"github.com/beatlabs/patron/log"
	"time"

	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	_ "nanomsg.org/go/mangos/v2/transport/all"
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

	//dt := time.Unix(int64(driver.RegistrationDate), 0)
	//fmt.Println("REGDATE", dt)
	return nil
}

func NewKafkaComponent(name, broker, topic, group string) (*KafkaComponent, error) {
	var sock mangos.Socket
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		log.Fatal("can't get new pub socket: %s", err)
	}
	if err = sock.Listen("inproc://driver-publisher"); err != nil {
		log.Fatal("can't listen on pub socket: %s", err.Error())
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

	return &kafkaCmp, nil
}

package driver

import (
	"bollobas"
	"bollobas/ingestion/injestionfakes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	_ "nanomsg.org/go/mangos/v2/transport/all"
	"sync"
	"testing"
)

func TestProcessing(t *testing.T) {
	cp, err := NewKafkaComponent("component 1", "broker a", "topic a", "group a")
	assert.Nil(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {

		sock, err := sub.NewSocket()

		assert.Nil(t, err)
		err = sock.Dial("inproc://driver-publisher")
		assert.Nil(t, err)
		err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
		assert.Nil(t, err)

		rsp, err := sock.Recv()
		assert.Nil(t, err)

		idt := &bollobas.Identity{}
		err = json.Unmarshal(rsp, idt)
		assert.Nil(t, err)

		assert.EqualValues(t, "abc", idt.ID)

		wg.Done()
	}(wg)

	msg := &injestionfakes.FakeMessage{}

	msg.DecodeStub = func(itf interface{}) error {
		dr := itf.(*Driver)
		dr.ID = "abc"

		return nil
	}

	err = cp.Process(msg)
	assert.Nil(t, err)

	assert.EqualValues(t, 1, msg.DecodeCallCount())

	wg.Wait()

	msg.DecodeReturns(errors.Errorf("error1"))

	err = cp.Process(msg)
	assert.NotNil(t, err)
}

func TestBusyPort(t *testing.T) {

	var sock mangos.Socket
	var err error
	sock, _ = pub.NewSocket()
	err = sock.Listen("inproc://driver-publisher")
assert.Nil(t, err)
	cp, err := NewKafkaComponent("component 1", "broker a", "topic a", "group a")
	assert.NotNil(t, err)
	assert.Nil(t,cp)
}

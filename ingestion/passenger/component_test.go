package passenger

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
	purl := "inproc://passenger-publisher"

	cp, err := NewKafkaComponent("component 1", "broker a", "topic a", "group a", purl)
	assert.Nil(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	wg2 := &sync.WaitGroup{}
	wg2.Add(1)

	go func(wg *sync.WaitGroup) {
		wg2.Done()
		sock, err := sub.NewSocket()

		assert.Nil(t, err)
		err = sock.Dial(purl)
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

	wg2.Wait()


	msg := &injestionfakes.FakeMessage{}

	msg.DecodeStub = func(itf interface{}) error {
		psg := itf.(*Passenger)
		psg.ID = "abc"

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

	purl := "inproc://passenger-publisher"
	var sock mangos.Socket
	var err error
	sock, _ = pub.NewSocket()
	sock.Listen(purl)

	cp, err := NewKafkaComponent("component 1", "broker a", "topic a", "group a", purl)
	assert.NotNil(t, err)
	assert.Nil(t,cp)
}

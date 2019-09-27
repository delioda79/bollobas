package ride

import (
	"bollobas"
	"bollobas/ingestion/injestionfakes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/beatlabs/patron/errors"
	"github.com/stretchr/testify/assert"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/protocol/sub"
)

func TestProcessing(t *testing.T) {
	os.Setenv("BOLLOBAS_LOCATION", "test")
	durl := fmt.Sprintf("inproc://%d", time.Now().UnixNano())
	cp, err := NewRideProcessor(durl)
	assert.Nil(t, err)
	cp.Activate(true)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	wg2 := &sync.WaitGroup{}
	wg2.Add(2)

	go func(wg *sync.WaitGroup) {
		wg2.Done()
		sock, err := sub.NewSocket()

		assert.Nil(t, err)
		err = sock.Dial(durl)
		assert.Nil(t, err)
		err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
		assert.Nil(t, err)
		rsp, err := sock.Recv()
		assert.Nil(t, err)
		idt := &bollobas.RideRequest{}
		err = json.Unmarshal(rsp, idt)
		assert.Nil(t, err)

		assert.EqualValues(t, 1, idt.RquestID)

		wg.Done()
	}(wg)

	cp.SetPipeEventHook(func(event mangos.PipeEvent, pipe mangos.Pipe) {
		if event == mangos.PipeEventAttached {
			wg2.Done()
		}
	})

	wg2.Wait()
	msg := &injestionfakes.FakeMessage{}

	msg.DecodeStub = func(itf interface{}) error {
		dr := itf.(*Ride)
		dr.RequestID = 1
		dr.Events = []RideEvent{}

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
	durl := fmt.Sprintf("inproc://%d", time.Now().UnixNano())
	var sock mangos.Socket
	var err error
	sock, _ = pub.NewSocket()
	err = sock.Listen(durl)
	assert.Nil(t, err)
	cp, err := NewRideProcessor(durl)
	assert.NotNil(t, err)
	assert.Nil(t, cp)
}

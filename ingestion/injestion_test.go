package ingestion

import (
	"bollobas"
	"bollobas/ingestion/injestionfakes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	"sync"
	"testing"
)

// CheckComponent is a test helper for the specific components
func CheckComponent(t *testing.T, cp KafkaComponent, dcf func(itf interface{}) error ) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {

		sock, err := sub.NewSocket()

		assert.Nil(t, err)
		err = sock.Dial("inproc://passenger-publisher")
		assert.Nil(t, err)
		err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
		assert.Nil(t, err)

		fmt.Println("Receiving?")
		rsp, err := sock.Recv()
		fmt.Println("Received")

		assert.Nil(t, err)

		idt := &bollobas.Identity{}
		err = json.Unmarshal(rsp, idt)
		assert.Nil(t, err)

		assert.EqualValues(t, "abc", idt.ID)

		wg.Done()
	}(wg)

	msg := &injestionfakes.FakeMessage{}

	msg.DecodeStub = dcf

	err := cp.Process(msg)
	assert.Nil(t, err)

	assert.EqualValues(t, 1, msg.DecodeCallCount())

	wg.Wait()

	msg.DecodeReturns(errors.Errorf("error1"))

	err = cp.Process(msg)
	assert.NotNil(t, err)
}
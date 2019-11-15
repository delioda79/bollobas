package cancellation

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/beatlabs/patron/errors"
	"github.com/beatlabs/patron/log"
	"github.com/dukex/mixpanel"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas"
	"github.com/taxibeat/bollobas/mixpanel/mixpanelfakes"
	_ "github.com/taxibeat/bollobas/mixpanel/mixpanelfakes"
	"github.com/taxibeat/bollobas/pkg/logging/store"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
)

func TestCancellationRequestProcessWrongFormat(t *testing.T) {

	wg := &sync.WaitGroup{}
	wg.Add(2)
	log.Setup(store.FactoryLogger, nil)

	cl := &mixpanelfakes.FakeMixpanel{}
	cl.UpdateStub = func(s string, update *mixpanel.Update) error {
		wg.Done()
		return nil
	}

	p := Processor{Mixpanel: cl}

	err := p.Process([]byte("dhjzjvkhvcxkjhvckjhvcxkjx"))

	assert.Equal(t, errors.New("error unmarshaling the data").Error(), err.Error()[0:27])
}

func TestProcessingCorrectCancellationRequestMessage(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	log.Setup(store.FactoryLogger, nil)

	cl := &mixpanelfakes.FakeMixpanel{}
	cl.UpdateStub = func(s string, update *mixpanel.Update) error {
		wg.Done()
		return nil
	}

	p := Processor{Mixpanel: cl}

	idt := &bollobas.RideRequestCancellation{}
	msg, err := json.Marshal(idt)
	assert.Nil(t, err)
	err = p.Process(msg)
	assert.Nil(t, err)

	assert.Equal(t, 1, cl.UpdateCallCount())

}

func TestImpossibleCancellationRequestUpdate(t *testing.T) {
	store.NewLogger()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	cl := &mixpanelfakes.FakeMixpanel{}
	cl.UpdateStub = func(s string, update *mixpanel.Update) error {
		wg.Done()
		return errors.New("impossible to update")
	}

	p := Processor{Mixpanel: cl}

	idt := &bollobas.RideRequestCancellation{}
	msg, err := json.Marshal(idt)
	assert.Nil(t, err)
	err = p.Process(msg)
	assert.Equal(t, errors.New("error while updating the CancellationRequest: impossible to update").Error(), err.Error())

	assert.Equal(t, 1, cl.UpdateCallCount())

}

package mixpanel

import (
	"bollobas"
	"bollobas/mixpanel/mixpanelfakes"
	_ "bollobas/mixpanel/mixpanelfakes"
	"bollobas/pkg/logging/store"
	"encoding/json"
	"fmt"
	"github.com/beatlabs/patron/log"
	"github.com/dukex/mixpanel"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
	"sync"
	"testing"
	"time"
)

func TestRecevedFormatError(t *testing.T) {
	store.NewLogger()
	cnt := 1
	wg := &sync.WaitGroup{}
	wg.Add(2)
	log.Setup(store.FactoryLogger, nil)
	sck := &mixpanelfakes.FakeSocket{}
	sck.RecvStub= func () ([]byte, error) {
		if cnt <1 {
			wg.Done()
			select{}
		}
		wg.Done()
		cnt--
		return []byte("dsghdsgjhgdsjhdsa"), nil
	}
	hdl := Handler{Socket: sck}
	hdl.Run()
	wg.Wait()
	errs := store.GetLogger().GetErrors()
	assert.Len(t, errs, 1)

	assert.Equal(t, "error", errs[0][0])
	assert.Equal(t, "error while receiving message", errs[0][1])
}

func TestReceivingError(t *testing.T) {
	store.NewLogger()
	cnt := 1
	wg := &sync.WaitGroup{}
	wg.Add(2)
	log.Setup(store.FactoryLogger, nil)
	sck := &mixpanelfakes.FakeSocket{}
	sck.RecvStub= func () ([]byte, error) {
		if cnt <1 {
			wg.Done()
			select{}
		}
		wg.Done()
		cnt--
		return []byte(""), errors.Errorf("An error")
	}
	hdl := Handler{Socket: sck}
	hdl.Run()
	wg.Wait()

	errs := store.GetLogger().GetErrors()
	assert.Len(t, errs, 1)

	assert.Equal(t, "error", errs[0][0])
	assert.Equal(t, "cannot recv: %s", errs[0][1])
}


func TestReceivingCorrectMessage(t *testing.T) {
	store.NewLogger()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	cl := &mixpanelfakes.FakeMixpanel{}
	cl.UpdateStub = func(s string, update *mixpanel.Update) error {
		wg.Done()
		return nil
	}
	cnt := 1

	log.Setup(store.FactoryLogger, nil)
	sck := &mixpanelfakes.FakeSocket{}
	sck.RecvStub= func () ([]byte, error) {
		if cnt <1 {
			select{}
		}
		wg.Done()
		cnt--

		idt := bollobas.Identity{}
		bts, _ := json.Marshal(idt)
		return bts, nil
	}
	hdl := Handler{Socket: sck, Mixpanel: cl}
	hdl.Run()
	wg.Wait()

	errs := store.GetLogger().GetErrors()
	assert.Len(t, errs, 0)

	assert.Equal(t, 1, cl.UpdateCallCount())

}

func TestImpossibleUpdate(t *testing.T) {
	store.NewLogger()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	cl := &mixpanelfakes.FakeMixpanel{}
	cl.UpdateStub = func(s string, update *mixpanel.Update) error {
		wg.Done()
		return errors.New("impossible to update")
	}
	cnt := 1

	log.Setup(store.FactoryLogger, nil)
	sck := &mixpanelfakes.FakeSocket{}
	sck.RecvStub= func () ([]byte, error) {
		if cnt <1 {
			wg.Done()
			select{}
		}
		wg.Done()
		cnt--

		idt := bollobas.Identity{}
		bts, _ := json.Marshal(idt)
		return bts, nil
	}
	hdl := Handler{Socket: sck, Mixpanel: cl}
	hdl.Run()
	wg.Wait()

	errs := store.GetLogger().GetErrors()
	assert.Len(t, errs, 1)

	assert.Equal(t, 1, cl.UpdateCallCount())

}


func TestGettingNewHandler(t *testing.T) {

	store.NewLogger()
	log.Setup(store.FactoryLogger, nil)


	var s, s2 mangos.Socket
	var err error
	s, err = pub.NewSocket()
	assert.Nil(t, err)
	s2, err = pub.NewSocket()
	assert.Nil(t, err)

	url1 := fmt.Sprintf("inproc://%d", time.Now().UnixNano())
	url2 := fmt.Sprintf("inproc://%d", time.Now().UnixNano())
	err = s.Listen(url1)
	assert.Nil(t, err)
	err = s2.Listen(url2)
	assert.Nil(t, err)

	NewHandler("atoken", []string{url1, url2})

	logs := store.GetLogger().GetErrors()

	assert.Len(t, logs, 2)
}
package mixpanel

import (
	"bollobas"
	"bollobas/mixpanel/mixpanelfakes"
	_ "bollobas/mixpanel/mixpanelfakes"
	"bollobas/pkg/logging/store"
	"encoding/json"
	"fmt"
	"github.com/beatlabs/patron/errors"
	"github.com/beatlabs/patron/log"
	"github.com/dukex/mixpanel"
	"github.com/stretchr/testify/assert"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	_ "nanomsg.org/go/mangos/v2/transport/inproc"
	"sync"
	"testing"
	"time"
)

type FakeProcessor struct {
	mixpanel.Mixpanel

	procreturn error
}

// Run starts the go routine which will receive the messages
func (p *FakeProcessor) Process(msg []byte) error {

	return p.procreturn
}



func TestReceivedFormatError(t *testing.T) {
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
		return []byte("dddfsfddsfds"), nil
	}

	idp := &FakeProcessor{procreturn: errors.Errorf("error while receiving message")}
	hdl := Handler{p: idp, Socket: sck}
	hdl.Run()
	wg.Wait()
	errs := store.GetLogger().GetErrors()
	assert.Len(t, errs, 1)

	assert.Equal(t, "error", errs[0][0])
	act := errs[0][1].(error)
	assert.Equal(t, errors.New("error while receiving message").Error(), act.Error())
}

func TestReceivingCorrectMessage(t *testing.T) {
	store.NewLogger()

	wg := &sync.WaitGroup{}
	wg.Add(1)

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
	idp := &FakeProcessor{procreturn: nil}
	hdl := Handler{p: idp, Socket: sck}
	hdl.Run()
	wg.Wait()

	errs := store.GetLogger().GetErrors()
	assert.Len(t, errs, 0)

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

	NewHandler(&FakeProcessor{}, []string{url1, url2})

	logs := store.GetLogger().GetErrors()

	assert.Len(t, logs, 2)
}
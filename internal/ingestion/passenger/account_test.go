package passenger

import (
	"encoding/json"
	"fmt"
	bollobas "github.com/taxibeat/bollobas/internal/mixpanel"
	"nanomsg.org/go/mangos/v2"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/taxibeat/bollobas/internal/ingestion"
	"github.com/taxibeat/bollobas/internal/ingestion/injestionfakes"
	"github.com/taxibeat/bollobas/internal/mixpanel/pkg/ciphrest"
	"github.com/taxibeat/bollobas/internal/mixpanel/pkg/parseid"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

func TestProcessing(t *testing.T) {
	purl := fmt.Sprintf("inproc://passenger-publisher-%d", time.Now().UnixNano())

	cp, err := NewAccountProcessor(purl, "", "", "")
	assert.Nil(t, err)
	cp.Activate(true)

	msg := &injestionfakes.FakeMessage{}

	msg.DecodeStub = func(itf interface{}) error {
		psg := itf.(*Passenger)
		psg.ID = "abc"

		return nil
	}

	HelpProcessing(t, purl, cp, msg)
}

func TestAccountBusyPorts(t *testing.T) {
	purl := fmt.Sprintf("inproc://passenger-publisher-%d", time.Now().UnixNano())
	HelpBusyPort(t, purl, func(url string) (ingestion.Processor, error) { return NewAccountProcessor(url, "", "", "") })
}

func HelpBusyPort(t *testing.T, url string, factory func(string) (ingestion.Processor, error)) {

	purl := fmt.Sprintf("inproc://passenger-publisher-%d", time.Now().UnixNano())
	var sock mangos.Socket
	var err error
	sock, _ = pub.NewSocket()
	sock.Listen(purl)

	cp, err := factory(purl)
	assert.NotNil(t, err)
	assert.Nil(t, cp)
}

func HelpProcessing(t *testing.T, purl string, cp ingestion.Processor, msg *injestionfakes.FakeMessage) {

	wg := &sync.WaitGroup{}
	wg.Add(1)

	wg2 := &sync.WaitGroup{}
	wg2.Add(2)

	go func(wg *sync.WaitGroup) {
		sock, err := sub.NewSocket()

		assert.Nil(t, err)
		err = sock.Dial(purl)
		assert.Nil(t, err)
		err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
		assert.Nil(t, err)

		wg2.Done()

		rsp, err := sock.Recv()
		assert.Nil(t, err)

		idt := &bollobas.Identity{}
		err = json.Unmarshal(rsp, idt)
		assert.Nil(t, err)

		wg.Done()
	}(wg)

	cp.SetPipeEventHook(func(event mangos.PipeEvent, pipe mangos.Pipe) {
		if event == mangos.PipeEventAttached {
			wg2.Done()
		}
	})
	wg2.Wait()

	err := cp.Process(msg)
	assert.Nil(t, err)

	assert.EqualValues(t, 1, msg.DecodeCallCount())
	wg.Wait()
	msg.DecodeReturns(errors.Errorf("error1"))

	err = cp.Process(msg)
	assert.NotNil(t, err)
}

func TestIDEnc(t *testing.T) {
	os.Setenv("BOLLOBAS_LOCATION", "sandbox")

	err := ciphrest.InitCipher("44441s111111R1222221", "11111111112222222222333333333344")
	assert.NoError(t, err)

	id2 := "VU90NVIwTzhWWHB0b0dSZ3dwcnVMZz09OjoRERERESIiIiIiMzMzMzNE-sandbox-pa"
	decryptedString, err := parseid.DecryptString(id2)
	assert.NoError(t, err)
	idInt, err := strconv.Atoi(decryptedString)
	assert.NoError(t, err)
	encd, err := parseid.EncryptString(idInt, "pa")
	assert.NoError(t, err)
	assert.Equal(t, id2, encd)
}

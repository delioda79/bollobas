package passenger

import (
	"bollobas/ingestion"
	"bollobas/ingestion/injestionfakes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestRequestProcessing(t *testing.T) {
	os.Setenv("BOLLOBAS_LOCATION", "local")
	purl :=  fmt.Sprintf("inproc://passenger-publisher-%d", time.Now().UnixNano())

	cp, err := NewRequestProcessor(purl)
	assert.Nil(t, err)
	cp.Activate(true)

	msg := &injestionfakes.FakeMessage{}

	msg.DecodeStub = func(itf interface{}) error {
		psg := itf.(*RequestRide)
		psg.RequestID = 1

		return nil
	}

	HelpProcessing(t, purl, cp, msg)
}


func TestRequestBusyPorts(t *testing.T) {
	purl := fmt.Sprintf("inproc://passenger-cancellation-%d", time.Now().UnixNano())

	HelpBusyPort(t, purl, func(url string) (ingestion.Processor, error) {return NewRequestProcessor(url)})
}
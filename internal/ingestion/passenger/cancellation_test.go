package passenger

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal/ingestion"
	"github.com/taxibeat/bollobas/internal/ingestion/injestionfakes"
)

func TestCancellationProcessing(t *testing.T) {
	os.Setenv("BOLLOBAS_LOCATION", "local")
	purl := fmt.Sprintf("inproc://passenger-publisher-%d", time.Now().UnixNano())

	cp, err := NewCancellationProcessor(purl, "", "")
	assert.Nil(t, err)
	cp.Activate(true)

	msg := &injestionfakes.FakeMessage{}

	msg.DecodeStub = func(itf interface{}) error {
		psg := itf.(*CancelRideRequest)
		psg.RequestID = 1

		return nil
	}

	HelpProcessing(t, purl, cp, msg)
}

func TestCancellationBusyPorts(t *testing.T) {
	purl := fmt.Sprintf("inproc://passenger-cancellation-%d", time.Now().UnixNano())

	HelpBusyPort(t, purl, func(url string) (ingestion.Processor, error) { return NewCancellationProcessor(url, "", "") })
}

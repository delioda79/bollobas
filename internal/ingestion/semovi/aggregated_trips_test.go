package semovi

import (
	"testing"

	"github.com/taxibeat/bollobas/internal/ingestion/injestionfakes"
	"github.com/taxibeat/bollobas/internal/internalfakes"

	"github.com/stretchr/testify/assert"
)

func TestNewAggregatedTripsProcessor(t *testing.T) {
	repo := &internalfakes.FakeAggregatedTripsRepository{}
	proc := NewAggregatedTripsProcessor(repo)

	assert.NotNil(t, proc)
}

func TestAggregatedTrips(t *testing.T) {
	repo := &internalfakes.FakeAggregatedTripsRepository{}
	repo.AddReturnsOnCall(0, nil)
	proc := NewAggregatedTripsProcessor(repo)
	proc.Activate(true)

	msg := &injestionfakes.FakeMessage{}
	msg.AckReturnsOnCall(0, nil)

	err := proc.Process(msg)
	assert.Nil(t, err)
	assert.Equal(t, 1, msg.AckCallCount())
	assert.Equal(t, 1, repo.AddCallCount())
}

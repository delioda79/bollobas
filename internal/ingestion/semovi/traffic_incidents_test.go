package semovi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal/ingestion/injestionfakes"
	"github.com/taxibeat/bollobas/internal/internalfakes"
)

func TestNewTrafficIncidentsProcessor(t *testing.T) {
	repo := &internalfakes.FakeTrafficIncidentsRepository{}
	proc := NewTrafficIncidentsProcessor(repo)

	assert.NotNil(t, proc)
}

func TestAddTrafficIncidents(t *testing.T) {
	repo := &internalfakes.FakeTrafficIncidentsRepository{}
	repo.AddReturnsOnCall(0, nil)
	proc := NewTrafficIncidentsProcessor(repo)
	proc.Activate(true)

	msg := &injestionfakes.FakeMessage{}
	msg.AckReturnsOnCall(0, nil)

	err := proc.Process(msg)
	assert.Nil(t, err)
	assert.Equal(t, 1, msg.AckCallCount())
	assert.Equal(t, 1, repo.AddCallCount())
}

package semovi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal/ingestion/injestionfakes"
	"github.com/taxibeat/bollobas/internal/internalfakes"
)

func TestNewOperatorStatsProcessor(t *testing.T) {
	repo := &internalfakes.FakeOperatorStatsRepository{}
	proc := NewOperatorStatsProcessor(repo)

	assert.NotNil(t, proc)
}

func TestAdd(t *testing.T) {
	repo := &internalfakes.FakeOperatorStatsRepository{}
	repo.AddReturnsOnCall(0, nil)
	proc := NewOperatorStatsProcessor(repo)
	proc.Activate(true)

	msg := &injestionfakes.FakeMessage{}
	msg.AckReturnsOnCall(0, nil)

	err := proc.Process(msg)
	assert.Nil(t, err)
	assert.Equal(t, 1, msg.AckCallCount())
	assert.Equal(t, 1, repo.AddCallCount())
}

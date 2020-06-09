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

func TestAddOperatorStats(t *testing.T) {
	repo := &internalfakes.FakeOperatorStatsRepository{}
	repo.AddReturnsOnCall(0, nil)
	proc := NewOperatorStatsProcessor(repo)
	proc.Activate(true)

	msg := &injestionfakes.FakeMessage{}
	msg.PayloadReturns([]byte(`{"fecha_produccion":1589279257759,"id_operador":"2d2ec778-b89e-4db5-9628-123fd99f0b91","genero":1,"cant_viajes":29,"tiempo_registro":44,"edad":"28-32","horas_conectado":"9-17","horas_viaje":"0-24","ingreso_totales":"$0-999"}`))
	msg.AckReturnsOnCall(0, nil)

	err := proc.Process(msg)
	assert.Nil(t, err)
	assert.Equal(t, 1, msg.AckCallCount())
	assert.Equal(t, 1, repo.AddCallCount())
}

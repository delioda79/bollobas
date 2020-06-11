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
	msg.PayloadReturns([]byte(`{"fecha": 1589279257759,"id_proveedor": "test","tot_viajes": 12,"tot_veh_viaje": 11,"tot_veh_disp": 1,"dist_pasajero": 1,"tiempo_pasajero": 1,"tiempo_solicitud": 1,"tiempo_vacio": 1,"multiplicador_eod": 1,"accesibilidad": 1,"operador_mujer": 1,"inicio_eod": 1,"fin_eod": 1,"dist_pasajero_eod": 1,"tiempo_pasajero_eod": 1,"dist_solicitud": 1,"dist_vacio": 1,"dist_solicitud_eod": 1,"tiempo_solicitud_eod": 1,"dist_vacio_eod": 1,"tiempo_vacio_eod": 1}`))
	msg.AckReturnsOnCall(0, nil)

	err := proc.Process(msg)
	assert.Nil(t, err)
	assert.Equal(t, 1, msg.AckCallCount())
	assert.Equal(t, 1, repo.AddCallCount())
}

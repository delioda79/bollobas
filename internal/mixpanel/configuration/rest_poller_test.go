package configuration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taxibeat/bollobas/internal/mixpanel"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPoll(t *testing.T) {

	tcs := []struct {
		status int
		cfg    []byte
		xp     bool
	}{
		{status: http.StatusOK, cfg: []byte(`{"mixpanel_passenger_enabled": true}`), xp: true},
		{status: http.StatusOK, cfg: []byte(`{"mixpanel_passenger_enabled": false}`), xp: false},
		{status: http.StatusNotFound, cfg: []byte(`{"mixpanel_passenger_enabled": true}`), xp: false},
	}

	for _, v := range tcs {
		t.Run(fmt.Sprintf("Status: %v", v.status), func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

				rw.WriteHeader(v.status)
				rw.Write(v.cfg)
			}))

			defer ts.Close()

			cnf := mixpanel.NewConfigurator()
			pll := RestPoller{
				Manager:       cnf,
				RestURL:       ts.URL,
				PollingPeriod: time.Millisecond * 100,
			}

			pll.UpdateSettings()

			time.Sleep(time.Millisecond * 150)
			assert.Equal(t, cnf.Check([]byte("")), v.xp)
		})
	}

}

package configuration

import (
	"bollobas"
	"bollobas/pkg/configclient"
	"context"
	"github.com/beatlabs/patron/log"
	"time"
)

// RestPoller is use dto poll rest for configuration updates
type RestPoller struct {
	Manager       bollobas.ConfigurationManager
	RestURL       string
	PollingPeriod time.Duration
	DefaultConf   map[string] interface{}
	RestKey       string
	Path          string
}

// UpdateSettings updates the settings
func (rp RestPoller) UpdateSettings() {
	if rp.RestURL == "" {
		rp.Manager.Configure(rp.DefaultConf)
		return
	}

	ticker := time.NewTicker(rp.PollingPeriod)
	cClient, err := configclient.New(rp.RestURL, rp.RestKey, rp.Path)
	if err != nil {
		log.Debugf("Couldn't create Configuration Client. Resolving to defaults: %v", rp.DefaultConf)
		return
	}

	start := time.Now()
	st, err := cClient.GetSettings(context.TODO())
	if err == nil {
		//Configure
		rp.Manager.Configure(st)
		ObserveCount("mixpanel", rp.RestURL, true)
		log.Debugf("Settings updated with: %v", st)
	} else {
		log.Infof("Failed to update settings: %v", err)
		ObserveCount("mixpanel", rp.RestURL, false)
	}

	ObserveLatency("mixpanel", rp.RestURL, time.Since(start))

	go func() {
		for {
			<-ticker.C
			//Logic to get configs here....
			start := time.Now()
			st, err := cClient.GetSettings(context.TODO())
			if err == nil {
				//Configure
				rp.Manager.Configure(st)
				log.Debugf("Settings updated with: %v", st)
				ObserveCount("mixpanel", rp.RestURL, true)
			} else {
				log.Infof("Failed to update settings: %v", err)
				ObserveCount("mixpanel", rp.RestURL, false)
			}
			ObserveLatency("mixpanel", rp.RestURL, time.Since(start))
		}
	}()
}
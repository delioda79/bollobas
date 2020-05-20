package config

import (
	"fmt"
	bsync "sync"
	"time"

	"github.com/beatlabs/harvester"
	"github.com/beatlabs/harvester/sync"
)

// Duration type with concurrent access support.
type Duration struct {
	rw    bsync.RWMutex
	value time.Duration
}

// Get returns the internal value.
func (d *Duration) Get() time.Duration {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.value
}

// Set a value.
func (d *Duration) Set(value time.Duration) {
	d.rw.Lock()
	defer d.rw.Unlock()
	d.value = value
}

// String returns string representation of value.
func (d *Duration) String() string {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.value.String()
}

// SetString parses and sets a value from string type.
func (d *Duration) SetString(val string) error {
	v, err := time.ParseDuration(val)
	if err != nil {
		return fmt.Errorf("env %s is not a duration: %v", val, err)
	}
	d.Set(v)
	return nil
}

// Configuration holds all the configuration for harvester
type Configuration struct {
	KafkaBroker                 sync.String `env:"BOLLOBAS_KAFKA_CONNECTION_STRING"`
	KafkaGroup                  sync.String `env:"BOLLOBAS_KAFKA_GROUP"`
	KafkaDriverIdentityTopic    sync.String `seed:"driver_analytics" env:"BOLLOBAS_KAFKA_DRIVER_TOPIC"`
	KafkaPassengerIdentityTopic sync.String `seed:"passenger_analytics" env:"BOLLOBAS_KAFKA_PASSENGER_TOPIC"`
	KafkaTimeout                Duration    `seed:"2s" env:"BOLLOBAS_KAFKA_TIMEOUT"`
	RestURL                     sync.String `seed:"" env:"REST_CONNECTION_STRING"`
	RestKey                     sync.String `seed:"" env:"REST_KEY"`
	MpToken                     sync.String `env:"BOLLOBAS_MIXPANEL_TOKEN"`
	KkPRRTopic                  sync.String `seed:"request" env:"BOLLOBAS_KAFKA_REQUEST_TOPIC"`
	KkPRCTopic                  sync.String `seed:"request_cancel" env:"BOLLOBAS_KAFKA_REQUEST_CANCEL_TOPIC"`
	KkRTopic                    sync.String `seed:"ride" env:"BOLLOBAS_KAFKA_RIDE_TOPIC"`
	KkSOTopic                   sync.String `seed:"stats_operador" env:"BOLLOBAS_KAFKA_STATS_OPERADOR_TOPIC"`
	BConf                       sync.String `seed:"{}" env:"BOLLOBAS_BASE_CONF"`
	RestMixpanelPath            sync.String `seed:"/taxidmin/bollobas/mixpanel-passenger-settings" env:"REST_MIXPANEL_PATH"`
	CipherKey                   sync.String `seed:"" env:"BOLLOBAS_CIPHER_KEY"`
	CipherInitVec               sync.String `seed:"" env:"BOLLOBAS_INIT_VECTOR"`
	Location                    sync.String `seed:"" env:"BOLLOBAS_LOCATION"`
	DBUsername                  sync.String `env:"MYSQL_USERNAME"`
	DBPassword                  sync.String `env:"MYSQL_PASS"`
	DBWriteHost                 sync.String `env:"MYSQL_WRITE"`
	DBReadHost                  sync.String `env:"MYSQL_READ"`
	DBPort                      sync.String `env:"MYSQL_PORT"`
	DBName                      sync.String `env:"MYSQL_DB"`
	SettingsPeriod              Duration    `seed:"60s" env:"BOLLOBAS_SETTINGS_DURATION"`
}

// NewConfig instantiates a new configuration object
func NewConfig(cfg *Configuration) (harvester.Harvester, error) {
	h, err := harvester.New(cfg).Create()

	return h, err
}

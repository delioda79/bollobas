package configuration

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	topicCounter  = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bollobas",
			Subsystem: "configuration",
			Name:      "messages_total",
			Help:      "Counts every time a configuration is attempted",
		},
		[]string{ "entity","url", "success"},
	)

	messageProcessingLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bollobas",
			Subsystem: "configuration",
			Name:      "latency_seconds",
			Help:      "Latency of configuration update",
		},
		[]string{"entity","url"},
	)
)

func init() {
	prometheus.MustRegister(topicCounter)
	prometheus.MustRegister(messageProcessingLatency)
}

// ObserveCount is responsible to count distance matrix calls
func ObserveCount(entity, url string, success bool) {
	topicCounter.WithLabelValues( entity, url, strconv.FormatBool(success)).Inc()
}

// ObserveLatency is responsible to observe latency.
func ObserveLatency(entity, url string, latency time.Duration) {
	messageProcessingLatency.WithLabelValues(entity, url).Observe(latency.Seconds())
}
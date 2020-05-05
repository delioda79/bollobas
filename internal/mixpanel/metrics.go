package mixpanel

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	topicCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bollobas",
			Subsystem: "mixpanel",
			Name:      "messages_total",
			Help:      "Counts every time a message gets receiver by the mixpanel services",
		},
		[]string{"provider", "topic", "processed", "success"},
	)

	messageProcessingLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bollobas",
			Subsystem: "mixpanel",
			Name:      "latency_seconds",
			Help:      "Latency of received messages processing",
		},
		[]string{"provider", "topic"},
	)
)

func init() {
	prometheus.MustRegister(topicCounter)
	prometheus.MustRegister(messageProcessingLatency)
}

// ObserveCount is responsible to count distance matrix calls
func ObserveCount(provider, topic string, processed, success bool) {
	topicCounter.WithLabelValues(provider, topic, strconv.FormatBool(processed), strconv.FormatBool(success)).Inc()
}

// ObserveLatency is responsible to observe latency.
func ObserveLatency(provider, topic string, latency time.Duration) {
	messageProcessingLatency.WithLabelValues(provider, topic).Observe(latency.Seconds())
}

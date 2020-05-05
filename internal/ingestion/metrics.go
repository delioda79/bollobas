package ingestion

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	topicCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bollobas",
			Subsystem: "injestion",
			Name:      "ingestion_messages_total",
			Help:      "Counts every time a message gets receiver by the ingestion services",
		},
		[]string{"provider", "topic", "processed", "success"},
	)

	repblishedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bollobas",
			Subsystem: "injestion",
			Name:      "republished_messages_total",
			Help:      "Counts every time a message is republished via mangos",
		},
		[]string{"message_type", "entity_type"},
	)

	messageProcessingLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bollobas",
			Subsystem: "injestion",
			Name:      "latency_seconds",
			Help:      "Latency of ingested messages processing",
		},
		[]string{"provider", "topic"},
	)
)

func init() {
	prometheus.MustRegister(topicCounter)
	prometheus.MustRegister(repblishedCounter)
	prometheus.MustRegister(messageProcessingLatency)
}

// ObserveCount is responsible to count distance matrix calls
func ObserveCount(provider, topic string, processed, success bool) {
	topicCounter.WithLabelValues(provider, topic, strconv.FormatBool(processed), strconv.FormatBool(success)).Inc()
}

// ObserveRepublishedCount is responsible to count the active distance matrix providers used
func ObserveRepublishedCount(messageType, entityType string) {
	repblishedCounter.WithLabelValues(messageType, entityType).Inc()
}

// ObserveLatency is responsible to observe latency.
func ObserveLatency(provider, topic string, latency time.Duration) {
	messageProcessingLatency.WithLabelValues(provider, topic).Observe(latency.Seconds())
}

package ingestion

import "github.com/beatlabs/patron/async"

type KafkaComponent interface {
	Process(msg async.Message) error
}
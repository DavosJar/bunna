package buffer

import "time"

// Config holds configuration options for the telemetry buffer system.
type Config struct {
	Capacity             int           `json:"capacity"`
	BatchSize            int           `json:"batchSize"`
	FlushIntervalSeconds int           `json:"flushIntervalSeconds"`
	MaxRetries           int           `json:"maxRetries"`
	BackoffBase          time.Duration `json:"backoffBase"`
	BackoffMax           time.Duration `json:"backoffMax"`
	KafkaBrokers         []string      `json:"kafkaBrokers"`
	KafkaTopic           string        `json:"kafkaTopic"`
}

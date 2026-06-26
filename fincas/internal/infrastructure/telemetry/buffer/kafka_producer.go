package buffer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer is a real producer for telemetry events using kafka-go.
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new real producer using configuration.
func NewKafkaProducer(cfg Config) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaBrokers...),
			Topic:    cfg.KafkaTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// Publish sends a batch of events to Kafka.
func (p *KafkaProducer) Publish(batch [][]byte) error {
	messages := make([]kafka.Message, len(batch))
	for i, b := range batch {
		messages[i] = kafka.Message{
			Value: b,
		}
	}
	return p.writer.WriteMessages(context.Background(), messages...)
}

// Close performs any cleanup.
func (p *KafkaProducer) Close() {
	p.writer.Close()
}

package eventpublisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaPublisher implementa application.EventPublisher usando Kafka.
type KafkaPublisher struct {
	writer *kafka.Writer
}

// NewKafkaPublisher crea un publisher que escribe a un topic Kafka.
func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}
	return &KafkaPublisher{writer: w}
}

// Publish serializa el evento como JSON y lo envía a Kafka.
// Usa routingKey como key del mensaje (para que consumidores filtren por key).
func (p *KafkaPublisher) Publish(ctx context.Context, routingKey string, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializando evento: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(routingKey),
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("error publicando en Kafka: %w", err)
	}

	log.Printf("[KAFKA] %s → %s", routingKey, string(data))
	return nil
}

// Close cierra el writer de Kafka.
func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

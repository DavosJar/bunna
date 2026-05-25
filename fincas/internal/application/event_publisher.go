package application

import "context"

// EventPublisher es la interfaz para publicar eventos en RabbitMQ.
// La implementación concreta vive en la capa de infraestructura.
type EventPublisher interface {
	Publish(ctx context.Context, routingKey string, event any) error
}

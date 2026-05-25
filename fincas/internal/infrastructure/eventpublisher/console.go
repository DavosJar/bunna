package eventpublisher

import (
	"context"
	"encoding/json"
	"log"
)

// ConsolePublisher implementa application.EventPublisher escribiendo
// los eventos a la consola (stdout) en formato JSON.
type ConsolePublisher struct{}

// NewConsolePublisher crea una nueva instancia de ConsolePublisher.
func NewConsolePublisher() *ConsolePublisher {
	return &ConsolePublisher{}
}

// Publish serializa el evento como JSON y lo escribe en el log.
func (p *ConsolePublisher) Publish(_ context.Context, routingKey string, event any) error {
	data, _ := json.Marshal(event)
	log.Printf("[EVENTO] %s → %s", routingKey, string(data))
	return nil
}

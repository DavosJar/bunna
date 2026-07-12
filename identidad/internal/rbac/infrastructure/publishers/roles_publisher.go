package publishers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventoRolActualizado struct {
	EventID   string    `json:"event_id"`
	Tipo      string    `json:"tipo"`
	RolID     string    `json:"rol_id"`
	TenantID  string    `json:"tenant_id"`
	Permisos  []string  `json:"permisos"`
	OcurredAt time.Time `json:"ocurred_at"`
}

type RolesPublisher struct {
	writer *kafka.Writer
}

func NewRolesPublisher(brokers []string, topic string) *RolesPublisher {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &RolesPublisher{writer: w}
}

func (p *RolesPublisher) PublicarRolActualizado(ctx context.Context, rolID, tenantID string, permisos []string) error {
	evento := EventoRolActualizado{
		EventID:   "", // Omitido o generar uno si es necesario, pero kafka-go no requiere ID
		Tipo:      "permisos.rol.actualizado",
		RolID:     rolID,
		TenantID:  tenantID,
		Permisos:  permisos,
		OcurredAt: time.Now(),
	}

	b, err := json.Marshal(evento)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(evento.RolID + "-" + evento.TenantID),
		Value: b,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("[RolesPublisher] Error publicando evento rol_actualizado para rol %s: %v", evento.RolID, err)
		return err
	}
	log.Printf("[RolesPublisher] Evento publicado exitosamente para rol: %s (Tenant: %s) con %d permisos", evento.RolID, evento.TenantID, len(evento.Permisos))
	return nil
}

// Writer expone el writer de Kafka para que el outbox worker pueda publicar eventos.
func (p *RolesPublisher) Writer() *kafka.Writer {
	return p.writer
}

func (p *RolesPublisher) Close() error {
	return p.writer.Close()
}

package consumers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	iampostgres "github.com/davosjar/bunna/services/fincas/internal/infrastructure/security/iam/postgres"
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

type RolesConsumer struct {
	reader  *kafka.Reader
	iamRepo *iampostgres.IAMRepositorio
}

func NewRolesConsumer(brokers []string, topic string, iamRepo *iampostgres.IAMRepositorio) *RolesConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "fincas-roles-sync",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &RolesConsumer{
		reader:  r,
		iamRepo: iamRepo,
	}
}

func (c *RolesConsumer) Start(ctx context.Context) {
	log.Printf("[RolesConsumer] Escuchando en tópico: %s", c.reader.Config().Topic)
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("[RolesConsumer] Error leyendo mensaje: %v", err)
			continue
		}

		var event EventoRolActualizado
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("[RolesConsumer] Error unmarshaling: %v", err)
			continue
		}

		if event.Tipo == "permisos.rol.actualizado" {
			log.Printf("[RolesConsumer] Actualizando permisos locales para Rol: %s, Tenant: %s (%d permisos)", event.RolID, event.TenantID, len(event.Permisos))
			if err := c.iamRepo.UpsertPermisos(ctx, event.RolID, event.TenantID, event.Permisos); err != nil {
				log.Printf("[RolesConsumer] Error haciendo upsert en DB local: %v", err)
			}
		}
	}
}

func (c *RolesConsumer) Close() error {
	return c.reader.Close()
}

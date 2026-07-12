package publishers

import (
	"context"
	"encoding/json"
	"time"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/outbox"
	"gorm.io/gorm"
)

// OutboxRolPublisher implementa rbac.RolPublisher escribiendo eventos en la
// tabla event_outbox dentro de la misma transacción de BD. El worker de outbox
// se encarga de publicar estos eventos a Kafka de forma asíncrona.
type OutboxRolPublisher struct {
	db          *gorm.DB
	generadorID shareddomain.GeneradorID
}

func NewOutboxRolPublisher(db *gorm.DB, generadorID shareddomain.GeneradorID) *OutboxRolPublisher {
	return &OutboxRolPublisher{db: db, generadorID: generadorID}
}

func (p *OutboxRolPublisher) PublicarRolActualizado(ctx context.Context, rolID, tenantID string, permisos []string) error {
	eventoID, err := p.generadorID.NextID(ctx)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"event_id":   eventoID,
		"tipo":       "permisos.rol.actualizado",
		"rol_id":     rolID,
		"tenant_id":  tenantID,
		"permisos":   permisos,
		"ocurred_at": time.Now().Format(time.RFC3339),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	evento := &outbox.EventoOutbox{
		ID:          eventoID,
		EventType:   "permisos.rol.actualizado",
		AggregateID: tenantID,
		Payload:     string(payloadJSON),
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
	return p.db.WithContext(ctx).Create(evento).Error
}

package outbox

import "time"

// EventoOutbox representa un evento pendiente de publicar en el topic de Kafka.
// Se persiste en la misma transacción de BD que la operación de negocio,
// garantizando consistencia sin depender de la disponibilidad de Kafka.
type EventoOutbox struct {
	ID          string     `gorm:"column:id;primaryKey;size:255"`
	EventType   string     `gorm:"column:event_type;not null"`
	AggregateID string     `gorm:"column:aggregate_id;not null"`
	Payload     string     `gorm:"column:payload;type:jsonb;not null"`
	Status      string     `gorm:"column:status;not null;default:pending"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	PublishedAt *time.Time `gorm:"column:published_at"`
	RetryCount  int        `gorm:"column:retry_count;not null;default:0"`
	LastError   string     `gorm:"column:last_error"`
}

func (EventoOutbox) TableName() string { return "event_outbox" }

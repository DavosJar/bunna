package outbox

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// OutboxRepositorio define las operaciones de persistencia para eventos del outbox.
type OutboxRepositorio interface {
	Insert(ctx context.Context, evento *EventoOutbox) error
	GetPending(ctx context.Context, limit int) ([]EventoOutbox, error)
	MarkPublished(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, lastError string) error
}

// OutboxRepositorioPostgres es la implementación con GORM para PostgreSQL.
type OutboxRepositorioPostgres struct {
	db *gorm.DB
}

func NewOutboxRepositorioPostgres(db *gorm.DB) *OutboxRepositorioPostgres {
	return &OutboxRepositorioPostgres{db: db}
}

func (r *OutboxRepositorioPostgres) Insert(ctx context.Context, evento *EventoOutbox) error {
	return r.db.WithContext(ctx).Create(evento).Error
}

func (r *OutboxRepositorioPostgres) GetPending(ctx context.Context, limit int) ([]EventoOutbox, error) {
	var eventos []EventoOutbox
	err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&eventos).Error
	return eventos, err
}

func (r *OutboxRepositorioPostgres) MarkPublished(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&EventoOutbox{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "published",
			"published_at": &now,
		}).Error
}

func (r *OutboxRepositorioPostgres) MarkFailed(ctx context.Context, id string, lastError string) error {
	return r.db.WithContext(ctx).
		Model(&EventoOutbox{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     "failed",
			"last_error": lastError,
		}).Error
}

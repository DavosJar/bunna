package idgenerator

import (
	"context"
	"fmt"
	"github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/google/uuid"
)

// Se usa v7 (time-ordered) como principal
type UUIDv7Generator struct{}

// NewUUIDv7Generator crea un nuevo generador de UUIDs
func NewUUIDv7Generator() domain.GeneradorID {
	return &UUIDv7Generator{}
}

// NextID genera un nuevo UUID v7 (Time-ordered)
func (g *UUIDv7Generator) NextID(ctx context.Context) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("error generando uuid v7: %w", err)
	}
	return id.String(), nil
}

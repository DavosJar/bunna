package idgenerator

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/google/uuid"
)

// UUIDv7Generator implementa domain.GeneradorID usando UUIDs v4
// Nota: google/uuid actualmente soporta v4, v5 y v6
// Se usa v4 (random) como fallback hasta que v7 esté disponible
type UUIDv7Generator struct{}

// NewUUIDv7Generator crea un nuevo generador de UUIDs
func NewUUIDv7Generator() domain.GeneradorID {
	return &UUIDv7Generator{}
}

// NextID genera un nuevo UUID
func (g *UUIDv7Generator) NextID(ctx context.Context) (string, error) {
	// TODO: Cambiar a uuid.NewV7() cuando google/uuid lo soporte
	// Por ahora usamos v4 (random) que es seguro
	return uuid.New().String(), nil
}

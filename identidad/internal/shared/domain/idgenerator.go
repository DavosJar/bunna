package domain

import "context"

// GeneradorID define el contrato para generar IDs únicos
// Puede ser implementado con UUID, Snowflake, Nanoid, etc.
// Esta interfaz es reutilizable en cualquier módulo (usuarios, sesiones, eventos, etc.)
type GeneradorID interface {
	// NextID genera el próximo ID único
	NextID(ctx context.Context) (string, error)
}

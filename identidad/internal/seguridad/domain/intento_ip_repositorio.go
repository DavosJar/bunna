package domain

import (
	"context"
	"time"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

// IntentoIPRepositorio define el contrato para persistir intentos de login por IP.
type IntentoIPRepositorio interface {
	// ObtenerPorIP retorna el registro de intentos para una IP en la ventana activa.
	// Retorna error si no existe registro para esa IP.
	ObtenerPorIP(ctx context.Context, ip string) (*IntentoPorIP, error)

	// Crear persiste un nuevo registro de intento para una IP.
	Crear(ctx context.Context, intento *IntentoPorIP) (*IntentoPorIP, error)

	// Actualizar actualiza un registro existente de intentos por IP.
	Actualizar(ctx context.Context, intento *IntentoPorIP) (*IntentoPorIP, error)

	// Listar retorna registros de intentos por IP filtrados y paginados.
	Listar(ctx context.Context, especificacion EspecificacionIntentoIP, paginacion shareddomain.Paginacion) ([]*IntentoPorIP, error)

	// EliminarExpirados elimina registros cuya ventana de tiempo ya expiró.
	// Se invoca periódicamente para mantener la tabla limpia.
	EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error
}
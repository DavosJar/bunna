package decorator

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
)

// AuthUseCase es la interfaz genérica para casos de uso autenticados.
// Firma estándar de la capa de aplicación de fincas.
type AuthUseCase[Cmd, Resp any] interface {
	Ejecutar(context.Context, *application.AuthContext, Cmd) (Resp, error)
}

// UseCase es la interfaz genérica para casos de uso internos sin auth
// (p. ej. consumer RabbitMQ).
type UseCase[Cmd, Resp any] interface {
	Ejecutar(context.Context, Cmd) (Resp, error)
}

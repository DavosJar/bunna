package decorator

import (
	"context"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/buffer"
)

type authUseCaseConstraint[Cmd, Resp any] interface {
	Ejecutar(context.Context, *application.AuthContext, Cmd) (Resp, error)
}

// AuthWrapper decora casos de uso autenticados (capa APO/negocio).
type AuthWrapper[Cmd, Resp any] struct {
	nombre  string
	writer  buffer.BufferWriter
	service telemetry.ServiceInfo
	inner   authUseCaseConstraint[Cmd, Resp]
}

// WrapAuth crea un decorador de telemetría para Ejecutar(ctx, auth, cmd).
func WrapAuth[Cmd, Resp any](
	nombre string,
	writer buffer.BufferWriter,
	service telemetry.ServiceInfo,
	inner authUseCaseConstraint[Cmd, Resp],
) *AuthWrapper[Cmd, Resp] {
	return &AuthWrapper[Cmd, Resp]{nombre: nombre, writer: writer, service: service, inner: inner}
}

func (w *AuthWrapper[Cmd, Resp]) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Cmd) (Resp, error) {
	start := time.Now()
	safeCmd := SafeCommand(cmd)
	resp, err := w.inner.Ejecutar(ctx, auth, cmd)
	duration := time.Since(start).Milliseconds()
	ReportarNegocio(ctx, w.writer, w.service, w.nombre, safeCmd, err, float64(duration))
	return resp, err
}

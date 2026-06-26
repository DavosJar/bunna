package decorator

import (
	"context"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/buffer"
)

type useCaseConstraint[Cmd, Resp any] interface {
	Ejecutar(context.Context, Cmd) (Resp, error)
}

// Wrapper decora casos de uso internos (sin AuthContext).
type Wrapper[Cmd, Resp any] struct {
	nombre  string
	writer  buffer.BufferWriter
	service telemetry.ServiceInfo
	inner   useCaseConstraint[Cmd, Resp]
}

// Wrap crea un decorador de telemetría para casos de uso con Ejecutar(ctx, cmd).
func Wrap[Cmd, Resp any](
	nombre string,
	writer buffer.BufferWriter,
	service telemetry.ServiceInfo,
	inner useCaseConstraint[Cmd, Resp],
) *Wrapper[Cmd, Resp] {
	return &Wrapper[Cmd, Resp]{nombre: nombre, writer: writer, service: service, inner: inner}
}

func (w *Wrapper[Cmd, Resp]) Ejecutar(ctx context.Context, cmd Cmd) (Resp, error) {
	start := time.Now()
	safeCmd := SafeCommand(cmd)
	resp, err := w.inner.Ejecutar(ctx, cmd)
	duration := time.Since(start).Milliseconds()
	ReportarNegocio(ctx, w.writer, w.service, w.nombre, safeCmd, err, float64(duration))
	return resp, err
}

package decorator

import (
	"context"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
)

// useCaseConstraint matches any use case with a single Ejecutar method.
type useCaseConstraint[Cmd, Resp any] interface {
	Ejecutar(context.Context, Cmd) (Resp, error)
}

// Wrapper is a generic telemetry decorator for use cases with exactly one
// Ejecutar method. It produces the same LogPayload as the per-use-case
// decorators it replaces.
type Wrapper[Cmd, Resp any] struct {
	nombre string
	writer buffer.BufferWriter
	inner  useCaseConstraint[Cmd, Resp]
}

// Wrap creates a telemetry-wrapped use case. The returned *Wrapper satisfies
// the same interface as the inner use case by structural typing, so it can
// be assigned to any existing LoginUseCase, RefreshUseCase, RegistroUseCase, etc.
func Wrap[Cmd, Resp any](nombre string, writer buffer.BufferWriter, inner useCaseConstraint[Cmd, Resp]) *Wrapper[Cmd, Resp] {
	return &Wrapper[Cmd, Resp]{nombre: nombre, writer: writer, inner: inner}
}

// Ejecutar delegates to the inner use case while emitting a NEGOCIO telemetry
// event with the same structure as the original per-use-case decorators.
func (w *Wrapper[Cmd, Resp]) Ejecutar(ctx context.Context, cmd Cmd) (Resp, error) {
	start := time.Now()
	safeCmd := SafeCommand(cmd)
	resp, err := w.inner.Ejecutar(ctx, cmd)
	duration := time.Since(start).Milliseconds()
	ReportarNegocio(ctx, w.writer, w.nombre, safeCmd, err, float64(duration))
	return resp, err
}

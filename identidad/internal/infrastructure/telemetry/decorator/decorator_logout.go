package decorator

import (
	"context"
	"encoding/json"
	"time"

	uc_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/middleware"
	presentation_middleware "github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
)

type decoratorLogout struct {
	uc     LogoutUseCase
	writer buffer.BufferWriter
}

func NewDecoratorLogout(uc LogoutUseCase, writer buffer.BufferWriter) *decoratorLogout {
	return &decoratorLogout{uc: uc, writer: writer}
}

func (d *decoratorLogout) Ejecutar(ctx context.Context, cmd uc_logout.ComandoCerrarSesion) (*uc_logout.RespuestaCerrarSesion, error) {
	return d.logWithTelemetry(ctx, "Logout", func() (interface{}, error) {
		return d.uc.Ejecutar(ctx, cmd)
	})
}

func (d *decoratorLogout) CerrarTodas(ctx context.Context, cmd uc_logout.ComandoCerrarTodasLasSesiones) (*uc_logout.RespuestaCerrarSesion, error) {
	return d.logWithTelemetry(ctx, "CerrarTodas", func() (interface{}, error) {
		return d.uc.CerrarTodas(ctx, cmd)
	})
}

func (d *decoratorLogout) logWithTelemetry(ctx context.Context, useCase string, fn func() (interface{}, error)) (*uc_logout.RespuestaCerrarSesion, error) {
	start := time.Now()
	traceID := middleware.GetTraceIDFromCtx(ctx)
	userID := presentation_middleware.GetUsuarioIDFromCtx(ctx)
	resp, err := fn()
	duration := time.Since(start).Milliseconds()
	result := classifyResult(err)
	level := determineLevel(result)
	spanID := middleware.GetSpanIDFromCtx(ctx)
	payload := telemetry.LogPayload{
		LogType:     "NEGOCIO",
		Level:       level,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		TraceID:     traceID,
		SpanID:      spanID,
		ServiceName: "identidad",
		Environment: "dev",
		Negocio: &telemetry.NegocioFields{
			UseCase:           useCase,
			Command:           map[string]any{},
			Result:            result,
			UserID:            userID,
			Details:           map[string]any{},
			DurationUsecaseMs: float64(duration),
		},
	}
	data, _ := json.Marshal(payload)
	prio := buffer.Media
	if level == "ERROR" {
		prio = buffer.Alta
	}
	_ = d.writer.Write(data, prio)
	if err != nil {
		return nil, err
	}
	return resp.(*uc_logout.RespuestaCerrarSesion), nil
}

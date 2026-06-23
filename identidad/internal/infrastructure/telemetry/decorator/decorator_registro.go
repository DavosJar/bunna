package decorator

import (
	"context"
	"encoding/json"
	"time"

	uc_register "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/register"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/middleware"
	presentation_middleware "github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
)

type decoratorRegistro struct {
	uc     RegistroUseCase
	writer buffer.BufferWriter
}

func NewDecoratorRegistro(uc RegistroUseCase, writer buffer.BufferWriter) *decoratorRegistro {
	return &decoratorRegistro{uc: uc, writer: writer}
}

func (d *decoratorRegistro) Ejecutar(ctx context.Context, cmd *uc_register.ComandoRegistrarUsuario) (*uc_register.RespuestaRegistrarUsuario, error) {
	start := time.Now()
	safeCmd := safeCommand(cmd)
	traceID := middleware.GetTraceIDFromCtx(ctx)
	userID := presentation_middleware.GetUsuarioIDFromCtx(ctx)
	resp, err := d.uc.Ejecutar(ctx, cmd)
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
			UseCase:           "Registro",
			Command:           safeCmd,
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
	return resp, err
}

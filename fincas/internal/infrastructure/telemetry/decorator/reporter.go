package decorator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/middleware"
)

// ReportarNegocio construye y encola un evento NEGOCIO correlacionado con la request HTTP.
func ReportarNegocio(
	ctx context.Context,
	writer buffer.BufferWriter,
	service telemetry.ServiceInfo,
	nombreUseCase string,
	comando map[string]any,
	err error,
	duracionMs float64,
) {
	traceID := middleware.GetTraceIDFromCtx(ctx)
	spanID := middleware.GetSpanIDFromCtx(ctx)
	userID := telemetry.GetUsuarioIDFromCtx(ctx)
	result := ClassifyResult(err)
	level := DetermineLevel(result)

	payload := telemetry.LogPayload{
		LogType:     "NEGOCIO",
		Level:       level,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		TraceID:     traceID,
		SpanID:      spanID,
		ServiceName: service.Name,
		Environment: service.Environment,
		Negocio: &telemetry.NegocioFields{
			UseCase:           nombreUseCase,
			Command:           comando,
			Result:            result,
			UserID:            userID,
			Details:           map[string]any{},
			DurationUsecaseMs: duracionMs,
		},
	}
	data, _ := json.Marshal(payload)
	prio := buffer.Media
	if level == "ERROR" {
		prio = buffer.Alta
	}
	_ = writer.Write(data, prio)
}

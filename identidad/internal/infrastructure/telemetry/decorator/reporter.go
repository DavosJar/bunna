package decorator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/middleware"
	presentation_middleware "github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
)

// ReportarNegocio builds and enqueues a NEGOCIO telemetry event.
// It extracts trace_id/span_id from context, classifies the result, and
// writes the exact same LogPayload that the per-use-case decorators produce.
// This is a reusable helper that any wrapper or manual caller can use.
func ReportarNegocio(
	ctx context.Context,
	writer buffer.BufferWriter,
	nombreUseCase string,
	comando map[string]any,
	err error,
	duracionMs float64,
) {
	traceID := middleware.GetTraceIDFromCtx(ctx)
	spanID := middleware.GetSpanIDFromCtx(ctx)
	userID := presentation_middleware.GetUsuarioIDFromCtx(ctx)
	result := ClassifyResult(err)
	level := DetermineLevel(result)

	payload := telemetry.LogPayload{
		LogType:     "NEGOCIO",
		Level:       level,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		TraceID:     traceID,
		SpanID:      spanID,
		ServiceName: "identidad",
		Environment: "dev",
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

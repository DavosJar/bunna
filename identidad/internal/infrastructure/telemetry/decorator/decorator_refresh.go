package decorator

import (
	"context"
	"encoding/json"
	"time"

	uc_refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/middleware"
	presentation_middleware "github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
)

type decoratorRefresh struct {
	uc     RefreshUseCase
	writer buffer.BufferWriter
}

func NewDecoratorRefresh(uc RefreshUseCase, writer buffer.BufferWriter) *decoratorRefresh {
	return &decoratorRefresh{uc: uc, writer: writer}
}

func (d *decoratorRefresh) Ejecutar(ctx context.Context, cmd uc_refresh.ComandoRenovarSesion) (*uc_refresh.RespuestaRenovarSesion, error) {
	start := time.Now()
	safeCmd := safeCommand(cmd)
	traceID := middleware.GetTraceIDFromCtx(ctx)
	userID := presentation_middleware.GetUsuarioIDFromCtx(ctx)
	resp, err := d.uc.Ejecutar(ctx, cmd)
	duration := time.Since(start).Milliseconds()
	result := classifyResult(err)
	level := determineLevel(result)
	out := map[string]any{
		"log_type":            "NEGOCIO",
		"use_case":            "Refresh",
		"command":             safeCmd,
		"result":              result,
		"user_id":             userID,
		"details":             map[string]any{},
		"duration_usecase_ms": duration,
		"trace_id":            traceID,
	}
	data, _ := json.Marshal(out)
	prio := buffer.Media
	if level == "ERROR" {
		prio = buffer.Alta
	}
	_ = d.writer.Write(data, prio)
	return resp, err
}

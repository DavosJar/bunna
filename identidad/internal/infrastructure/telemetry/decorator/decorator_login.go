package decorator

import (
	"context"
	"encoding/json"
	"time"

	uc_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/middleware"
	presentation_middleware "github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
)

type decoratorLogin struct {
	uc     LoginUseCase
	writer buffer.BufferWriter
}

func NewDecoratorLogin(uc LoginUseCase, writer buffer.BufferWriter) *decoratorLogin {
	return &decoratorLogin{uc: uc, writer: writer}
}

func (d *decoratorLogin) Ejecutar(ctx context.Context, cmd uc_login.ComandoIniciarSesion) (*uc_login.RespuestaIniciarSesion, error) {
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
		"use_case":            "Login",
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

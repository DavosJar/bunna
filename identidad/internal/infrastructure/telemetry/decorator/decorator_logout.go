package decorator

import (
	"context"
	"time"

	uc_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
)

type decoratorLogout struct {
	uc     LogoutUseCase
	writer buffer.BufferWriter
}

func NewDecoratorLogout(uc LogoutUseCase, writer buffer.BufferWriter) *decoratorLogout {
	return &decoratorLogout{uc: uc, writer: writer}
}

func (d *decoratorLogout) Ejecutar(ctx context.Context, cmd uc_logout.ComandoCerrarSesion) (*uc_logout.RespuestaCerrarSesion, error) {
	start := time.Now()
	safeCmd := SafeCommand(cmd)
	resp, err := d.uc.Ejecutar(ctx, cmd)
	duration := time.Since(start).Milliseconds()
	ReportarNegocio(ctx, d.writer, "Logout", safeCmd, err, float64(duration))
	return resp, err
}

func (d *decoratorLogout) CerrarTodas(ctx context.Context, cmd uc_logout.ComandoCerrarTodasLasSesiones) (*uc_logout.RespuestaCerrarSesion, error) {
	start := time.Now()
	safeCmd := SafeCommand(cmd)
	resp, err := d.uc.CerrarTodas(ctx, cmd)
	duration := time.Since(start).Milliseconds()
	ReportarNegocio(ctx, d.writer, "CerrarTodas", safeCmd, err, float64(duration))
	return resp, err
}

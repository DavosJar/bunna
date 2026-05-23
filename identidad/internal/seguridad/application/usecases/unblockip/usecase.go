package unblockip

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

type DesbloquearIPCasoDeUso struct {
	intentoRepo seguridad.IntentoIPRepositorio
	authSvc     rbac.AuthorizationService
}

func NewDesbloquearIPCasoDeUso(
	intentoRepo seguridad.IntentoIPRepositorio,
	authSvc rbac.AuthorizationService,
) *DesbloquearIPCasoDeUso {
	return &DesbloquearIPCasoDeUso{intentoRepo: intentoRepo, authSvc: authSvc}
}

func (uc *DesbloquearIPCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoDesbloquearIP) (*RespuestaDesbloquearIP, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoIPDesbloquear)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	intento, err := uc.intentoRepo.ObtenerPorIP(ctx, cmd.IP)
	if err != nil {
		return nil, fmt.Errorf("IP no encontrada: %w", err)
	}

	if !intento.EstaBloqueada(time.Now()) {
		return nil, fmt.Errorf("la IP %s no está bloqueada", cmd.IP)
	}

	intento.Desbloquear()
	if _, err := uc.intentoRepo.Actualizar(ctx, intento); err != nil {
		return nil, fmt.Errorf("error al desbloquear IP: %w", err)
	}

	ahora := time.Now().Format("2006-01-02T15:04:05Z")
	return &RespuestaDesbloquearIP{
		IP:             cmd.IP,
		DesbloqueadoEn: ahora,
	}, nil
}

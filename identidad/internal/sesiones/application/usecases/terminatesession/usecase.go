package terminatesession

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	sesiones "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

type ForzarCierreSesionCasoDeUso struct {
	sessionRepo sesiones.SesionRepositorio
	authSvc     rbac.AuthorizationService
}

func NewForzarCierreSesionCasoDeUso(
	sessionRepo sesiones.SesionRepositorio,
	authSvc rbac.AuthorizationService,
) *ForzarCierreSesionCasoDeUso {
	return &ForzarCierreSesionCasoDeUso{sessionRepo: sessionRepo, authSvc: authSvc}
}

func (uc *ForzarCierreSesionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoForzarCierreSesion) (*RespuestaForzarCierreSesion, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoSesionForzarCierre)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	sesion, err := uc.sessionRepo.ObtenerPorID(ctx, cmd.SesionID)
	if err != nil {
		return nil, fmt.Errorf("sesión no encontrada: %w", err)
	}

	sesion.Revocar()
	_, err = uc.sessionRepo.Actualizar(ctx, sesion)
	if err != nil {
		return nil, fmt.Errorf("error al persistir revocación: %w", err)
	}

	return &RespuestaForzarCierreSesion{
		SesionID:   sesion.ID(),
		Estado:     string(sesion.Estado()),
		RevocadoEn: time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

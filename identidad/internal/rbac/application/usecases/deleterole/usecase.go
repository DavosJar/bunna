package deleterole

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type EliminarRolCasoDeUso struct {
	rolRepo rbac.RolRepositorio
	authSvc rbac.AuthorizationService
}

func NewEliminarRolCasoDeUso(
	rolRepo rbac.RolRepositorio,
	authSvc rbac.AuthorizationService,
) *EliminarRolCasoDeUso {
	return &EliminarRolCasoDeUso{rolRepo: rolRepo, authSvc: authSvc}
}

func (uc *EliminarRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoEliminarRol) (*RespuestaEliminarRol, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoRolEliminar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	rol, err := uc.rolRepo.ObtenerPorID(ctx, cmd.RolID)
	if err != nil {
		return nil, fmt.Errorf("rol no encontrado: %w", err)
	}

	if rol.EsSistema {
		return nil, rbac.ErrRolInmutable
	}

	if err := uc.rolRepo.Eliminar(ctx, cmd.RolID); err != nil {
		return nil, fmt.Errorf("error al eliminar rol: %w", err)
	}

	return &RespuestaEliminarRol{
		RolID:       cmd.RolID,
		EliminadoEn: time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

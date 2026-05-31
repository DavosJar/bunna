package updaterole

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type ModificarRolCasoDeUso struct {
	rolRepo rbac.RolRepositorio
	authSvc rbac.AuthorizationService
}

func NewModificarRolCasoDeUso(
	rolRepo rbac.RolRepositorio,
	authSvc rbac.AuthorizationService,
) *ModificarRolCasoDeUso {
	return &ModificarRolCasoDeUso{rolRepo: rolRepo, authSvc: authSvc}
}

func (uc *ModificarRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoModificarRol) (*RespuestaModificarRol, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoRolModificar)
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

	if cmd.Descripcion != "" {
		if err := uc.rolRepo.ActualizarDescripcion(ctx, cmd.RolID, cmd.Descripcion); err != nil {
			return nil, fmt.Errorf("error al actualizar descripción: %w", err)
		}
	}

	return &RespuestaModificarRol{
		ID:           cmd.RolID,
		Nombre:       cmd.Nombre,
		Descripcion:  cmd.Descripcion,
		ModificadoEn: time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

package assignpermissiontorole

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type AsignarPermisoARolCasoDeUso struct {
	rolRepo        rbac.RolRepositorio
	permisoRepo    rbac.PermisoRepositorio
	rolPermisoRepo rbac.RolPermisoRepositorio
	authSvc        rbac.AuthorizationService
}

func NewAsignarPermisoARolCasoDeUso(
	rolRepo rbac.RolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
	rolPermisoRepo rbac.RolPermisoRepositorio,
	authSvc rbac.AuthorizationService,
) *AsignarPermisoARolCasoDeUso {
	return &AsignarPermisoARolCasoDeUso{
		rolRepo:        rolRepo,
		permisoRepo:    permisoRepo,
		rolPermisoRepo: rolPermisoRepo,
		authSvc:        authSvc,
	}
}

func (uc *AsignarPermisoARolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoAsignarPermisoARol) (*RespuestaAsignarPermisoARol, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoRolPermisoAsignar)
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

	permisoDB, err := uc.permisoRepo.ObtenerPorCodigo(ctx, cmd.PermisoCodigo)
	if err != nil {
		return nil, fmt.Errorf("permiso no encontrado: %w", err)
	}

	if err := uc.rolPermisoRepo.AsignarPermiso(ctx, cmd.RolID, permisoDB.ID); err != nil {
		return nil, fmt.Errorf("error al asignar permiso al rol: %w", err)
	}

	return &RespuestaAsignarPermisoARol{
		RolID:         cmd.RolID,
		PermisoCodigo: cmd.PermisoCodigo,
		AsignadoEn:    time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

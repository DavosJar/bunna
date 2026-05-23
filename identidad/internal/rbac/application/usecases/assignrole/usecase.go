package assignrole

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type AsignarRolCasoDeUso struct {
	usuarioRolRepo     rbac.UsuarioRolRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
	rolRepo            rbac.RolRepositorio
	authSvc            rbac.AuthorizationService
}

func NewAsignarRolCasoDeUso(
	usuarioRolRepo rbac.UsuarioRolRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
	rolRepo rbac.RolRepositorio,
	authSvc rbac.AuthorizationService,
) *AsignarRolCasoDeUso {
	return &AsignarRolCasoDeUso{
		usuarioRolRepo:     usuarioRolRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
		rolRepo:            rolRepo,
		authSvc:            authSvc,
	}
}

func (uc *AsignarRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoAsignarRol) (*RespuestaAsignarRol, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoRolAsignar)
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

	if rol.Nombre == rbac.RolSysAdmin && cmd.TenantID != "" {
		return nil, rbac.ErrSysAdminRequiereTenantVacio
	}

	if cmd.TenantID == "" {
		if err := uc.usuarioRolRepo.Crear(ctx, cmd.UsuarioID, cmd.RolID); err != nil {
			return nil, fmt.Errorf("error al asignar rol global: %w", err)
		}
	} else {
		if err := uc.usuarioTenantRolRepo.Crear(ctx, cmd.UsuarioID, cmd.TenantID, cmd.RolID); err != nil {
			return nil, fmt.Errorf("error al asignar rol en tenant: %w", err)
		}
	}

	return &RespuestaAsignarRol{
		UsuarioID:  cmd.UsuarioID,
		RolID:      cmd.RolID,
		TenantID:   cmd.TenantID,
		AsignadoEn: time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

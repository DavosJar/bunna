package revokerole

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type RevocarRolCasoDeUso struct {
	usuarioRolRepo     rbac.UsuarioRolRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
	authSvc            rbac.AuthorizationService
}

func NewRevocarRolCasoDeUso(
	usuarioRolRepo rbac.UsuarioRolRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
	authSvc rbac.AuthorizationService,
) *RevocarRolCasoDeUso {
	return &RevocarRolCasoDeUso{
		usuarioRolRepo:     usuarioRolRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
		authSvc:            authSvc,
	}
}

func (uc *RevocarRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoRevocarRol) (*RespuestaRevocarRol, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoRolRevocar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	if cmd.TenantID == "" {
		if err := uc.usuarioRolRepo.Eliminar(ctx, cmd.UsuarioID, cmd.RolID); err != nil {
			return nil, fmt.Errorf("error al revocar rol global: %w", err)
		}
	} else {
		if err := uc.usuarioTenantRolRepo.Eliminar(ctx, cmd.UsuarioID, cmd.TenantID, cmd.RolID); err != nil {
			return nil, fmt.Errorf("error al revocar rol en tenant: %w", err)
		}
	}

	return &RespuestaRevocarRol{
		UsuarioID:  cmd.UsuarioID,
		RolID:      cmd.RolID,
		TenantID:   cmd.TenantID,
		RevocadoEn: time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

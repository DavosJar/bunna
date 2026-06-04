package revokepermissionfromrole

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type RevocarPermisoDeRolCasoDeUso struct {
	rolRepo        rbac.RolRepositorio
	permisoRepo    rbac.PermisoRepositorio
	rolPermisoRepo rbac.RolPermisoRepositorio
	authSvc        rbac.AuthorizationService
}

func NewRevocarPermisoDeRolCasoDeUso(
	rolRepo rbac.RolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
	rolPermisoRepo rbac.RolPermisoRepositorio,
	authSvc rbac.AuthorizationService,
) *RevocarPermisoDeRolCasoDeUso {
	return &RevocarPermisoDeRolCasoDeUso{
		rolRepo:        rolRepo,
		permisoRepo:    permisoRepo,
		rolPermisoRepo: rolPermisoRepo,
		authSvc:        authSvc,
	}
}

func (uc *RevocarPermisoDeRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoRevocarPermisoDeRol) (*RespuestaRevocarPermisoDeRol, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoRolPermisoRevocar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	// Also verify the executor possesses the specific permission they're trying to revoke
	ok2, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, cmd.PermisoCodigo)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso específico: %w", err)
	}
	if !ok2 {
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

	if err := uc.rolPermisoRepo.EliminarPermiso(ctx, cmd.RolID, permisoDB.ID, cmd.TenantID); err != nil {
		return nil, fmt.Errorf("error al revocar permiso del rol: %w", err)
	}

	return &RespuestaRevocarPermisoDeRol{
		RolID:         cmd.RolID,
		PermisoCodigo: cmd.PermisoCodigo,
		RevocadoEn:    time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

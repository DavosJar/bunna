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
	publisher      rbac.RolPublisher
}

func NewRevocarPermisoDeRolCasoDeUso(
	rolRepo rbac.RolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
	rolPermisoRepo rbac.RolPermisoRepositorio,
	authSvc rbac.AuthorizationService,
	publisher rbac.RolPublisher,
) *RevocarPermisoDeRolCasoDeUso {
	return &RevocarPermisoDeRolCasoDeUso{
		rolRepo:        rolRepo,
		permisoRepo:    permisoRepo,
		rolPermisoRepo: rolPermisoRepo,
		authSvc:        authSvc,
		publisher:      publisher,
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

	// Publish role update event
	permisosDB, err := uc.rolPermisoRepo.ListarPorRolYTenant(ctx, cmd.RolID, cmd.TenantID)
	if err == nil {
		var codigos []string
		for _, p := range permisosDB {
			codigos = append(codigos, p.Codigo)
		}
		_ = uc.publisher.PublicarRolActualizado(ctx, cmd.RolID, cmd.TenantID, codigos)
	}

	return &RespuestaRevocarPermisoDeRol{
		RolID:         cmd.RolID,
		PermisoCodigo: cmd.PermisoCodigo,
		RevocadoEn:    time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

package checkpermission

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

// VerificarPermisoCasoDeUso verifica si un usuario tiene un permiso específico.
// Si tenantID no está vacío, revisa primero roles del tenant y luego roles globales.
// Si tenantID está vacío, solo revisa roles globales.
type VerificarPermisoCasoDeUso struct {
	usuarioRolRepo       rbac.UsuarioRolRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
	permisoRepo          rbac.PermisoRepositorio
}

func NewVerificarPermisoCasoDeUso(
	usuarioRolRepo rbac.UsuarioRolRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
) *VerificarPermisoCasoDeUso {
	return &VerificarPermisoCasoDeUso{
		usuarioRolRepo:       usuarioRolRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
		permisoRepo:          permisoRepo,
	}
}

// TienePermiso satisface rbac.AuthorizationService.
// Evalúa permisos en este orden:
//  1. Roles del tenant (si tenantID != "")
//  2. Roles globales (sys_admin)
func (uc *VerificarPermisoCasoDeUso) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	if tenantID != "" {
		roles, err := uc.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, usuarioID, tenantID)
		if err != nil {
			return false, err
		}
		for _, rol := range roles {
			if rbac.TienePermisoEnRol(ctx, uc.permisoRepo, rol, codigoPermiso, tenantID) {
				return true, nil
			}
		}
	}

	roles, err := uc.usuarioRolRepo.ListarRolesPorUsuario(ctx, usuarioID)
	if err != nil {
		return false, err
	}
	for _, rol := range roles {
		if rbac.TienePermisoEnRol(ctx, uc.permisoRepo, rol, codigoPermiso, tenantID) {
			return true, nil
		}
	}

	return false, nil
}

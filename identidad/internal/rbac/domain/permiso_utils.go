package rbac

import (
	"context"
	"sync"
)

var (
	initPermisosSistemaOnce sync.Once
	permisosSistemaMap      map[string]map[string]bool // rol.Nombre → codigoPermiso → true
)

func permisosSistemaInit() {
	initPermisosSistemaOnce.Do(func() {
		permisosSistemaMap = make(map[string]map[string]bool, len(RolesDeSistema))
		for _, r := range RolesDeSistema {
			perms := make(map[string]bool, len(r.Permisos))
			for _, p := range r.Permisos {
				perms[p] = true
			}
			permisosSistemaMap[r.Nombre] = perms
		}
	})
}

// TienePermisoEnRol verifica si un rol tiene un permiso específico.
// Para roles de sistema usa un mapa en memoria (O(1)).
// Para roles personalizados consulta el repositorio.
func TienePermisoEnRol(ctx context.Context, permisoRepo PermisoRepositorio, rol *RolDB, codigoPermiso string, tenantID string) bool {
	if rol.EsSistema {
		permisosSistemaInit()
		if perms, ok := permisosSistemaMap[rol.Nombre]; ok && len(perms) > 0 {
			return perms[codigoPermiso]
		}
		// System roles with empty permissions = dynamic, fall through to DB
	}

	permisos, err := permisoRepo.ListarPorRol(ctx, rol.ID, tenantID)
	if err != nil {
		return false
	}
	for _, p := range permisos {
		if p.Codigo == codigoPermiso {
			return true
		}
	}
	return false
}

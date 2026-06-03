package autorizacion

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type ServicioAutorizacion struct {
	usuarioRolRepo rbac.UsuarioRolRepositorio
	permisoRepo    rbac.PermisoRepositorio
}

func NewServicioAutorizacion(
	usuarioRolRepo rbac.UsuarioRolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
) *ServicioAutorizacion {
	return &ServicioAutorizacion{
		usuarioRolRepo: usuarioRolRepo,
		permisoRepo:    permisoRepo,
	}
}

func (s *ServicioAutorizacion) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	roles, err := s.usuarioRolRepo.ListarRolesPorUsuario(ctx, usuarioID)
	if err != nil {
		return false, err
	}

	for _, rol := range roles {
		if rbac.TienePermisoEnRol(ctx, s.permisoRepo, rol, codigoPermiso, tenantID) {
			return true, nil
		}
	}

	return false, nil
}

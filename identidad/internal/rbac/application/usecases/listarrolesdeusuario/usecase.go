package listarrolesdeusuario

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type ListarRolesDeUsuarioCasoDeUso struct {
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
}

func NewListarRolesDeUsuarioCasoDeUso(
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
) *ListarRolesDeUsuarioCasoDeUso {
	return &ListarRolesDeUsuarioCasoDeUso{
		usuarioTenantRolRepo: usuarioTenantRolRepo,
	}
}

type ComandoListarRolesDeUsuario struct {
	UsuarioID string
	TenantID  string
}

type RolDeUsuarioDTO struct {
	RolID  string `json:"rol_id"`
	Nombre string `json:"nombre"`
}

type RespuestaListarRolesDeUsuario struct {
	Roles []RolDeUsuarioDTO `json:"roles"`
}

func (uc *ListarRolesDeUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoListarRolesDeUsuario) (*RespuestaListarRolesDeUsuario, error) {
	dbRoles, err := uc.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, cmd.UsuarioID, cmd.TenantID)
	if err != nil {
		return nil, fmt.Errorf("error al listar roles del usuario: %w", err)
	}

	roles := make([]RolDeUsuarioDTO, len(dbRoles))
	for i, r := range dbRoles {
		roles[i] = RolDeUsuarioDTO{RolID: r.ID, Nombre: r.Nombre}
	}

	return &RespuestaListarRolesDeUsuario{Roles: roles}, nil
}

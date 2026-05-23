package application

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	tenantdomain "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

// AutorizacionServicio verifica permisos y genera claims JWT
type AutorizacionServicio struct {
	rolRepo            rbac.RolRepositorio
	permisoRepo        rbac.PermisoRepositorio
	usuarioRolRepo     rbac.UsuarioRolRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
	tenantRepo         tenantdomain.TenantRepositorio
}

func NuevoAutorizacionServicio(
	rolRepo rbac.RolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
	usuarioRolRepo rbac.UsuarioRolRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
	tenantRepo tenantdomain.TenantRepositorio,
) *AutorizacionServicio {
	return &AutorizacionServicio{
		rolRepo:              rolRepo,
		permisoRepo:          permisoRepo,
		usuarioRolRepo:       usuarioRolRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
		tenantRepo:           tenantRepo,
	}
}

// TienePermiso verifica si un usuario tiene un permiso en un contexto de tenant
func (s *AutorizacionServicio) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	// 1. Verificar si es SYS_ADMIN (global)
	esSysAdmin, err := s.usuarioRolRepo.TieneRol(ctx, usuarioID, rbac.RolSysAdmin)
	if err != nil {
		return false, err
	}
	if esSysAdmin {
		return true, nil
	}

	// 2. Sin tenant y no es SYS_ADMIN → denegar
	if tenantID == "" {
		return false, nil
	}

	// 3. Obtener roles del usuario en el tenant
	roles, err := s.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, usuarioID, tenantID)
	if err != nil {
		return false, err
	}

	// 4. Verificar permisos de cada rol
	for _, rol := range roles {
		permisos, err := s.permisoRepo.ListarPorRol(ctx, rol.ID)
		if err != nil {
			return false, err
		}
		for _, p := range permisos {
			if p.Codigo == codigoPermiso {
				return true, nil
			}
		}
	}

	return false, nil
}

// ObtenerClaimsUsuario construye los claims JWT para un usuario
func (s *AutorizacionServicio) ObtenerClaimsUsuario(ctx context.Context, usuarioID string) (*rbac.UsuarioClaims, error) {
	claims := &rbac.UsuarioClaims{
		UsuarioID: usuarioID,
		Tenants:   make(map[string]rbac.TenantClaims),
	}

	// 1. Verificar si es SYS_ADMIN
	esSysAdmin, err := s.usuarioRolRepo.TieneRol(ctx, usuarioID, rbac.RolSysAdmin)
	if err != nil {
		return nil, err
	}
	if esSysAdmin {
		claims.Global = true
		return claims, nil
	}

	// 2. Cargar tenants del usuario
	tenants, err := s.tenantRepo.ListarPorUsuario(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	// 3. Para cada tenant, cargar roles y permisos
	for _, t := range tenants {
		roles, err := s.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, usuarioID, t.ID())
		if err != nil {
			return nil, err
		}

		nombresRoles := make([]string, 0)
		codigosPermisos := make([]string, 0)
		permisosVistos := make(map[string]bool)

		for _, rol := range roles {
			nombresRoles = append(nombresRoles, rol.Nombre)
			permisos, err := s.permisoRepo.ListarPorRol(ctx, rol.ID)
			if err != nil {
				return nil, err
			}
			for _, p := range permisos {
				if !permisosVistos[p.Codigo] {
					codigosPermisos = append(codigosPermisos, p.Codigo)
					permisosVistos[p.Codigo] = true
				}
			}
		}

		claims.Tenants[t.ID()] = rbac.TenantClaims{
			Slug:     t.Slug(),
			Roles:    nombresRoles,
			Permisos: codigosPermisos,
		}
	}

	return claims, nil
}
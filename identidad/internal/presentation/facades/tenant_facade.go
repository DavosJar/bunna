package facades

import (
	"context"
	"sort"

	decorator "github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/decorator"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	uc_listarmistenants "github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/uc_listarmistenants"
	tenant_domain "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

// TenantConRol representa un tenant con el rol del usuario en él.
type TenantConRol struct {
	ID       string
	Nombre   string
	Slug     string
	Rol      string
	EsPropio bool
}

// RespuestaListarMisTenants contiene la lista de tenants y el ID del tenant propio.
type RespuestaListarMisTenants struct {
	Tenants  []TenantConRol
	PropioID string
}

type TenantFacade interface {
	ListarMisTenants(ctx context.Context, usuarioID string) (*RespuestaListarMisTenants, error)
}

type tenantFacadeImpl struct {
	listarMisTenantsUC   decorator.UseCase[string, *uc_listarmistenants.RespuestaListarMisTenants]
	membresiaRepo        tenant_domain.MembresiaRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
}

func NewTenantFacade(
	listarMisTenantsUC decorator.UseCase[string, *uc_listarmistenants.RespuestaListarMisTenants],
	membresiaRepo tenant_domain.MembresiaRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
) TenantFacade {
	return &tenantFacadeImpl{
		listarMisTenantsUC:   listarMisTenantsUC,
		membresiaRepo:        membresiaRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
	}
}

// ListarMisTenants retorna los tenants del usuario con su rol, e identifica el tenant propio.
func (f *tenantFacadeImpl) ListarMisTenants(ctx context.Context, usuarioID string) (*RespuestaListarMisTenants, error) {
	resp, err := f.listarMisTenantsUC.Ejecutar(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	sort.Slice(resp.Tenants, func(i, j int) bool {
		return resp.Tenants[i].FechaCreacion.Before(resp.Tenants[j].FechaCreacion)
	})

	propioID := ""
	if len(resp.Tenants) > 0 {
		propioID = resp.Tenants[0].ID
	}

	var result []TenantConRol
	for _, t := range resp.Tenants {
		rol := ""
		roles, err := f.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, usuarioID, t.ID)
		if err == nil && len(roles) > 0 {
			rol = roles[0].Nombre
		}

		result = append(result, TenantConRol{
			ID:       t.ID,
			Nombre:   t.Nombre,
			Slug:     t.Slug,
			Rol:      rol,
			EsPropio: t.ID == propioID,
		})
	}

	return &RespuestaListarMisTenants{
		Tenants:  result,
		PropioID: propioID,
	}, nil
}

package facades

import (
	"context"
	"sort"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	uc_updatetenant "github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/updatetenant"
	tenant_domain "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type ComandoConfigurarTenant struct {
	TenantID   string
	Nombre     string
	Slug       string
	EjecutorID string
}

type RespuestaConfigurarTenant struct {
	TenantID     string
	Nombre       string
	Slug         string
	ModificadoEn string
}

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
	ConfigurarTenant(ctx context.Context, cmd ComandoConfigurarTenant) (*RespuestaConfigurarTenant, error)
	ListarMisTenants(ctx context.Context, usuarioID string) (*RespuestaListarMisTenants, error)
}

type tenantFacadeImpl struct {
	configurarTenant    *uc_updatetenant.ConfigurarTenantCasoDeUso
	tenantRepo          tenant_domain.TenantRepositorio
	membresiaRepo       tenant_domain.MembresiaRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
}

func NewTenantFacade(
	configurarTenant *uc_updatetenant.ConfigurarTenantCasoDeUso,
	tenantRepo tenant_domain.TenantRepositorio,
	membresiaRepo tenant_domain.MembresiaRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
) TenantFacade {
	return &tenantFacadeImpl{
		configurarTenant:    configurarTenant,
		tenantRepo:          tenantRepo,
		membresiaRepo:       membresiaRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
	}
}

func (f *tenantFacadeImpl) ConfigurarTenant(ctx context.Context, cmd ComandoConfigurarTenant) (*RespuestaConfigurarTenant, error) {
	resp, err := f.configurarTenant.Ejecutar(ctx, &uc_updatetenant.ComandoConfigurarTenant{
		TenantID:   cmd.TenantID,
		Nombre:     cmd.Nombre,
		Slug:       cmd.Slug,
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaConfigurarTenant{
		TenantID:     resp.TenantID,
		Nombre:       resp.Nombre,
		Slug:         resp.Slug,
		ModificadoEn: resp.ModificadoEn,
	}, nil
}

// ListarMisTenants retorna los tenants del usuario con su rol, e identifica el tenant propio.
func (f *tenantFacadeImpl) ListarMisTenants(ctx context.Context, usuarioID string) (*RespuestaListarMisTenants, error) {
	tenants, err := f.tenantRepo.ListarPorUsuario(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	// Ordenar por fecha de creación para determinar el tenant propio (el primero)
	sort.Slice(tenants, func(i, j int) bool {
		return tenants[i].FechaCreacion().Before(tenants[j].FechaCreacion())
	})

	propioID := ""
	if len(tenants) > 0 {
		propioID = tenants[0].ID()
	}

	var result []TenantConRol
	for _, t := range tenants {
		rol := ""
		roles, err := f.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, usuarioID, t.ID())
		if err == nil && len(roles) > 0 {
			rol = roles[0].Nombre
		}

		result = append(result, TenantConRol{
			ID:       t.ID(),
			Nombre:   t.Nombre(),
			Slug:     t.Slug(),
			Rol:      rol,
			EsPropio: t.ID() == propioID,
		})
	}

	return &RespuestaListarMisTenants{
		Tenants:  result,
		PropioID: propioID,
	}, nil
}

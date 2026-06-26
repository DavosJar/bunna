package uc_listarmistenants_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/uc_listarmistenants"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type mockTenantRepoMisTenants struct {
	listarPorUsuarioFunc func(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error)
}

func (m *mockTenantRepoMisTenants) Crear(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	return t, nil
}
func (m *mockTenantRepoMisTenants) ObtenerPorID(ctx context.Context, id string) (*tenant.Tenant, error) {
	return nil, tenant.ErrTenantNoEncontrado
}
func (m *mockTenantRepoMisTenants) ObtenerPorSlug(ctx context.Context, slug string) (*tenant.Tenant, error) {
	return nil, tenant.ErrTenantNoEncontrado
}
func (m *mockTenantRepoMisTenants) Actualizar(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	return t, nil
}
func (m *mockTenantRepoMisTenants) Listar(ctx context.Context) ([]*tenant.Tenant, error) {
	return nil, nil
}
func (m *mockTenantRepoMisTenants) ListarPorUsuario(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) {
	return m.listarPorUsuarioFunc(ctx, usuarioID)
}

func nuevoTenant(id, nombre, slug string) *tenant.Tenant {
	t, err := tenant.NuevoTenant(id, nombre, slug)
	if err != nil {
		panic(err)
	}
	return t
}

func TestListarMisTenantsExitoso(t *testing.T) {
	repo := &mockTenantRepoMisTenants{
		listarPorUsuarioFunc: func(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) {
			if usuarioID != "user-001" {
				t.Errorf("usuarioID esperado user-001, got %s", usuarioID)
			}
			return []*tenant.Tenant{
				nuevoTenant("t-001", "Mi Tenant", "mi-tenant"),
				nuevoTenant("t-002", "Otro Tenant", "otro-tenant"),
			}, nil
		},
	}
	uc := uc_listarmistenants.NewListarMisTenantsCasoDeUso(repo)

	resp, err := uc.Ejecutar(context.Background(), "user-001")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(resp.Tenants) != 2 {
		t.Fatalf("esperaba 2 tenants, got %d", len(resp.Tenants))
	}
	if resp.Tenants[0].Nombre != "Mi Tenant" {
		t.Errorf("Nombre[0] esperado 'Mi Tenant', got %s", resp.Tenants[0].Nombre)
	}
	if resp.Tenants[1].Slug != "otro-tenant" {
		t.Errorf("Slug[1] esperado 'otro-tenant', got %s", resp.Tenants[1].Slug)
	}
}

func TestListarMisTenantsVacio(t *testing.T) {
	repo := &mockTenantRepoMisTenants{
		listarPorUsuarioFunc: func(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) {
			return []*tenant.Tenant{}, nil
		},
	}
	uc := uc_listarmistenants.NewListarMisTenantsCasoDeUso(repo)

	resp, err := uc.Ejecutar(context.Background(), "user-002")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(resp.Tenants) != 0 {
		t.Fatalf("esperaba 0 tenants, got %d", len(resp.Tenants))
	}
}

func TestListarMisTenantsRepoError(t *testing.T) {
	repo := &mockTenantRepoMisTenants{
		listarPorUsuarioFunc: func(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) {
			return nil, errors.New("error de BD")
		},
	}
	uc := uc_listarmistenants.NewListarMisTenantsCasoDeUso(repo)

	_, err := uc.Ejecutar(context.Background(), "user-001")
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

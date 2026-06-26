package uc_obtenertenantporslug_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/uc_obtenertenantporslug"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type mockTenantRepoGetBySlug struct {
	obtenerPorSlugFunc func(ctx context.Context, slug string) (*tenant.Tenant, error)
}

func (m *mockTenantRepoGetBySlug) Crear(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	return t, nil
}
func (m *mockTenantRepoGetBySlug) ObtenerPorID(ctx context.Context, id string) (*tenant.Tenant, error) {
	return nil, tenant.ErrTenantNoEncontrado
}
func (m *mockTenantRepoGetBySlug) ObtenerPorSlug(ctx context.Context, slug string) (*tenant.Tenant, error) {
	return m.obtenerPorSlugFunc(ctx, slug)
}
func (m *mockTenantRepoGetBySlug) Actualizar(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	return t, nil
}
func (m *mockTenantRepoGetBySlug) Listar(ctx context.Context) ([]*tenant.Tenant, error) {
	return nil, nil
}
func (m *mockTenantRepoGetBySlug) ListarPorUsuario(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) {
	return nil, nil
}

func nuevoTenant(id, nombre, slug string) *tenant.Tenant {
	t, err := tenant.NuevoTenant(id, nombre, slug)
	if err != nil {
		panic(err)
	}
	return t
}

func TestObtenerTenantPorSlugExitoso(t *testing.T) {
	repo := &mockTenantRepoGetBySlug{
		obtenerPorSlugFunc: func(ctx context.Context, slug string) (*tenant.Tenant, error) {
			if slug != "mi-tenant" {
				t.Errorf("slug esperado 'mi-tenant', got %s", slug)
			}
			return nuevoTenant("t-001", "Mi Tenant", "mi-tenant"), nil
		},
	}
	uc := uc_obtenertenantporslug.NewObtenerTenantPorSlugCasoDeUso(repo)

	resp, err := uc.Ejecutar(context.Background(), "mi-tenant")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.ID != "t-001" {
		t.Errorf("ID esperado t-001, got %s", resp.ID)
	}
	if resp.Slug != "mi-tenant" {
		t.Errorf("Slug esperado 'mi-tenant', got %s", resp.Slug)
	}
}

func TestObtenerTenantPorSlugNoEncontrado(t *testing.T) {
	repo := &mockTenantRepoGetBySlug{
		obtenerPorSlugFunc: func(ctx context.Context, slug string) (*tenant.Tenant, error) {
			return nil, tenant.ErrTenantNoEncontrado
		},
	}
	uc := uc_obtenertenantporslug.NewObtenerTenantPorSlugCasoDeUso(repo)

	_, err := uc.Ejecutar(context.Background(), "no-existe")
	if !errors.Is(err, tenant.ErrTenantNoEncontrado) {
		t.Errorf("esperaba ErrTenantNoEncontrado, got %v", err)
	}
}

func TestObtenerTenantPorSlugRepoError(t *testing.T) {
	repo := &mockTenantRepoGetBySlug{
		obtenerPorSlugFunc: func(ctx context.Context, slug string) (*tenant.Tenant, error) {
			return nil, errors.New("error de BD")
		},
	}
	uc := uc_obtenertenantporslug.NewObtenerTenantPorSlugCasoDeUso(repo)

	_, err := uc.Ejecutar(context.Background(), "mi-tenant")
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

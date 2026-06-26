package uc_obtenertenantporid_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/uc_obtenertenantporid"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type mockTenantRepoGetByID struct {
	obtenerPorIDFunc func(ctx context.Context, id string) (*tenant.Tenant, error)
}

func (m *mockTenantRepoGetByID) Crear(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	return t, nil
}
func (m *mockTenantRepoGetByID) ObtenerPorID(ctx context.Context, id string) (*tenant.Tenant, error) {
	return m.obtenerPorIDFunc(ctx, id)
}
func (m *mockTenantRepoGetByID) ObtenerPorSlug(ctx context.Context, slug string) (*tenant.Tenant, error) {
	return nil, tenant.ErrTenantNoEncontrado
}
func (m *mockTenantRepoGetByID) Actualizar(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	return t, nil
}
func (m *mockTenantRepoGetByID) Listar(ctx context.Context) ([]*tenant.Tenant, error) {
	return nil, nil
}
func (m *mockTenantRepoGetByID) ListarPorUsuario(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) {
	return nil, nil
}

func nuevoTenant(id, nombre, slug string) *tenant.Tenant {
	t, err := tenant.NuevoTenant(id, nombre, slug)
	if err != nil {
		panic(err)
	}
	return t
}

func TestObtenerTenantPorIDExitoso(t *testing.T) {
	repo := &mockTenantRepoGetByID{
		obtenerPorIDFunc: func(ctx context.Context, id string) (*tenant.Tenant, error) {
			if id != "t-001" {
				t.Errorf("id esperado t-001, got %s", id)
			}
			return nuevoTenant("t-001", "Mi Tenant", "mi-tenant"), nil
		},
	}
	uc := uc_obtenertenantporid.NewObtenerTenantPorIDCasoDeUso(repo)

	resp, err := uc.Ejecutar(context.Background(), "t-001")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.ID != "t-001" {
		t.Errorf("ID esperado t-001, got %s", resp.ID)
	}
	if resp.Nombre != "Mi Tenant" {
		t.Errorf("Nombre esperado 'Mi Tenant', got %s", resp.Nombre)
	}
	if resp.Slug != "mi-tenant" {
		t.Errorf("Slug esperado 'mi-tenant', got %s", resp.Slug)
	}
}

func TestObtenerTenantPorIDNoEncontrado(t *testing.T) {
	repo := &mockTenantRepoGetByID{
		obtenerPorIDFunc: func(ctx context.Context, id string) (*tenant.Tenant, error) {
			return nil, tenant.ErrTenantNoEncontrado
		},
	}
	uc := uc_obtenertenantporid.NewObtenerTenantPorIDCasoDeUso(repo)

	_, err := uc.Ejecutar(context.Background(), "t-999")
	if !errors.Is(err, tenant.ErrTenantNoEncontrado) {
		t.Errorf("esperaba ErrTenantNoEncontrado, got %v", err)
	}
}

func TestObtenerTenantPorIDRepoError(t *testing.T) {
	repo := &mockTenantRepoGetByID{
		obtenerPorIDFunc: func(ctx context.Context, id string) (*tenant.Tenant, error) {
			return nil, errors.New("error de BD")
		},
	}
	uc := uc_obtenertenantporid.NewObtenerTenantPorIDCasoDeUso(repo)

	_, err := uc.Ejecutar(context.Background(), "t-001")
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

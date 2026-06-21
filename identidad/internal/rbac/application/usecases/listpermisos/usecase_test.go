package listpermisos_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listpermisos"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type mockPermisoRepo struct {
	permisos []*rbac.PermisoDB
	err      error
}

func (m *mockPermisoRepo) ObtenerPorCodigo(ctx context.Context, codigo string) (*rbac.PermisoDB, error) { return nil, nil }
func (m *mockPermisoRepo) Listar(ctx context.Context) ([]*rbac.PermisoDB, error) { return m.permisos, m.err }
func (m *mockPermisoRepo) Crear(ctx context.Context, p *rbac.PermisoDB) error { return nil }
func (m *mockPermisoRepo) ActualizarNombreDescripcion(ctx context.Context, id, nombre, desc string) error { return nil }
func (m *mockPermisoRepo) ListarPorRol(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) { return nil, nil }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

func TestListarPermisosExitoso(t *testing.T) {
	repo := &mockPermisoRepo{permisos: []*rbac.PermisoDB{
		{Codigo: "test:1", Nombre: "Test1", Modulo: "mod"},
		{Codigo: "test:2", Nombre: "Test2", Modulo: "mod"},
	}}
	uc := listpermisos.NewListarPermisosCasoDeUso(repo, &mockAuthSvc{permiso: true})
	resp, err := uc.Ejecutar(context.Background(), "admin-1", "t-1")
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.Total != 2 { t.Errorf("Total incorrecto: %v", resp.Total) }
}

func TestListarPermisosSinPermiso(t *testing.T) {
	uc := listpermisos.NewListarPermisosCasoDeUso(&mockPermisoRepo{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), "e-1", "t-1")
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestListarPermisosErrorRepo(t *testing.T) {
	uc := listpermisos.NewListarPermisosCasoDeUso(
		&mockPermisoRepo{err: errors.New("fallo bd")},
		&mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), "e-1", "t-1")
	if err == nil { t.Fatal("esperaba error de repositorio") }
}

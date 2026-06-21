package listarmispermisos_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listarmispermisos"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockRolRepo struct {
	rol *rbac.RolDB
	err error
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) { return m.rol, m.err }
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Crear(ctx context.Context, r *rbac.RolDB) error { return nil }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, id, desc string) error { return nil }

type mockRolPermisoRepo struct {
	permisos []*rbac.PermisoDB
	err      error
}

func (m *mockRolPermisoRepo) AsignarPermiso(ctx context.Context, rolID, permisoID, tenantID, asignadoPor string) error { return nil }
func (m *mockRolPermisoRepo) EliminarPermiso(ctx context.Context, rolID, permisoID, tenantID string) error { return nil }
func (m *mockRolPermisoRepo) ListarPorRolYTenant(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) {
	return m.permisos, m.err
}

func TestListarMisPermisosExitoso(t *testing.T) {
	uc := listarmispermisos.NewListarMisPermisosCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", Nombre: "editor"}},
		&mockRolPermisoRepo{permisos: []*rbac.PermisoDB{
			{Codigo: "test:1", Nombre: "Permiso 1", Modulo: "mod"},
		}},
	)
	resp, err := uc.Ejecutar(context.Background(), "editor", "t-1")
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if len(resp) != 1 { t.Errorf("esperaba 1 permiso, got %v", len(resp)) }
	if resp[0].Codigo != "test:1" { t.Errorf("Codigo incorrecto: %v", resp[0].Codigo) }
}

func TestListarMisPermisosCaeASistema(t *testing.T) {
	uc := listarmispermisos.NewListarMisPermisosCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", Nombre: "editor"}},
		&mockRolPermisoRepo{
			permisos: nil,
			err:      nil,
		},
	)
	resp, err := uc.Ejecutar(context.Background(), "editor", "t-1")
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	_ = resp
}

func TestListarMisPermisosRolNoEncontrado(t *testing.T) {
	uc := listarmispermisos.NewListarMisPermisosCasoDeUso(
		&mockRolRepo{err: errors.New("no encontrado")},
		&mockRolPermisoRepo{},
	)
	_, err := uc.Ejecutar(context.Background(), "rol-x", "t-1")
	if err == nil { t.Fatal("esperaba error por rol no encontrado") }
}

func TestListarMisPermisosErrorRepo(t *testing.T) {
	uc := listarmispermisos.NewListarMisPermisosCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1"}},
		&mockRolPermisoRepo{err: errors.New("fallo bd")},
	)
	_, err := uc.Ejecutar(context.Background(), "editor", "t-1")
	if err == nil { t.Fatal("esperaba error de repositorio") }
}

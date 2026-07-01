package listroles_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listroles"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockRolRepo struct {
	roles []*rbac.RolDB
	err   error
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) {
	return m.roles, m.err
}
func (m *mockRolRepo) Crear(ctx context.Context, r *rbac.RolDB) error { return nil }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, id, desc string) error { return nil }
func (m *mockRolRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockPermisoRepo struct {
	permisos []*rbac.PermisoDB
	err      error
}

func (m *mockPermisoRepo) ObtenerPorCodigo(ctx context.Context, codigo string) (*rbac.PermisoDB, error) { return nil, nil }
func (m *mockPermisoRepo) Listar(ctx context.Context) ([]*rbac.PermisoDB, error) { return nil, nil }
func (m *mockPermisoRepo) Crear(ctx context.Context, p *rbac.PermisoDB) error { return nil }
func (m *mockPermisoRepo) ActualizarNombreDescripcion(ctx context.Context, id, nombre, desc string) error { return nil }
func (m *mockPermisoRepo) ListarPorRol(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) {
	return m.permisos, m.err
}

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

func TestListarRolesExitoso(t *testing.T) {
	uc := listroles.NewListarRolesCasoDeUso(
		&mockRolRepo{roles: []*rbac.RolDB{
			{ID: "r-1", Nombre: "admin", EsSistema: true},
			{ID: "r-2", Nombre: "editor", EsSistema: false},
		}},
		&mockPermisoRepo{permisos: []*rbac.PermisoDB{
			{Codigo: "test:1"},
		}},
		&mockAuthSvc{permiso: true},
	)
	resp, err := uc.Ejecutar(context.Background(), &listroles.ComandoListarRoles{
		Paginacion: shareddomain.Paginacion{Pagina: 1, TamanoPagina: 10},
		TenantID:   "t-1", EjecutorID: "admin-1",
	})
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.Total != 2 { t.Errorf("Total incorrecto: %v", resp.Total) }
	if resp.Pagina != 1 { t.Errorf("Pagina incorrecto: %v", resp.Pagina) }
}

func TestListarRolesSinPermiso(t *testing.T) {
	uc := listroles.NewListarRolesCasoDeUso(&mockRolRepo{}, &mockPermisoRepo{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), &listroles.ComandoListarRoles{TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestListarRolesErrorRepo(t *testing.T) {
	uc := listroles.NewListarRolesCasoDeUso(
		&mockRolRepo{err: errors.New("fallo bd")},
		&mockPermisoRepo{}, &mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &listroles.ComandoListarRoles{TenantID: "t-1", EjecutorID: "e-1"})
	if err == nil { t.Fatal("esperaba error al listar roles") }
}

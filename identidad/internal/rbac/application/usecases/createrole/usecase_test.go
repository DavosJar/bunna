package createrole_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/createrole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockRolRepo struct {
	errCrear error
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) ObtenerPorNombreYTenant(ctx context.Context, nombre string, tenantID string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Crear(ctx context.Context, r *rbac.RolDB) error { return m.errCrear }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, id, desc string) error { return nil }
func (m *mockRolRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockPermisoRepo struct{}

func (m *mockPermisoRepo) ObtenerPorCodigo(ctx context.Context, codigo string) (*rbac.PermisoDB, error) { return nil, nil }
func (m *mockPermisoRepo) Listar(ctx context.Context) ([]*rbac.PermisoDB, error) { return nil, nil }
func (m *mockPermisoRepo) Crear(ctx context.Context, p *rbac.PermisoDB) error { return nil }
func (m *mockPermisoRepo) ActualizarNombreDescripcion(ctx context.Context, id, nombre, desc string) error { return nil }
func (m *mockPermisoRepo) ListarPorRol(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) { return nil, nil }

type mockRolPermisoRepo struct{}

func (m *mockRolPermisoRepo) AsignarPermiso(ctx context.Context, rolID, permisoID, tenantID, asignadoPor string) error { return nil }
func (m *mockRolPermisoRepo) EliminarPermiso(ctx context.Context, rolID, permisoID, tenantID string) error { return nil }
func (m *mockRolPermisoRepo) ListarPorRolYTenant(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) { return nil, nil }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

func TestCrearRolExitoso(t *testing.T) {
	uc := createrole.NewCrearRolCasoDeUso(&mockRolRepo{}, &mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: true})
	resp, err := uc.Ejecutar(context.Background(), &createrole.ComandoCrearRol{
		Nombre: "editor", Descripcion: "Puede editar", TenantID: "t-1", EjecutorID: "admin-1",
	})
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.Nombre != "editor" { t.Errorf("Nombre incorrecto: %v", resp.Nombre) }
	if resp.EsSistema { t.Error("rol creado no debe ser de sistema") }
	if resp.CreadoEn == "" { t.Error("CreadoEn vacío") }
}

func TestCrearRolSinPermiso(t *testing.T) {
	uc := createrole.NewCrearRolCasoDeUso(&mockRolRepo{}, &mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), &createrole.ComandoCrearRol{TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestCrearRolFalloAlCrear(t *testing.T) {
	uc := createrole.NewCrearRolCasoDeUso(
		&mockRolRepo{errCrear: errors.New("fallo bd")},
		&mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &createrole.ComandoCrearRol{
		Nombre: "editor", TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error al crear rol") }
}

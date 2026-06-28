package assignpermissiontorole_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/assignpermissiontorole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockRolRepo struct {
	rol        *rbac.RolDB
	errObtener error
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) { return m.rol, m.errObtener }
func (m *mockRolRepo) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Crear(ctx context.Context, r *rbac.RolDB) error { return nil }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, id, desc string) error { return nil }

type mockPermisoRepo struct {
	permiso      *rbac.PermisoDB
	errObtener   error
}

func (m *mockPermisoRepo) ObtenerPorCodigo(ctx context.Context, codigo string) (*rbac.PermisoDB, error) {
	return m.permiso, m.errObtener
}
func (m *mockPermisoRepo) Listar(ctx context.Context) ([]*rbac.PermisoDB, error) { return nil, nil }
func (m *mockPermisoRepo) Crear(ctx context.Context, p *rbac.PermisoDB) error { return nil }
func (m *mockPermisoRepo) ActualizarNombreDescripcion(ctx context.Context, id, nombre, desc string) error { return nil }
func (m *mockPermisoRepo) ListarPorRol(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) { return nil, nil }

type mockRolPermisoRepo struct {
	errAsignar error
}

func (m *mockRolPermisoRepo) AsignarPermiso(ctx context.Context, rolID, permisoID, tenantID, asignadoPor string) error {
	return m.errAsignar
}
func (m *mockRolPermisoRepo) EliminarPermiso(ctx context.Context, rolID, permisoID, tenantID string) error { return nil }
func (m *mockRolPermisoRepo) ListarPorRolYTenant(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) { return nil, nil }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

type mockRolPublisher struct{}

func (m *mockRolPublisher) PublicarRolActualizado(ctx context.Context, rolID, tenantID string, permisos []string) error { return nil }

func TestAsignarPermisoARolExitoso(t *testing.T) {
	authSvc := &mockAuthSvc{permiso: true}
	uc := assignpermissiontorole.NewAsignarPermisoARolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: false}},
		&mockPermisoRepo{permiso: &rbac.PermisoDB{ID: "p-1"}},
		&mockRolPermisoRepo{},
		authSvc,
		&mockRolPublisher{},
	)
	resp, err := uc.Ejecutar(context.Background(), &assignpermissiontorole.ComandoAsignarPermisoARol{
		RolID: "r-1", PermisoCodigo: rbac.PermisoRolCrear,
		TenantID: "t-1", EjecutorID: "admin-1", AsignadoPor: "admin-1",
	})
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.RolID != "r-1" { t.Errorf("RolID incorrecto: %v", resp.RolID) }
	if resp.AsignadoEn == "" { t.Error("AsignadoEn vacío") }
}

func TestAsignarPermisoARolSinPermiso(t *testing.T) {
	uc := assignpermissiontorole.NewAsignarPermisoARolCasoDeUso(&mockRolRepo{}, &mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: false}, &mockRolPublisher{})
	_, err := uc.Ejecutar(context.Background(), &assignpermissiontorole.ComandoAsignarPermisoARol{TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestAsignarPermisoARolRolInmutable(t *testing.T) {
	uc := assignpermissiontorole.NewAsignarPermisoARolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: true}},
		&mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: true}, &mockRolPublisher{},
	)
	_, err := uc.Ejecutar(context.Background(), &assignpermissiontorole.ComandoAsignarPermisoARol{
		RolID: "r-1", PermisoCodigo: rbac.PermisoRolCrear, TenantID: "t-1", EjecutorID: "e-1",
	})
	if !errors.Is(err, rbac.ErrRolInmutable) {
		t.Errorf("esperaba ErrRolInmutable, got %v", err)
	}
}

func TestAsignarPermisoARolRolNoEncontrado(t *testing.T) {
	uc := assignpermissiontorole.NewAsignarPermisoARolCasoDeUso(
		&mockRolRepo{errObtener: errors.New("no encontrado")},
		&mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: true}, &mockRolPublisher{},
	)
	_, err := uc.Ejecutar(context.Background(), &assignpermissiontorole.ComandoAsignarPermisoARol{
		RolID: "r-x", PermisoCodigo: rbac.PermisoRolCrear, TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error por rol no encontrado") }
}

func TestAsignarPermisoARolPermisoNoEncontrado(t *testing.T) {
	authSvc := &mockAuthSvc{permiso: true}
	uc := assignpermissiontorole.NewAsignarPermisoARolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: false}},
		&mockPermisoRepo{errObtener: errors.New("no encontrado")},
		&mockRolPermisoRepo{}, authSvc, &mockRolPublisher{},
	)
	_, err := uc.Ejecutar(context.Background(), &assignpermissiontorole.ComandoAsignarPermisoARol{
		RolID: "r-1", PermisoCodigo: "bad-code", TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error por permiso no encontrado") }
}

func TestAsignarPermisoARolFalloAsignacion(t *testing.T) {
	authSvc := &mockAuthSvc{permiso: true}
	uc := assignpermissiontorole.NewAsignarPermisoARolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: false}},
		&mockPermisoRepo{permiso: &rbac.PermisoDB{ID: "p-1"}},
		&mockRolPermisoRepo{errAsignar: errors.New("fallo bd")}, authSvc, &mockRolPublisher{},
	)
	_, err := uc.Ejecutar(context.Background(), &assignpermissiontorole.ComandoAsignarPermisoARol{
		RolID: "r-1", PermisoCodigo: rbac.PermisoRolCrear, TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error al asignar permiso") }
}

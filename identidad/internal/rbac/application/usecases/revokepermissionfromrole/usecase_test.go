package revokepermissionfromrole_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/revokepermissionfromrole"
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
	permiso    *rbac.PermisoDB
	errObtener error
}

func (m *mockPermisoRepo) ObtenerPorCodigo(ctx context.Context, codigo string) (*rbac.PermisoDB, error) {
	return m.permiso, m.errObtener
}
func (m *mockPermisoRepo) Listar(ctx context.Context) ([]*rbac.PermisoDB, error) { return nil, nil }
func (m *mockPermisoRepo) Crear(ctx context.Context, p *rbac.PermisoDB) error { return nil }
func (m *mockPermisoRepo) ActualizarNombreDescripcion(ctx context.Context, id, nombre, desc string) error { return nil }
func (m *mockPermisoRepo) ListarPorRol(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) { return nil, nil }

type mockRolPermisoRepo struct {
	errEliminar error
}

func (m *mockRolPermisoRepo) AsignarPermiso(ctx context.Context, rolID, permisoID, tenantID, asignadoPor string) error { return nil }
func (m *mockRolPermisoRepo) EliminarPermiso(ctx context.Context, rolID, permisoID, tenantID string) error {
	return m.errEliminar
}
func (m *mockRolPermisoRepo) ListarPorRolYTenant(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) { return nil, nil }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

func TestRevocarPermisoDeRolExitoso(t *testing.T) {
	authSvc := &mockAuthSvc{permiso: true}
	uc := revokepermissionfromrole.NewRevocarPermisoDeRolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: false}},
		&mockPermisoRepo{permiso: &rbac.PermisoDB{ID: "p-1"}},
		&mockRolPermisoRepo{},
		authSvc,
	)
	resp, err := uc.Ejecutar(context.Background(), &revokepermissionfromrole.ComandoRevocarPermisoDeRol{
		RolID: "r-1", PermisoCodigo: rbac.PermisoRolCrear,
		TenantID: "t-1", EjecutorID: "admin-1",
	})
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.RolID != "r-1" { t.Errorf("RolID incorrecto: %v", resp.RolID) }
	if resp.RevocadoEn == "" { t.Error("RevocadoEn vacío") }
}

func TestRevocarPermisoDeRolSinPermiso(t *testing.T) {
	uc := revokepermissionfromrole.NewRevocarPermisoDeRolCasoDeUso(&mockRolRepo{}, &mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), &revokepermissionfromrole.ComandoRevocarPermisoDeRol{TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestRevocarPermisoDeRolRolInmutable(t *testing.T) {
	uc := revokepermissionfromrole.NewRevocarPermisoDeRolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: true}},
		&mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &revokepermissionfromrole.ComandoRevocarPermisoDeRol{
		RolID: "r-1", PermisoCodigo: rbac.PermisoRolCrear, TenantID: "t-1", EjecutorID: "e-1",
	})
	if !errors.Is(err, rbac.ErrRolInmutable) {
		t.Errorf("esperaba ErrRolInmutable, got %v", err)
	}
}

func TestRevocarPermisoDeRolRolNoEncontrado(t *testing.T) {
	uc := revokepermissionfromrole.NewRevocarPermisoDeRolCasoDeUso(
		&mockRolRepo{errObtener: errors.New("no encontrado")},
		&mockPermisoRepo{}, &mockRolPermisoRepo{}, &mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &revokepermissionfromrole.ComandoRevocarPermisoDeRol{
		RolID: "r-x", PermisoCodigo: rbac.PermisoRolCrear, TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error por rol no encontrado") }
}

func TestRevocarPermisoDeRolPermisoNoEncontrado(t *testing.T) {
	authSvc := &mockAuthSvc{permiso: true}
	uc := revokepermissionfromrole.NewRevocarPermisoDeRolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: false}},
		&mockPermisoRepo{errObtener: errors.New("no encontrado")},
		&mockRolPermisoRepo{}, authSvc,
	)
	_, err := uc.Ejecutar(context.Background(), &revokepermissionfromrole.ComandoRevocarPermisoDeRol{
		RolID: "r-1", PermisoCodigo: "bad", TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error por permiso no encontrado") }
}

func TestRevocarPermisoDeRolFalloRevocacion(t *testing.T) {
	authSvc := &mockAuthSvc{permiso: true}
	uc := revokepermissionfromrole.NewRevocarPermisoDeRolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: false}},
		&mockPermisoRepo{permiso: &rbac.PermisoDB{ID: "p-1"}},
		&mockRolPermisoRepo{errEliminar: errors.New("fallo bd")}, authSvc,
	)
	_, err := uc.Ejecutar(context.Background(), &revokepermissionfromrole.ComandoRevocarPermisoDeRol{
		RolID: "r-1", PermisoCodigo: rbac.PermisoRolCrear, TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error al revocar permiso") }
}

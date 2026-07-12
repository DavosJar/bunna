package assignrole_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/assignrole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockUsuarioRolRepo struct {
	errCrear error
}

func (m *mockUsuarioRolRepo) Crear(ctx context.Context, uid, rid string) error    { return m.errCrear }
func (m *mockUsuarioRolRepo) Eliminar(ctx context.Context, uid, rid string) error { return nil }
func (m *mockUsuarioRolRepo) ListarRolesPorUsuario(ctx context.Context, uid string) ([]*rbac.RolDB, error) {
	return nil, nil
}
func (m *mockUsuarioRolRepo) TieneRol(ctx context.Context, uid, rol string) (bool, error) {
	return false, nil
}
func (m *mockUsuarioRolRepo) ObtenerUsuarioConRol(ctx context.Context, rolNombre string) (string, bool, error) {
	return "", false, nil
}

type mockUTRR struct {
	errCrear error
}

func (m *mockUTRR) Crear(ctx context.Context, uid, tid, rid string) error    { return m.errCrear }
func (m *mockUTRR) Eliminar(ctx context.Context, uid, tid, rid string) error { return nil }
func (m *mockUTRR) ListarRolesPorUsuarioEnTenant(ctx context.Context, uid, tid string) ([]*rbac.RolDB, error) {
	return nil, nil
}
func (m *mockUTRR) TieneRolEnTenant(ctx context.Context, uid, tid, rol string) (bool, error) {
	return false, nil
}

type mockRolRepo struct {
	rol        *rbac.RolDB
	errObtener error
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) {
	return nil, nil
}
func (m *mockRolRepo) ObtenerPorNombreYTenant(ctx context.Context, nombre string, tenantID string) (*rbac.RolDB, error) {
	return nil, nil
}
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) {
	return m.rol, m.errObtener
}
func (m *mockRolRepo) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) {
	return nil, nil
}
func (m *mockRolRepo) Crear(ctx context.Context, r *rbac.RolDB) error                   { return nil }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, id, desc string) error { return nil }
func (m *mockRolRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

func TestAsignarRolExitosoEnTenant(t *testing.T) {
	uc := assignrole.NewAsignarRolCasoDeUso(
		&mockUsuarioRolRepo{},
		&mockUTRR{},
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", Nombre: "editor"}},
		&mockAuthSvc{permiso: true},
	)
	resp, err := uc.Ejecutar(context.Background(), &assignrole.ComandoAsignarRol{
		UsuarioID: "user-1", RolID: "r-1", TenantID: "t-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.AsignadoEn == "" {
		t.Error("AsignadoEn vacío")
	}
}

func TestAsignarRolGlobalExitoso(t *testing.T) {
	uc := assignrole.NewAsignarRolCasoDeUso(
		&mockUsuarioRolRepo{},
		&mockUTRR{},
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", Nombre: "sys_admin"}},
		&mockAuthSvc{permiso: true},
	)
	resp, err := uc.Ejecutar(context.Background(), &assignrole.ComandoAsignarRol{
		UsuarioID: "user-1", RolID: "r-1", TenantID: "", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.AsignadoEn == "" {
		t.Error("AsignadoEn vacío")
	}
}

func TestAsignarRolSinPermiso(t *testing.T) {
	uc := assignrole.NewAsignarRolCasoDeUso(&mockUsuarioRolRepo{}, &mockUTRR{}, &mockRolRepo{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), &assignrole.ComandoAsignarRol{UsuarioID: "u-1", RolID: "r-1", TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestAsignarRolSysAdminConTenant(t *testing.T) {
	uc := assignrole.NewAsignarRolCasoDeUso(
		&mockUsuarioRolRepo{}, &mockUTRR{},
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", Nombre: rbac.RolSysAdmin}},
		&mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &assignrole.ComandoAsignarRol{
		UsuarioID: "u-1", RolID: "r-1", TenantID: "t-1", EjecutorID: "e-1",
	})
	if !errors.Is(err, rbac.ErrSysAdminRequiereTenantVacio) {
		t.Errorf("esperaba ErrSysAdminRequiereTenantVacio, got %v", err)
	}
}

func TestAsignarRolNoEncontrado(t *testing.T) {
	uc := assignrole.NewAsignarRolCasoDeUso(
		&mockUsuarioRolRepo{}, &mockUTRR{},
		&mockRolRepo{errObtener: errors.New("no encontrado")},
		&mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &assignrole.ComandoAsignarRol{
		UsuarioID: "u-1", RolID: "r-x", TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil {
		t.Fatal("esperaba error por rol no encontrado")
	}
}

func TestAsignarRolFalloAlCrear(t *testing.T) {
	uc := assignrole.NewAsignarRolCasoDeUso(
		&mockUsuarioRolRepo{errCrear: errors.New("fallo bd")}, &mockUTRR{},
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1"}},
		&mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &assignrole.ComandoAsignarRol{
		UsuarioID: "u-1", RolID: "r-1", TenantID: "", EjecutorID: "e-1",
	})
	if err == nil {
		t.Fatal("esperaba error al crear rol global")
	}
}

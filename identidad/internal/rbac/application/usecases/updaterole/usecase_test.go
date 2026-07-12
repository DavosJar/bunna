package updaterole_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/updaterole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockRolRepo struct {
	rol              *rbac.RolDB
	errObtener       error
	errActualizarDesc error
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) ObtenerPorNombreYTenant(ctx context.Context, nombre string, tenantID string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) { return m.rol, m.errObtener }
func (m *mockRolRepo) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Crear(ctx context.Context, r *rbac.RolDB) error { return nil }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, id, desc string) error { return m.errActualizarDesc }
func (m *mockRolRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

func TestModificarRolExitoso(t *testing.T) {
	rolRepo := &mockRolRepo{rol: &rbac.RolDB{ID: "r-1", Nombre: "editor", Descripcion: "old", EsSistema: false}}
	uc := updaterole.NewModificarRolCasoDeUso(rolRepo, &mockAuthSvc{permiso: true})
	resp, err := uc.Ejecutar(context.Background(), &updaterole.ComandoModificarRol{
		RolID: "r-1", Nombre: "editor", Descripcion: "new desc",
		TenantID: "t-1", EjecutorID: "admin-1",
	})
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.ID != "r-1" { t.Errorf("ID incorrecto: %v", resp.ID) }
	if resp.ModificadoEn == "" { t.Error("ModificadoEn vacío") }
}

func TestModificarRolSinPermiso(t *testing.T) {
	uc := updaterole.NewModificarRolCasoDeUso(&mockRolRepo{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), &updaterole.ComandoModificarRol{RolID: "r-1", TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestModificarRolSistemaInmutable(t *testing.T) {
	rolRepo := &mockRolRepo{rol: &rbac.RolDB{ID: "r-1", Nombre: "sys_admin", EsSistema: true}}
	uc := updaterole.NewModificarRolCasoDeUso(rolRepo, &mockAuthSvc{permiso: true})
	_, err := uc.Ejecutar(context.Background(), &updaterole.ComandoModificarRol{RolID: "r-1", TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrRolInmutable) {
		t.Errorf("esperaba ErrRolInmutable, got %v", err)
	}
}

func TestModificarRolNoEncontrado(t *testing.T) {
	rolRepo := &mockRolRepo{errObtener: errors.New("no encontrado")}
	uc := updaterole.NewModificarRolCasoDeUso(rolRepo, &mockAuthSvc{permiso: true})
	_, err := uc.Ejecutar(context.Background(), &updaterole.ComandoModificarRol{RolID: "r-x", TenantID: "t-1", EjecutorID: "e-1"})
	if err == nil { t.Fatal("esperaba error por rol no encontrado") }
}

func TestModificarRolFalloAlActualizar(t *testing.T) {
	rolRepo := &mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: false}, errActualizarDesc: errors.New("fallo bd")}
	uc := updaterole.NewModificarRolCasoDeUso(rolRepo, &mockAuthSvc{permiso: true})
	_, err := uc.Ejecutar(context.Background(), &updaterole.ComandoModificarRol{
		RolID: "r-1", Descripcion: "new", TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error al actualizar") }
}

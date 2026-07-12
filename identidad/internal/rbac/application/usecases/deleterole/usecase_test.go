package deleterole_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/deleterole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockRolRepo struct {
	rol        *rbac.RolDB
	errObtener error
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) ObtenerPorNombreYTenant(ctx context.Context, nombre string, tenantID string) (*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) { return m.rol, m.errObtener }
func (m *mockRolRepo) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Crear(ctx context.Context, r *rbac.RolDB) error { return nil }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, id, desc string) error { return nil }
func (m *mockRolRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

func TestEliminarRolExitoso(t *testing.T) {
	uc := deleterole.NewEliminarRolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: false}},
		&mockAuthSvc{permiso: true},
	)
	resp, err := uc.Ejecutar(context.Background(), &deleterole.ComandoEliminarRol{
		RolID: "r-1", TenantID: "t-1", EjecutorID: "admin-1",
	})
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.RolID != "r-1" { t.Errorf("RolID incorrecto: %v", resp.RolID) }
	if resp.EliminadoEn == "" { t.Error("EliminadoEn vacío") }
}

func TestEliminarRolSinPermiso(t *testing.T) {
	uc := deleterole.NewEliminarRolCasoDeUso(&mockRolRepo{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), &deleterole.ComandoEliminarRol{RolID: "r-1", TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestEliminarRolSistemaInmutable(t *testing.T) {
	uc := deleterole.NewEliminarRolCasoDeUso(
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1", EsSistema: true}},
		&mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &deleterole.ComandoEliminarRol{RolID: "r-1", TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrRolInmutable) {
		t.Errorf("esperaba ErrRolInmutable, got %v", err)
	}
}

func TestEliminarRolNoEncontrado(t *testing.T) {
	uc := deleterole.NewEliminarRolCasoDeUso(
		&mockRolRepo{errObtener: errors.New("no encontrado")},
		&mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &deleterole.ComandoEliminarRol{RolID: "r-x", TenantID: "t-1", EjecutorID: "e-1"})
	if err == nil { t.Fatal("esperaba error por rol no encontrado") }
}

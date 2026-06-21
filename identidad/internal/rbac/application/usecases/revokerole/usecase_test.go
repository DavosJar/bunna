package revokerole_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/revokerole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type mockUsuarioRolRepo struct {
	errEliminar error
}

func (m *mockUsuarioRolRepo) Crear(ctx context.Context, uid, rid string) error { return nil }
func (m *mockUsuarioRolRepo) Eliminar(ctx context.Context, uid, rid string) error { return m.errEliminar }
func (m *mockUsuarioRolRepo) ListarRolesPorUsuario(ctx context.Context, uid string) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockUsuarioRolRepo) TieneRol(ctx context.Context, uid, rol string) (bool, error) { return false, nil }

type mockUTRR struct {
	errEliminar error
}

func (m *mockUTRR) Crear(ctx context.Context, uid, tid, rid string) error { return nil }
func (m *mockUTRR) Eliminar(ctx context.Context, uid, tid, rid string) error { return m.errEliminar }
func (m *mockUTRR) ListarRolesPorUsuarioEnTenant(ctx context.Context, uid, tid string) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockUTRR) TieneRolEnTenant(ctx context.Context, uid, tid, rol string) (bool, error) { return false, nil }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, uid, tid, permiso string) (bool, error) {
	return m.permiso, m.err
}

func TestRevocarRolExitosoEnTenant(t *testing.T) {
	uc := revokerole.NewRevocarRolCasoDeUso(
		&mockUsuarioRolRepo{}, &mockUTRR{}, &mockAuthSvc{permiso: true},
	)
	resp, err := uc.Ejecutar(context.Background(), &revokerole.ComandoRevocarRol{
		UsuarioID: "user-1", RolID: "r-1", TenantID: "t-1", EjecutorID: "admin-1",
	})
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.RevocadoEn == "" { t.Error("RevocadoEn vacío") }
}

func TestRevocarRolGlobalExitoso(t *testing.T) {
	uc := revokerole.NewRevocarRolCasoDeUso(
		&mockUsuarioRolRepo{}, &mockUTRR{}, &mockAuthSvc{permiso: true},
	)
	resp, err := uc.Ejecutar(context.Background(), &revokerole.ComandoRevocarRol{
		UsuarioID: "user-1", RolID: "r-1", TenantID: "", EjecutorID: "admin-1",
	})
	if err != nil { t.Fatalf("error inesperado: %v", err) }
	if resp.RevocadoEn == "" { t.Error("RevocadoEn vacío") }
}

func TestRevocarRolSinPermiso(t *testing.T) {
	uc := revokerole.NewRevocarRolCasoDeUso(&mockUsuarioRolRepo{}, &mockUTRR{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), &revokerole.ComandoRevocarRol{UsuarioID: "u-1", RolID: "r-1", TenantID: "t-1", EjecutorID: "e-1"})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestRevocarRolFalloGlobal(t *testing.T) {
	uc := revokerole.NewRevocarRolCasoDeUso(
		&mockUsuarioRolRepo{errEliminar: errors.New("fallo bd")},
		&mockUTRR{}, &mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &revokerole.ComandoRevocarRol{
		UsuarioID: "u-1", RolID: "r-1", TenantID: "", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error al revocar rol global") }
}

func TestRevocarRolFalloEnTenant(t *testing.T) {
	uc := revokerole.NewRevocarRolCasoDeUso(
		&mockUsuarioRolRepo{},
		&mockUTRR{errEliminar: errors.New("fallo bd")},
		&mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &revokerole.ComandoRevocarRol{
		UsuarioID: "u-1", RolID: "r-1", TenantID: "t-1", EjecutorID: "e-1",
	})
	if err == nil { t.Fatal("esperaba error al revocar rol en tenant") }
}

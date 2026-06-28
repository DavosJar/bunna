package checkpermission_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/checkpermission"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type mockUsuarioRolRepo struct {
	roles []*rbac.RolDB
	err   error
}

func (m *mockUsuarioRolRepo) Crear(ctx context.Context, uid, rid string) error    { return nil }
func (m *mockUsuarioRolRepo) Eliminar(ctx context.Context, uid, rid string) error { return nil }
func (m *mockUsuarioRolRepo) ListarRolesPorUsuario(ctx context.Context, uid string) ([]*rbac.RolDB, error) {
	return m.roles, m.err
}
func (m *mockUsuarioRolRepo) TieneRol(ctx context.Context, uid, rol string) (bool, error) {
	return false, nil
}
func (m *mockUsuarioRolRepo) ObtenerUsuarioConRol(ctx context.Context, rolNombre string) (string, bool, error) {
	return "", false, nil
}

type mockUTRR struct {
	roles []*rbac.RolDB
	err   error
}

func (m *mockUTRR) Crear(ctx context.Context, uid, tid, rid string) error    { return nil }
func (m *mockUTRR) Eliminar(ctx context.Context, uid, tid, rid string) error { return nil }
func (m *mockUTRR) ListarRolesPorUsuarioEnTenant(ctx context.Context, uid, tid string) ([]*rbac.RolDB, error) {
	return m.roles, m.err
}
func (m *mockUTRR) TieneRolEnTenant(ctx context.Context, uid, tid, rol string) (bool, error) {
	return false, nil
}

type mockPermisoRepo struct {
	permisos []*rbac.PermisoDB
	err      error
}

func (m *mockPermisoRepo) ObtenerPorCodigo(ctx context.Context, codigo string) (*rbac.PermisoDB, error) {
	return nil, m.err
}
func (m *mockPermisoRepo) Listar(ctx context.Context) ([]*rbac.PermisoDB, error) {
	return m.permisos, m.err
}
func (m *mockPermisoRepo) Crear(ctx context.Context, p *rbac.PermisoDB) error { return nil }
func (m *mockPermisoRepo) ActualizarNombreDescripcion(ctx context.Context, id, nombre, desc string) error {
	return nil
}
func (m *mockPermisoRepo) ListarPorRol(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) {
	return m.permisos, m.err
}

func TestTienePermisoEnTenant(t *testing.T) {
	uc := checkpermission.NewVerificarPermisoCasoDeUso(
		&mockUsuarioRolRepo{},
		&mockUTRR{roles: []*rbac.RolDB{{Nombre: rbac.RolAdministrador, EsSistema: true}}},
		&mockPermisoRepo{},
	)
	ok, err := uc.TienePermiso(context.Background(), "user-1", "tenant-1", rbac.PermisoRolCrear)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !ok {
		t.Error("esperaba true (sys_admin tiene todos los permisos)")
	}
}

func TestTienePermisoSinTenantRolGlobal(t *testing.T) {
	uc := checkpermission.NewVerificarPermisoCasoDeUso(
		&mockUsuarioRolRepo{roles: []*rbac.RolDB{{Nombre: "sys_admin", EsSistema: true}}},
		&mockUTRR{},
		&mockPermisoRepo{},
	)
	ok, err := uc.TienePermiso(context.Background(), "user-1", "", rbac.PermisoRolCrear)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !ok {
		t.Error("esperaba true (sys_admin global)")
	}
}

func TestTienePermisoDenegado(t *testing.T) {
	uc := checkpermission.NewVerificarPermisoCasoDeUso(
		&mockUsuarioRolRepo{},
		&mockUTRR{},
		&mockPermisoRepo{},
	)
	ok, err := uc.TienePermiso(context.Background(), "user-1", "tenant-1", rbac.PermisoRolCrear)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if ok {
		t.Error("esperaba false (sin roles)")
	}
}

func TestTienePermisoErrorTenantRoles(t *testing.T) {
	uc := checkpermission.NewVerificarPermisoCasoDeUso(
		&mockUsuarioRolRepo{},
		&mockUTRR{err: errors.New("fallo bd")},
		&mockPermisoRepo{},
	)
	_, err := uc.TienePermiso(context.Background(), "user-1", "tenant-1", rbac.PermisoRolCrear)
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

func TestTienePermisoErrorGlobalRoles(t *testing.T) {
	uc := checkpermission.NewVerificarPermisoCasoDeUso(
		&mockUsuarioRolRepo{err: errors.New("fallo bd")},
		&mockUTRR{},
		&mockPermisoRepo{},
	)
	_, err := uc.TienePermiso(context.Background(), "user-1", "", rbac.PermisoRolCrear)
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

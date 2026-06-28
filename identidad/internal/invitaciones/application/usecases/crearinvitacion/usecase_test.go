package crearInvitacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/crearinvitacion"
	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type mockInvRepo struct {
	errCrear error
}

func (m *mockInvRepo) Crear(ctx context.Context, inv *invitaciones.Invitacion) error { return m.errCrear }
func (m *mockInvRepo) ObtenerPorTokenHash(ctx context.Context, hash string) (*invitaciones.Invitacion, error) { return nil, nil }
func (m *mockInvRepo) MarcarAceptada(ctx context.Context, id string) error { return nil }
func (m *mockInvRepo) ObtenerPorID(ctx context.Context, id string) (*invitaciones.Invitacion, error) { return nil, nil }
func (m *mockInvRepo) ListarPorTenant(ctx context.Context, tenantID string, pag shareddomain.Paginacion, estado string) ([]*invitaciones.Invitacion, int, error) { return nil, 0, nil }
func (m *mockInvRepo) ActualizarTokenHash(ctx context.Context, id string, tokenHash string) error { return nil }

type mockTenantRepo struct {
	tenant *tenant.Tenant
	err    error
}

func (m *mockTenantRepo) Crear(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) { return t, nil }
func (m *mockTenantRepo) ObtenerPorID(ctx context.Context, id string) (*tenant.Tenant, error) { return m.tenant, m.err }
func (m *mockTenantRepo) ObtenerPorSlug(ctx context.Context, slug string) (*tenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) Actualizar(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) { return t, nil }
func (m *mockTenantRepo) Listar(ctx context.Context) ([]*tenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) ListarPorUsuario(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) { return nil, nil }

type mockRolRepo struct {
	rol *rbac.RolDB
	err error
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) { return m.rol, m.err }
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) { return m.rol, m.err }
func (m *mockRolRepo) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockRolRepo) Crear(ctx context.Context, r *rbac.RolDB) error { return nil }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, id, desc string) error { return nil }

type mockGenID struct{ id string }

func (m *mockGenID) NextID(ctx context.Context) (string, error) {
	return m.id, nil
}

type mockAuthSvc struct {
	ok  bool
	err error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func TestCrearInvitacionExitoso(t *testing.T) {
	now := time.Now()
	tenantObj := tenant.NuevoTenantDesdeBD("t-1", "Mi Tenant", "mi-tenant", true, now, now)
	invRepo := &mockInvRepo{}
	tenantRepo := &mockTenantRepo{tenant: tenantObj}
	rolRepo := &mockRolRepo{rol: &rbac.RolDB{ID: "r-1", Nombre: "admin"}}
	emailSvc := &notificaciones.MockEmailServicio{}
	genID := &mockGenID{id: "inv-id-1"}
	uc := crearInvitacion.NewCrearInvitacionCasoDeUso(&mockAuthSvc{ok: true}, invRepo, tenantRepo, rolRepo, emailSvc, genID, "http://frontend", 48*time.Hour)

	resp, err := uc.Ejecutar(context.Background(), &crearInvitacion.ComandoCrearInvitacion{
		TenantID: "t-1", RolID: "r-1", Correo: "invitado@test.com",
		Nombre: "Invitado", CreadoPor: "user-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.ID != "inv-id-1" {
		t.Errorf("ID incorrecto: %v", resp.ID)
	}
	if resp.Token == "" {
		t.Error("Token vacío")
	}
	ultimo := emailSvc.UltimoEmail()
	if ultimo == nil {
		t.Fatal("no se envió email")
	}
	if ultimo.Destinatario != "invitado@test.com" {
		t.Errorf("destinatario incorrecto: %v", ultimo.Destinatario)
	}
	if ultimo.Tipo != notificaciones.TipoInvitacion {
		t.Errorf("tipo incorrecto: %v", ultimo.Tipo)
	}
}

func TestCrearInvitacionCorreoVacio(t *testing.T) {
	uc := crearInvitacion.NewCrearInvitacionCasoDeUso(&mockAuthSvc{ok: true}, &mockInvRepo{}, &mockTenantRepo{}, &mockRolRepo{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, "", 0)
	_, err := uc.Ejecutar(context.Background(), &crearInvitacion.ComandoCrearInvitacion{
		TenantID: "t-1", RolID: "r-1", Correo: "",
	})
	if !errors.Is(err, invitaciones.ErrEmailRequerido) {
		t.Errorf("esperaba ErrEmailRequerido, got %v", err)
	}
}

func TestCrearInvitacionCorreoInvalido(t *testing.T) {
	uc := crearInvitacion.NewCrearInvitacionCasoDeUso(&mockAuthSvc{ok: true}, &mockInvRepo{}, &mockTenantRepo{}, &mockRolRepo{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, "", 0)
	_, err := uc.Ejecutar(context.Background(), &crearInvitacion.ComandoCrearInvitacion{
		TenantID: "t-1", RolID: "r-1", Correo: "not-an-email",
	})
	if err == nil {
		t.Fatal("esperaba error por correo inválido")
	}
}

func TestCrearInvitacionRolVacio(t *testing.T) {
	uc := crearInvitacion.NewCrearInvitacionCasoDeUso(&mockAuthSvc{ok: true}, &mockInvRepo{}, &mockTenantRepo{}, &mockRolRepo{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, "", 0)
	_, err := uc.Ejecutar(context.Background(), &crearInvitacion.ComandoCrearInvitacion{
		TenantID: "t-1", RolID: "", Correo: "a@b.com",
	})
	if !errors.Is(err, invitaciones.ErrRolRequerido) {
		t.Errorf("esperaba ErrRolRequerido, got %v", err)
	}
}

func TestCrearInvitacionTenantNoEncontrado(t *testing.T) {
	uc := crearInvitacion.NewCrearInvitacionCasoDeUso(
		&mockAuthSvc{ok: true},
		&mockInvRepo{}, &mockTenantRepo{err: errors.New("no encontrado")},
		&mockRolRepo{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, "", 0,
	)
	_, err := uc.Ejecutar(context.Background(), &crearInvitacion.ComandoCrearInvitacion{
		TenantID: "t-x", RolID: "r-1", Correo: "a@b.com",
	})
	if err == nil {
		t.Fatal("esperaba error por tenant no encontrado")
	}
}

func TestCrearInvitacionRolNoEncontrado(t *testing.T) {
	now := time.Now()
	tenantObj := tenant.NuevoTenantDesdeBD("t-1", "T", "t", true, now, now)
	uc := crearInvitacion.NewCrearInvitacionCasoDeUso(
		&mockAuthSvc{ok: true},
		&mockInvRepo{}, &mockTenantRepo{tenant: tenantObj},
		&mockRolRepo{err: errors.New("no encontrado")},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, "", 0,
	)
	_, err := uc.Ejecutar(context.Background(), &crearInvitacion.ComandoCrearInvitacion{
		TenantID: "t-1", RolID: "r-x", Correo: "a@b.com",
	})
	if err == nil {
		t.Fatal("esperaba error por rol no encontrado")
	}
}

func TestCrearInvitacionFalloAlCrear(t *testing.T) {
	now := time.Now()
	tenantObj := tenant.NuevoTenantDesdeBD("t-1", "T", "t", true, now, now)
	uc := crearInvitacion.NewCrearInvitacionCasoDeUso(
		&mockAuthSvc{ok: true},
		&mockInvRepo{errCrear: errors.New("fallo bd")}, &mockTenantRepo{tenant: tenantObj},
		&mockRolRepo{rol: &rbac.RolDB{ID: "r-1"}},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, "", 0,
	)
	_, err := uc.Ejecutar(context.Background(), &crearInvitacion.ComandoCrearInvitacion{
		TenantID: "t-1", RolID: "r-1", Correo: "a@b.com",
	})
	if err == nil {
		t.Fatal("esperaba error al crear invitación")
	}
}

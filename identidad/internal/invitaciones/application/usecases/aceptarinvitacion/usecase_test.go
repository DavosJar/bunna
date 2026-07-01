package aceptarInvitacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/aceptarinvitacion"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepo struct {
	usuario *usuario.Usuario
	err     error
}

func (m *mockUsuarioRepo) Crear(ctx context.Context, u *usuario.Usuario) (*usuario.Usuario, error) { return u, nil }
func (m *mockUsuarioRepo) Actualizar(ctx context.Context, u *usuario.Usuario) (*usuario.Usuario, error) { return u, nil }
func (m *mockUsuarioRepo) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepo) ObtenerPorID(ctx context.Context, id string) (*usuario.Usuario, error) { return m.usuario, m.err }
func (m *mockUsuarioRepo) ObtenerPorCorreo(ctx context.Context, correo string) (*usuario.Usuario, error) { return m.usuario, m.err }
func (m *mockUsuarioRepo) Listar(ctx context.Context, spec usuario.EspecificacionUsuario, pag shareddomain.Paginacion) ([]*usuario.Usuario, error) { return nil, nil }

type mockInvRepo struct {
	invitacion   *invitaciones.Invitacion
	errObtener   error
	errMarcar    error
}

func (m *mockInvRepo) Crear(ctx context.Context, inv *invitaciones.Invitacion) error { return nil }
func (m *mockInvRepo) ObtenerPorTokenHash(ctx context.Context, hash string) (*invitaciones.Invitacion, error) {
	return m.invitacion, m.errObtener
}
func (m *mockInvRepo) MarcarAceptada(ctx context.Context, id string) error { return m.errMarcar }
func (m *mockInvRepo) ObtenerPorID(ctx context.Context, id string) (*invitaciones.Invitacion, error) { return nil, nil }
func (m *mockInvRepo) ListarPorTenant(ctx context.Context, tenantID string, pag shareddomain.Paginacion, estado string) ([]*invitaciones.Invitacion, int, error) { return nil, 0, nil }
func (m *mockInvRepo) ActualizarTokenHash(ctx context.Context, id string, tokenHash string) error { return nil }
func (m *mockInvRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockMembresiaRepo struct {
	errCrear error
}

func (m *mockMembresiaRepo) Crear(ctx context.Context, memb *tenant.Membresia) error { return m.errCrear }
func (m *mockMembresiaRepo) Eliminar(ctx context.Context, uid, tid string) error { return nil }
func (m *mockMembresiaRepo) ExisteMiembro(ctx context.Context, uid, tid string) (bool, error) { return false, nil }
func (m *mockMembresiaRepo) ListarUsuariosPorTenant(ctx context.Context, tid string) ([]string, error) { return nil, nil }
func (m *mockMembresiaRepo) ListarTenantsPorUsuario(ctx context.Context, uid string) ([]string, error) { return nil, nil }

type mockUTRR struct {
	errCrear error
}

func (m *mockUTRR) Crear(ctx context.Context, uid, tid, rid string) error { return m.errCrear }
func (m *mockUTRR) Eliminar(ctx context.Context, uid, tid, rid string) error { return nil }
func (m *mockUTRR) ListarRolesPorUsuarioEnTenant(ctx context.Context, uid, tid string) ([]*rbac.RolDB, error) { return nil, nil }
func (m *mockUTRR) TieneRolEnTenant(ctx context.Context, uid, tid, rol string) (bool, error) { return false, nil }

func usuarioValido() *usuario.Usuario {
	u, _ := usuario.NuevoUsuario("user-1", "inv@test.com", "Invitado", "", "")
	return u
}

func invitacionValida() *invitaciones.Invitacion {
	return invitaciones.NuevaInvitacionDesdeBD("inv-1", "tenant-1", "rol-1",
		"inv@test.com", "Invitado", "hash:token-raw",
		time.Now().Add(48*time.Hour), false, time.Now().Add(-1*time.Hour), nil)
}

func TestAceptarInvitacionExitoso(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida()}
	membRepo := &mockMembresiaRepo{}
	utrRepo := &mockUTRR{}
	userRepo := &mockUsuarioRepo{usuario: usuarioValido()}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, membRepo, utrRepo, userRepo)

	resp, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.TenantID != "tenant-1" {
		t.Errorf("TenantID incorrecto: %v", resp.TenantID)
	}
	if resp.RolID != "rol-1" {
		t.Errorf("RolID incorrecto: %v", resp.RolID)
	}
}

func TestAceptarInvitacionTokenVacio(t *testing.T) {
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(&mockInvRepo{}, &mockMembresiaRepo{}, &mockUTRR{}, &mockUsuarioRepo{})
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "",
	})
	if !errors.Is(err, invitaciones.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestAceptarInvitacionTokenInvalido(t *testing.T) {
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(
		&mockInvRepo{errObtener: errors.New("no encontrado")},
		&mockMembresiaRepo{}, &mockUTRR{}, &mockUsuarioRepo{},
	)
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "bad-token",
	})
	if !errors.Is(err, invitaciones.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestAceptarInvitacionYaAceptada(t *testing.T) {
	inv := invitacionValida()
	inv.Aceptar()
	invRepo := &mockInvRepo{invitacion: inv}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, &mockMembresiaRepo{}, &mockUTRR{}, &mockUsuarioRepo{})
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw",
	})
	if !errors.Is(err, invitaciones.ErrYaAceptada) {
		t.Errorf("esperaba ErrYaAceptada, got %v", err)
	}
}

func TestAceptarInvitacionExpirada(t *testing.T) {
	inv := invitaciones.NuevaInvitacionDesdeBD("inv-1", "tenant-1", "rol-1",
		"inv@test.com", "Invitado", "hash:token-raw",
		time.Now().Add(-48*time.Hour), false, time.Now().Add(-96*time.Hour), nil)
	invRepo := &mockInvRepo{invitacion: inv}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, &mockMembresiaRepo{}, &mockUTRR{}, &mockUsuarioRepo{})
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw",
	})
	if !errors.Is(err, invitaciones.ErrEnlaceExpirado) {
		t.Errorf("esperaba ErrEnlaceExpirado, got %v", err)
	}
}

func TestAceptarInvitacionFalloAlMarcar(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida(), errMarcar: errors.New("fallo bd")}
	userRepo := &mockUsuarioRepo{usuario: usuarioValido()}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, &mockMembresiaRepo{}, &mockUTRR{}, userRepo)
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw",
	})
	if err == nil {
		t.Fatal("esperaba error al marcar aceptada")
	}
}

func TestAceptarInvitacionUsuarioNoRegistrado(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida()}
	userRepo := &mockUsuarioRepo{err: errors.New("no encontrado")}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, &mockMembresiaRepo{}, &mockUTRR{}, userRepo)
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw",
	})
	if !errors.Is(err, invitaciones.ErrUsuarioNoRegistrado) {
		t.Errorf("esperaba ErrUsuarioNoRegistrado, got %v", err)
	}
}

func TestAceptarInvitacionFalloAlCrearMembresia(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida()}
	userRepo := &mockUsuarioRepo{usuario: usuarioValido()}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(
		invRepo, &mockMembresiaRepo{errCrear: errors.New("fallo bd")}, &mockUTRR{}, userRepo,
	)
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw",
	})
	if err == nil {
		t.Fatal("esperaba error al crear membresía")
	}
}

func TestAceptarInvitacionFalloAlAsignarRol(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida()}
	userRepo := &mockUsuarioRepo{usuario: usuarioValido()}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(
		invRepo, &mockMembresiaRepo{},
		&mockUTRR{errCrear: errors.New("fallo bd")}, userRepo,
	)
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw",
	})
	if err == nil {
		t.Fatal("esperaba error al asignar rol")
	}
}

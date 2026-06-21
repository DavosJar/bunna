package aceptarInvitacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/aceptarinvitacion"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

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

func invitacionValida() *invitaciones.Invitacion {
	return invitaciones.NuevaInvitacionDesdeBD("inv-1", "tenant-1", "rol-1",
		"inv@test.com", "Invitado", "hash:token-raw",
		time.Now().Add(48*time.Hour), false, time.Now().Add(-1*time.Hour), nil)
}

func TestAceptarInvitacionExitoso(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida()}
	membRepo := &mockMembresiaRepo{}
	utrRepo := &mockUTRR{}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, membRepo, utrRepo)

	resp, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw", UsuarioID: "user-1",
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
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(&mockInvRepo{}, &mockMembresiaRepo{}, &mockUTRR{})
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "", UsuarioID: "user-1",
	})
	if !errors.Is(err, invitaciones.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestAceptarInvitacionTokenInvalido(t *testing.T) {
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(
		&mockInvRepo{errObtener: errors.New("no encontrado")},
		&mockMembresiaRepo{}, &mockUTRR{},
	)
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "bad-token", UsuarioID: "user-1",
	})
	if !errors.Is(err, invitaciones.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestAceptarInvitacionYaAceptada(t *testing.T) {
	inv := invitacionValida()
	now := time.Now()
	inv.Aceptar()
	_ = now
	invRepo := &mockInvRepo{invitacion: inv}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, &mockMembresiaRepo{}, &mockUTRR{})
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw", UsuarioID: "user-1",
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
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, &mockMembresiaRepo{}, &mockUTRR{})
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw", UsuarioID: "user-1",
	})
	if !errors.Is(err, invitaciones.ErrEnlaceExpirado) {
		t.Errorf("esperaba ErrEnlaceExpirado, got %v", err)
	}
}

func TestAceptarInvitacionFalloAlMarcar(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida(), errMarcar: errors.New("fallo bd")}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(invRepo, &mockMembresiaRepo{}, &mockUTRR{})
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw", UsuarioID: "user-1",
	})
	if err == nil {
		t.Fatal("esperaba error al marcar aceptada")
	}
}

func TestAceptarInvitacionFalloAlCrearMembresia(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida()}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(
		invRepo, &mockMembresiaRepo{errCrear: errors.New("fallo bd")}, &mockUTRR{},
	)
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw", UsuarioID: "user-1",
	})
	if err == nil {
		t.Fatal("esperaba error al crear membresía")
	}
}

func TestAceptarInvitacionFalloAlAsignarRol(t *testing.T) {
	invRepo := &mockInvRepo{invitacion: invitacionValida()}
	uc := aceptarInvitacion.NewAceptarInvitacionCasoDeUso(
		invRepo, &mockMembresiaRepo{},
		&mockUTRR{errCrear: errors.New("fallo bd")},
	)
	_, err := uc.Ejecutar(context.Background(), &aceptarInvitacion.ComandoAceptarInvitacion{
		Token: "token-raw", UsuarioID: "user-1",
	})
	if err == nil {
		t.Fatal("esperaba error al asignar rol")
	}
}

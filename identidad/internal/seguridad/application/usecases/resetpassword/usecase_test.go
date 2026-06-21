package resetpassword_test

import (
	"context"
	"errors"
	"testing"
	"time"

	rbacdomain "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/resetpassword"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockCredRepo struct {
	credenciales  *seguridad.CredencialesUsuario
	errObtener    error
	errActualizar error
	actualizado   bool
}

func (m *mockCredRepo) Crear(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error) { return c, nil }
func (m *mockCredRepo) Actualizar(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error) {
	if m.errActualizar != nil { return nil, m.errActualizar }
	m.actualizado = true
	return c, nil
}
func (m *mockCredRepo) ObtenerPorUsuarioID(ctx context.Context, id string) (*seguridad.CredencialesUsuario, error) { return m.credenciales, m.errObtener }
func (m *mockCredRepo) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockCredRepo) Find(ctx context.Context, _ seguridad.EspecificacionCredenciales, _ shareddomain.Paginacion) ([]*seguridad.CredencialesUsuario, error) { return nil, nil }

type mockSesionRepo struct {
	errInvalidar error
	invalidado   bool
}

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) { return s, nil }
func (m *mockSesionRepo) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) { return s, nil }
func (m *mockSesionRepo) ObtenerPorID(ctx context.Context, id string) (*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) ListarActivasPorUsuarioID(ctx context.Context, usuarioID string, ahora time.Time) ([]*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) Listar(ctx context.Context, spec sesiones_domain.EspecificacionSesion, pag shareddomain.Paginacion) ([]*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) InvalidarTodasPorUsuarioID(ctx context.Context, usuarioID string) error {
	if m.errInvalidar != nil { return m.errInvalidar }
	m.invalidado = true
	return nil
}
func (m *mockSesionRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockEncriptacion struct {
	hash string
}

func (m *mockEncriptacion) Hashear(password string) (string, error) { return m.hash, nil }
func (m *mockEncriptacion) Verificar(password, hash string) bool { return hash == "hash:"+password }

type mockAuthSvc struct {
	permiso bool
	err     error
}

func (m *mockAuthSvc) TienePermiso(ctx context.Context, usuarioID, tenantID, permiso string) (bool, error) {
	return m.permiso, m.err
}

func credencialesValidas() *seguridad.CredencialesUsuario {
	return seguridad.NuevaCredencialesUsuarioDesdeBD("user-id-target", "hash:old", true, false, 0, time.Time{})
}

func TestResetearContrasenaExitoso(t *testing.T) {
	credRepo := &mockCredRepo{credenciales: credencialesValidas()}
	sesionRepo := &mockSesionRepo{}
	uc := resetpassword.NewResetearContrasenaCasoDeUso(credRepo, sesionRepo, &mockEncriptacion{hash: "hash:NuevoPass1!"}, &mockAuthSvc{permiso: true})

	resp, err := uc.Ejecutar(context.Background(), &resetpassword.ComandoResetearContrasena{
		UsuarioID: "user-id-target", NuevaPassword: "NuevoPass1!",
		TenantID: "tenant-id-1", EjecutorID: "admin-id-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.UsuarioID != "user-id-target" {
		t.Errorf("UsuarioID incorrecto: %v", resp.UsuarioID)
	}
	if resp.ModificadoEn == "" {
		t.Error("ModificadoEn vacío")
	}
	if !credRepo.actualizado {
		t.Error("credenciales no actualizadas")
	}
	if !sesionRepo.invalidado {
		t.Error("sesiones no invalidadas")
	}
}

func TestResetearContrasenaSinPermiso(t *testing.T) {
	uc := resetpassword.NewResetearContrasenaCasoDeUso(&mockCredRepo{}, &mockSesionRepo{}, &mockEncriptacion{}, &mockAuthSvc{permiso: false})
	_, err := uc.Ejecutar(context.Background(), &resetpassword.ComandoResetearContrasena{
		UsuarioID: "uid", NuevaPassword: "NuevoPass1!", TenantID: "tid", EjecutorID: "eid",
	})
	if !errors.Is(err, rbacdomain.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestResetearContrasenaPasswordDebil(t *testing.T) {
	uc := resetpassword.NewResetearContrasenaCasoDeUso(&mockCredRepo{}, &mockSesionRepo{}, &mockEncriptacion{}, &mockAuthSvc{permiso: true})
	_, err := uc.Ejecutar(context.Background(), &resetpassword.ComandoResetearContrasena{
		UsuarioID: "uid", NuevaPassword: "short", TenantID: "tid", EjecutorID: "eid",
	})
	if err == nil {
		t.Fatal("esperaba error por password débil")
	}
}

func TestResetearContrasenaUsuarioNoEncontrado(t *testing.T) {
	uc := resetpassword.NewResetearContrasenaCasoDeUso(
		&mockCredRepo{errObtener: errors.New("no encontrado")},
		&mockSesionRepo{}, &mockEncriptacion{}, &mockAuthSvc{permiso: true},
	)
	_, err := uc.Ejecutar(context.Background(), &resetpassword.ComandoResetearContrasena{
		UsuarioID: "uid", NuevaPassword: "NuevoPass1!", TenantID: "tid", EjecutorID: "eid",
	})
	if err == nil {
		t.Fatal("esperaba error por usuario no encontrado")
	}
}

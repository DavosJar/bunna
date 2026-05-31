package logout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/logout"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	"github.com/stretchr/testify/mock"
)

// ── Mocks ────────────────────────────────────────────────────────────────────

type mockUnitOfWork struct {
	sesionRepo sesiones_domain.SesionRepositorio
}

func (m *mockUnitOfWork) Transaccional(ctx context.Context, fn func(tx sesiones_domain.UnitOfWork) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fn(m)
}
func (m *mockUnitOfWork) SesionRepositorio() sesiones_domain.SesionRepositorio { return m.sesionRepo }
func (m *mockUnitOfWork) CredencialesRepositorio() seguridad_domain.CredencialesRepositorio {
	return nil
}
func (m *mockUnitOfWork) UsuarioRepositorio() usuario_domain.UsuarioRepositorio       { return nil }
func (m *mockUnitOfWork) EncriptacionServicio() seguridad_domain.EncriptacionServicio { return nil }
func (m *mockUnitOfWork) TokenServicio() sesiones_domain.TokenServicio                { return nil }
func (m *mockUnitOfWork) GeneradorID() shared_domain.GeneradorID                      { return nil }

// mockSesionRepo
type mockSesionRepo struct {
	mock.Mock
	sesionPorID   *sesiones_domain.Sesion
	errPorID      error
	errActualizar error
	sesionesLista []*sesiones_domain.Sesion
	actualizado   bool
}

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	return s, nil
}
func (m *mockSesionRepo) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	if m.errActualizar != nil {
		return nil, m.errActualizar
	}
	m.actualizado = true
	return s, nil
}
func (m *mockSesionRepo) ObtenerPorID(ctx context.Context, id string) (*sesiones_domain.Sesion, error) {
	return m.sesionPorID, m.errPorID
}
func (m *mockSesionRepo) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) ListarActivasPorUsuarioID(ctx context.Context, usuarioID string, ahora time.Time) ([]*sesiones_domain.Sesion, error) {
	return m.sesionesLista, nil
}
func (m *mockSesionRepo) Listar(ctx context.Context, spec sesiones_domain.EspecificacionSesion, pag shared_domain.Paginacion) ([]*sesiones_domain.Sesion, error) {
	args := m.Called(ctx, spec, pag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sesiones_domain.Sesion), args.Error(1)
}
func (m *mockSesionRepo) InvalidarTodasPorUsuarioID(ctx context.Context, usuarioID string) error {
	return nil
}
func (m *mockSesionRepo) Eliminar(ctx context.Context, id string) error { return nil }

// ── Helpers ───────────────────────────────────────────────────────────────────

// sesionActivaDeUsuario crea una sesión activa perteneciente al usuario dado.
func sesionActivaDeUsuario(sesionID, usuarioID string) *sesiones_domain.Sesion {
	ahora := time.Now()
	return sesiones_domain.NuevaSesionDesdeBD(
		sesionID, usuarioID, "hash:access", "hash:refresh",
		sesiones_domain.EstadoActiva, "192.168.1.1",
		ahora, ahora,
		ahora.Add(15*time.Minute), ahora.Add(24*time.Hour),
		ahora, 0,
	)
}

// ── Tests ────────────────────────────────────────────────────────────────────

// Escenario 1: Logout sesión específica exitoso
func TestLogout_SesionEspecifica(t *testing.T) {
	sesion := sesionActivaDeUsuario("sesion-id-1", "user-id-1")
	repo := &mockSesionRepo{sesionPorID: sesion}
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{sesionRepo: repo})

	resp, err := svc.Ejecutar(context.Background(), logout.ComandoLogout{
		SesionID:  "sesion-id-1",
		UsuarioID: "user-id-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.SesionesRevocadas != 1 {
		t.Errorf("esperaba 1 sesión revocada, got %d", resp.SesionesRevocadas)
	}
	if sesion.Estado() != sesiones_domain.EstadoRevocada {
		t.Errorf("esperaba estado REVOCADA, got %v", sesion.Estado())
	}
}

// Escenario 2: Logout de todas las sesiones
func TestLogout_CerrarTodas(t *testing.T) {
	sesiones := []*sesiones_domain.Sesion{
		sesionActivaDeUsuario("s1", "user-id-1"),
		sesionActivaDeUsuario("s2", "user-id-1"),
		sesionActivaDeUsuario("s3", "user-id-1"),
	}
	repo := &mockSesionRepo{sesionesLista: sesiones}
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{sesionRepo: repo})

	resp, err := svc.CerrarTodas(context.Background(), logout.ComandoCerrarTodas{
		UsuarioID: "user-id-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.SesionesRevocadas != 3 {
		t.Errorf("esperaba 3 sesiones revocadas, got %d", resp.SesionesRevocadas)
	}
	for _, s := range sesiones {
		if s.Estado() != sesiones_domain.EstadoRevocada {
			t.Errorf("sesión %v no fue revocada", s.ID())
		}
	}
}

// Escenario 4: Logout de sesión ya expirada → no-op
func TestLogout_SesionExpirada_NoOp(t *testing.T) {
	ahora := time.Now()
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh",
		sesiones_domain.EstadoExpirada, "", ahora, ahora,
		ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora, 0,
	)
	repo := &mockSesionRepo{sesionPorID: sesion}
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{sesionRepo: repo})

	resp, err := svc.Ejecutar(context.Background(), logout.ComandoLogout{
		SesionID:  "sesion-id-1",
		UsuarioID: "user-id-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.SesionesRevocadas != 0 {
		t.Errorf("esperaba 0 sesiones revocadas en no-op, got %d", resp.SesionesRevocadas)
	}
}

// Escenario 5: Logout de sesión ya revocada → no-op
func TestLogout_SesionRevocada_NoOp(t *testing.T) {
	ahora := time.Now()
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh",
		sesiones_domain.EstadoRevocada, "", ahora, ahora,
		ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora, 0,
	)
	repo := &mockSesionRepo{sesionPorID: sesion}
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{sesionRepo: repo})

	resp, err := svc.Ejecutar(context.Background(), logout.ComandoLogout{
		SesionID:  "sesion-id-1",
		UsuarioID: "user-id-1",
	})
	if err != nil {
		t.Fatalf("error inesperado en no-op: %v", err)
	}
	if resp.SesionesRevocadas != 0 {
		t.Errorf("esperaba 0 sesiones revocadas en no-op, got %d", resp.SesionesRevocadas)
	}
}

// Escenario 6: Logout de sesión de otro usuario → no autorizado
func TestLogout_SesionDeOtroUsuario(t *testing.T) {
	sesion := sesionActivaDeUsuario("sesion-id-1", "user-id-otro")
	repo := &mockSesionRepo{sesionPorID: sesion}
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{sesionRepo: repo})

	_, err := svc.Ejecutar(context.Background(), logout.ComandoLogout{
		SesionID:  "sesion-id-1",
		UsuarioID: "user-id-1",
	})
	if !errors.Is(err, logout.ErrNoAutorizado) {
		t.Errorf("esperaba ErrNoAutorizado, got %v", err)
	}
}

// Escenario 7: Sesión no encontrada
func TestLogout_SesionNoEncontrada(t *testing.T) {
	repo := &mockSesionRepo{errPorID: errors.New("no encontrada")}
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{sesionRepo: repo})

	_, err := svc.Ejecutar(context.Background(), logout.ComandoLogout{
		SesionID:  "sesion-inexistente",
		UsuarioID: "user-id-1",
	})
	if !errors.Is(err, logout.ErrSesionNoEncontrada) {
		t.Errorf("esperaba ErrSesionNoEncontrada, got %v", err)
	}
}

// Escenario 8: Timeout de inactividad marca sesión como expirada
func TestLogout_TimeoutInactividad(t *testing.T) {
	hace1hora := time.Now().Add(-1 * time.Hour)
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh",
		sesiones_domain.EstadoActiva, "", hace1hora, hace1hora,
		time.Now().Add(15*time.Minute), time.Now().Add(24*time.Hour),
		hace1hora, 0,
	)
	repo := &mockSesionRepo{sesionPorID: sesion}
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{sesionRepo: repo})

	err := svc.VerificarTimeout(context.Background(), "sesion-id-1", 30*time.Minute)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if sesion.Estado() != sesiones_domain.EstadoExpirada {
		t.Errorf("esperaba estado EXPIRADA tras timeout, got %v", sesion.Estado())
	}
}

// Escenario 9: Timeout configurable respeta el valor dado
func TestLogout_TimeoutConfigurable(t *testing.T) {
	hace10min := time.Now().Add(-10 * time.Minute)
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh",
		sesiones_domain.EstadoActiva, "", hace10min, hace10min,
		time.Now().Add(15*time.Minute), time.Now().Add(24*time.Hour),
		hace10min, 0,
	)
	repo := &mockSesionRepo{sesionPorID: sesion}
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{sesionRepo: repo})

	// timeout de 30 min, solo han pasado 10 min → no debe expirar
	err := svc.VerificarTimeout(context.Background(), "sesion-id-1", 30*time.Minute)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if sesion.Estado() != sesiones_domain.EstadoActiva {
		t.Errorf("esperaba estado ACTIVA con timeout no excedido, got %v", sesion.Estado())
	}
}

// Validación: SesionID vacío
func TestLogout_SesionIDVacio(t *testing.T) {
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{})
	_, err := svc.Ejecutar(context.Background(), logout.ComandoLogout{
		SesionID:  "",
		UsuarioID: "user-id-1",
	})
	if !errors.Is(err, logout.ErrSesionIDRequerido) {
		t.Errorf("esperaba ErrSesionIDRequerido, got %v", err)
	}
}

// Validación: UsuarioID vacío en Ejecutar
func TestLogout_UsuarioIDVacio(t *testing.T) {
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{})
	_, err := svc.Ejecutar(context.Background(), logout.ComandoLogout{
		SesionID:  "sesion-id-1",
		UsuarioID: "",
	})
	if !errors.Is(err, logout.ErrUsuarioIDRequerido) {
		t.Errorf("esperaba ErrUsuarioIDRequerido, got %v", err)
	}
}

// Validación: UsuarioID vacío en CerrarTodas
func TestLogout_CerrarTodas_UsuarioIDVacio(t *testing.T) {
	svc := logout.NuevoServicioLogout(&mockUnitOfWork{})
	_, err := svc.CerrarTodas(context.Background(), logout.ComandoCerrarTodas{UsuarioID: ""})
	if !errors.Is(err, logout.ErrUsuarioIDRequerido) {
		t.Errorf("esperaba ErrUsuarioIDRequerido, got %v", err)
	}
}

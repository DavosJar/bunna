package logout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/stretchr/testify/mock"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUnitOfWork struct {
	sesionRepo    sesiones_domain.SesionRepositorio
	transaccional func(ctx context.Context, fn func(tx sesiones_domain.UnitOfWork) error) error
}

func (m *mockUnitOfWork) Transaccional(ctx context.Context, fn func(tx sesiones_domain.UnitOfWork) error) error {
	if m.transaccional != nil {
		return m.transaccional(ctx, fn)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return fn(m)
}
func (m *mockUnitOfWork) SesionRepositorio() sesiones_domain.SesionRepositorio             { return m.sesionRepo }
func (m *mockUnitOfWork) CredencialesRepositorio() seguridad_domain.CredencialesRepositorio { return nil }
func (m *mockUnitOfWork) UsuarioRepositorio() usuario_domain.UsuarioRepositorio             { return nil }
func (m *mockUnitOfWork) EncriptacionServicio() seguridad_domain.EncriptacionServicio       { return nil }
func (m *mockUnitOfWork) TokenServicio() sesiones_domain.TokenServicio                      { return nil }
func (m *mockUnitOfWork) GeneradorID() shared_domain.GeneradorID                            { return nil }

type mockSesionRepo struct {
	mock.Mock
	sesion      *sesiones_domain.Sesion
	sesiones    []*sesiones_domain.Sesion
	err         error
	actualizado bool
}

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	m.actualizado = true
	return s, nil
}
func (m *mockSesionRepo) ObtenerPorID(ctx context.Context, id string) (*sesiones_domain.Sesion, error) {
	return m.sesion, m.err
}
func (m *mockSesionRepo) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) ListarActivasPorUsuarioID(ctx context.Context, usuarioID string, ahora time.Time) ([]*sesiones_domain.Sesion, error) {
	return m.sesiones, m.err
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

func sesionActiva() *sesiones_domain.Sesion {
	s, _ := sesiones_domain.NuevaSesion(
		"sesion-id", "user-id-1", "access-hash", "refresh-hash",
		"10.0.0.1", time.Now(), time.Now().Add(15*time.Minute), time.Now().Add(24*time.Hour),
	)
	return s
}

func TestCerrarSesionExitoso(t *testing.T) {
	sesion := sesionActiva()
	repo := &mockSesionRepo{sesion: sesion}
	uow := &mockUnitOfWork{sesionRepo: repo}
	uc := logout.NewCerrarSesionCasoDeUso(uow)

	resp, err := uc.Ejecutar(context.Background(), logout.ComandoCerrarSesion{
		SesionID: "sesion-id", UsuarioID: "user-id-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.SesionesRevocadas != 1 {
		t.Errorf("esperaba 1 sesión revocada, got %d", resp.SesionesRevocadas)
	}
	if !repo.actualizado {
		t.Error("esperaba persistencia de la sesión")
	}
}

func TestCerrarSesionSesionIDVacio(t *testing.T) {
	uc := logout.NewCerrarSesionCasoDeUso(&mockUnitOfWork{})
	_, err := uc.Ejecutar(context.Background(), logout.ComandoCerrarSesion{
		SesionID: "", UsuarioID: "user-id-1",
	})
	if !errors.Is(err, logout.ErrSesionIDRequerido) {
		t.Errorf("esperaba ErrSesionIDRequerido, got %v", err)
	}
}

func TestCerrarSesionUsuarioIDVacio(t *testing.T) {
	uc := logout.NewCerrarSesionCasoDeUso(&mockUnitOfWork{})
	_, err := uc.Ejecutar(context.Background(), logout.ComandoCerrarSesion{
		SesionID: "sesion-id", UsuarioID: "",
	})
	if !errors.Is(err, logout.ErrUsuarioIDRequerido) {
		t.Errorf("esperaba ErrUsuarioIDRequerido, got %v", err)
	}
}

func TestCerrarSesionNoEncontrada(t *testing.T) {
	repo := &mockSesionRepo{err: errors.New("not found")}
	uow := &mockUnitOfWork{sesionRepo: repo}
	uc := logout.NewCerrarSesionCasoDeUso(uow)
	_, err := uc.Ejecutar(context.Background(), logout.ComandoCerrarSesion{
		SesionID: "no-existe", UsuarioID: "user-id-1",
	})
	if !errors.Is(err, logout.ErrSesionNoEncontrada) {
		t.Errorf("esperaba ErrSesionNoEncontrada, got %v", err)
	}
}

func TestCerrarSesionNoAutorizado(t *testing.T) {
	sesion := sesionActiva()
	repo := &mockSesionRepo{sesion: sesion}
	uow := &mockUnitOfWork{sesionRepo: repo}
	uc := logout.NewCerrarSesionCasoDeUso(uow)
	_, err := uc.Ejecutar(context.Background(), logout.ComandoCerrarSesion{
		SesionID: "sesion-id", UsuarioID: "otro-usuario",
	})
	if !errors.Is(err, logout.ErrNoAutorizado) {
		t.Errorf("esperaba ErrNoAutorizado, got %v", err)
	}
}

func TestCerrarSesionYaRevocada(t *testing.T) {
	sesion := sesionActiva()
	sesion.Revocar()
	repo := &mockSesionRepo{sesion: sesion}
	uow := &mockUnitOfWork{sesionRepo: repo}
	uc := logout.NewCerrarSesionCasoDeUso(uow)
	resp, err := uc.Ejecutar(context.Background(), logout.ComandoCerrarSesion{
		SesionID: "sesion-id", UsuarioID: "user-id-1",
	})
	if err != nil {
		t.Fatalf("no debería fallar al cerrar sesión ya revocada: %v", err)
	}
	if resp.SesionesRevocadas != 0 {
		t.Errorf("esperaba 0 sesiones revocadas para sesión ya revocada, got %d", resp.SesionesRevocadas)
	}
}

func TestCerrarTodasLasSesiones(t *testing.T) {
	s1, _ := sesiones_domain.NuevaSesion("id-1", "user-1", "ah1", "rh1", "ip1", time.Now(), time.Now().Add(15*time.Minute), time.Now().Add(24*time.Hour))
	s2, _ := sesiones_domain.NuevaSesion("id-2", "user-1", "ah2", "rh2", "ip2", time.Now(), time.Now().Add(15*time.Minute), time.Now().Add(24*time.Hour))
	s3, _ := sesiones_domain.NuevaSesion("id-3", "user-1", "ah3", "rh3", "ip3", time.Now(), time.Now().Add(15*time.Minute), time.Now().Add(24*time.Hour))

	repo := &mockSesionRepo{sesiones: []*sesiones_domain.Sesion{s1, s2, s3}}
	uow := &mockUnitOfWork{sesionRepo: repo}
	uc := logout.NewCerrarSesionCasoDeUso(uow)

	resp, err := uc.CerrarTodas(context.Background(), logout.ComandoCerrarTodasLasSesiones{UsuarioID: "user-1"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.SesionesRevocadas != 3 {
		t.Errorf("esperaba 3 sesiones revocadas, got %d", resp.SesionesRevocadas)
	}
}

func TestCerrarTodasUsuarioIDVacio(t *testing.T) {
	uc := logout.NewCerrarSesionCasoDeUso(&mockUnitOfWork{})
	_, err := uc.CerrarTodas(context.Background(), logout.ComandoCerrarTodasLasSesiones{UsuarioID: ""})
	if !errors.Is(err, logout.ErrUsuarioIDRequerido) {
		t.Errorf("esperaba ErrUsuarioIDRequerido, got %v", err)
	}
}

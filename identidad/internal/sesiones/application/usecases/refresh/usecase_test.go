package refresh_test

import (
	"context"
	"errors"
	"testing"
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	"github.com/stretchr/testify/mock"
)

type mockUnitOfWork struct {
	sesionRepo    sesiones_domain.SesionRepositorio
	credRepo      seguridad_domain.CredencialesRepositorio
	usuarioRepo   usuario_domain.UsuarioRepositorio
	encriptacion  seguridad_domain.EncriptacionServicio
	tokenServicio sesiones_domain.TokenServicio
	generadorID   shared_domain.GeneradorID
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
func (m *mockUnitOfWork) SesionRepositorio() sesiones_domain.SesionRepositorio { return m.sesionRepo }
func (m *mockUnitOfWork) CredencialesRepositorio() seguridad_domain.CredencialesRepositorio {
	return m.credRepo
}
func (m *mockUnitOfWork) UsuarioRepositorio() usuario_domain.UsuarioRepositorio { return m.usuarioRepo }
func (m *mockUnitOfWork) EncriptacionServicio() seguridad_domain.EncriptacionServicio {
	return m.encriptacion
}
func (m *mockUnitOfWork) TokenServicio() sesiones_domain.TokenServicio { return m.tokenServicio }
func (m *mockUnitOfWork) GeneradorID() shared_domain.GeneradorID       { return m.generadorID }

type mockSesionRepo struct {
	mock.Mock
	sesion   *sesiones_domain.Sesion
	err      error
	errInval error
	invocado bool
}

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	return s, nil
}
func (m *mockSesionRepo) ObtenerPorID(ctx context.Context, id string) (*sesiones_domain.Sesion, error) {
	return m.sesion, m.err
}
func (m *mockSesionRepo) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones_domain.Sesion, error) {
	return m.sesion, m.err
}
func (m *mockSesionRepo) ListarActivasPorUsuarioID(ctx context.Context, usuarioID string, ahora time.Time) ([]*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) Listar(ctx context.Context, spec sesiones_domain.EspecificacionSesion, pag shared_domain.Paginacion) ([]*sesiones_domain.Sesion, error) {
	args := m.Called(ctx, spec, pag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sesiones_domain.Sesion), args.Error(1)
}
func (m *mockSesionRepo) InvalidarTodasPorUsuarioID(ctx context.Context, usuarioID string) error {
	m.invocado = true
	return m.errInval
}
func (m *mockSesionRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockTokenServicio struct {
	claims     *sesiones_domain.TokenClaims
	errValid   error
	accessTok  string
	refreshTok string
	expAccess  time.Time
	expRefresh time.Time
}

func (m *mockTokenServicio) GenerarAccessToken(usuarioID, sesionID string, tenantID string, rol string) (string, time.Time, error) {
	return m.accessTok, m.expAccess, nil
}
func (m *mockTokenServicio) GenerarRefreshToken(usuarioID, sesionID string) (string, time.Time, error) {
	return m.refreshTok, m.expRefresh, nil
}
func (m *mockTokenServicio) ValidarAccessToken(token string) (*sesiones_domain.TokenClaims, error) {
	return nil, nil
}
func (m *mockTokenServicio) ValidarRefreshToken(token string) (*sesiones_domain.TokenClaims, error) {
	return m.claims, m.errValid
}
func (m *mockTokenServicio) HashearToken(token string) string { return "hash:" + token }

type mockGeneradorID struct{}

func (m *mockGeneradorID) NextID(ctx context.Context) (string, error) { return "id", nil }

func sesionActiva() *sesiones_domain.Sesion {
	s, _ := sesiones_domain.NuevaSesion(
		"sesion-id", "user-id-1", "access-hash", "refresh-hash",
		"10.0.0.1", time.Now().Add(-1*time.Hour),
		time.Now().Add(14*time.Hour), time.Now().Add(7*24*time.Hour),
	)
	return s
}

func TestRenovarSesionExitoso(t *testing.T) {
	ahora := time.Now()
	sesion := sesionActiva()
	tokenSvc := &mockTokenServicio{
		claims:    &sesiones_domain.TokenClaims{UsuarioID: "user-id-1", SesionID: "sesion-id"},
		accessTok: "new-access", expAccess: ahora.Add(15 * time.Minute),
		refreshTok: "new-refresh", expRefresh: ahora.Add(24 * time.Hour),
	}
	uow := &mockUnitOfWork{
		sesionRepo:    &mockSesionRepo{sesion: sesion},
		tokenServicio: tokenSvc,
		generadorID:   &mockGeneradorID{},
	}
	uc := refresh.NewRenovarSesionCasoDeUso(uow, refresh.ConfigRefresh{MaxRefrescos: 10}, nil, nil)
	resp, err := uc.Ejecutar(context.Background(), refresh.ComandoRenovarSesion{RefreshToken: "valid-token"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.AccessToken != "new-access" || resp.RefreshToken != "new-refresh" {
		t.Error("tokens incorrectos en respuesta")
	}
}

func TestRenovarSesionTokenVacio(t *testing.T) {
	uc := refresh.NewRenovarSesionCasoDeUso(&mockUnitOfWork{}, refresh.ConfigRefresh{}, nil, nil)
	_, err := uc.Ejecutar(context.Background(), refresh.ComandoRenovarSesion{RefreshToken: ""})
	if !errors.Is(err, refresh.ErrRefreshTokenRequerido) {
		t.Errorf("esperaba ErrRefreshTokenRequerido, got %v", err)
	}
}

func TestRenovarSesionTokenInvalido(t *testing.T) {
	uow := &mockUnitOfWork{
		tokenServicio: &mockTokenServicio{errValid: errors.New("invalido")},
	}
	uc := refresh.NewRenovarSesionCasoDeUso(uow, refresh.ConfigRefresh{}, nil, nil)
	_, err := uc.Ejecutar(context.Background(), refresh.ComandoRenovarSesion{RefreshToken: "bad-token"})
	if !errors.Is(err, refresh.ErrTokenInvalido) {
		t.Errorf("esperaba ErrTokenInvalido, got %v", err)
	}
}

func TestRenovarSesionSesionNoValida(t *testing.T) {
	tokenSvc := &mockTokenServicio{
		claims: &sesiones_domain.TokenClaims{UsuarioID: "user-id-1"},
	}
	uow := &mockUnitOfWork{
		sesionRepo:    &mockSesionRepo{err: errors.New("not found")},
		tokenServicio: tokenSvc,
	}
	sesionRepoMock := &mockSesionRepo{err: errors.New("not found")}
	uow.sesionRepo = sesionRepoMock

	uc := refresh.NewRenovarSesionCasoDeUso(uow, refresh.ConfigRefresh{}, nil, nil)
	_, err := uc.Ejecutar(context.Background(), refresh.ComandoRenovarSesion{RefreshToken: "valid-token"})
	if !errors.Is(err, refresh.ErrTokenInvalido) {
		t.Errorf("esperaba ErrTokenInvalido cuando no se encuentra sesión, got %v", err)
	}
	if !sesionRepoMock.invocado {
		t.Error("esperaba que se invalidaran todas las sesiones (posible robo)")
	}
}

func TestRenovarSesionLimiteRefrescos(t *testing.T) {
	sesion := sesionActiva()
	tokenSvc := &mockTokenServicio{
		claims: &sesiones_domain.TokenClaims{UsuarioID: "user-id-1", SesionID: "sesion-id"},
	}
	uow := &mockUnitOfWork{
		sesionRepo:    &mockSesionRepo{sesion: sesion},
		tokenServicio: tokenSvc,
	}
	uc := refresh.NewRenovarSesionCasoDeUso(uow, refresh.ConfigRefresh{MaxRefrescos: 0}, nil, nil)
	_, err := uc.Ejecutar(context.Background(), refresh.ComandoRenovarSesion{RefreshToken: "valid-token"})
	if err != nil {
		t.Fatalf("no debería fallar con MaxRefrescos=0: %v", err)
	}
}

func TestRenovarSesionTimeoutAbsoluto(t *testing.T) {
	ahora := time.Now()
	sesion := sesionActiva()
	tokenSvc := &mockTokenServicio{
		claims:    &sesiones_domain.TokenClaims{UsuarioID: "user-id-1", SesionID: "sesion-id"},
		accessTok: "new-access", expAccess: ahora.Add(15 * time.Minute),
		refreshTok: "new-refresh", expRefresh: ahora.Add(24 * time.Hour),
	}
	uow := &mockUnitOfWork{
		sesionRepo:    &mockSesionRepo{sesion: sesion},
		tokenServicio: tokenSvc,
		generadorID:   &mockGeneradorID{},
	}
	uc := refresh.NewRenovarSesionCasoDeUso(uow, refresh.ConfigRefresh{
		MaxRefrescos:    10,
		TimeoutAbsoluto: 30 * time.Minute,
	}, nil, nil)
	_, err := uc.Ejecutar(context.Background(), refresh.ComandoRenovarSesion{RefreshToken: "valid-token"})
	if !errors.Is(err, refresh.ErrSesionAbsolutaExpirada) {
		t.Errorf("esperaba ErrSesionAbsolutaExpirada, got %v", err)
	}
}

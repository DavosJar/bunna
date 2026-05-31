package login_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/login"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	"github.com/stretchr/testify/mock"
)

// ── Mocks ────────────────────────────────────────────────────────────────────

type mockUnitOfWork struct {
	sesionRepo    sesiones_domain.SesionRepositorio
	credRepo      seguridad_domain.CredencialesRepositorio
	usuarioRepo   usuario_domain.UsuarioRepositorio
	encriptacion  seguridad_domain.EncriptacionServicio
	tokenServicio sesiones_domain.TokenServicio
	generadorID   shared_domain.GeneradorID
}

func (m *mockUnitOfWork) Transaccional(ctx context.Context, fn func(tx sesiones_domain.UnitOfWork) error) error {
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

// mockUsuarioRepo
type mockUsuarioRepo struct {
	usuarios []*usuario_domain.Usuario
	err      error
}

func (m *mockUsuarioRepo) Crear(ctx context.Context, u *usuario_domain.Usuario) (*usuario_domain.Usuario, error) {
	return nil, nil
}
func (m *mockUsuarioRepo) Actualizar(ctx context.Context, u *usuario_domain.Usuario) (*usuario_domain.Usuario, error) {
	return nil, nil
}
func (m *mockUsuarioRepo) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepo) ObtenerPorID(ctx context.Context, id string) (*usuario_domain.Usuario, error) {
	return nil, nil
}
func (m *mockUsuarioRepo) Listar(ctx context.Context, spec usuario_domain.EspecificacionUsuario, pag shared_domain.Paginacion) ([]*usuario_domain.Usuario, error) {
	return m.usuarios, m.err
}

// mockCredencialesRepo
type mockCredencialesRepo struct {
	credenciales  *seguridad_domain.CredencialesUsuario
	err           error
	errActualizar error
	actualizado   bool
}

func (m *mockCredencialesRepo) Crear(ctx context.Context, c *seguridad_domain.CredencialesUsuario) (*seguridad_domain.CredencialesUsuario, error) {
	return nil, nil
}
func (m *mockCredencialesRepo) Actualizar(ctx context.Context, c *seguridad_domain.CredencialesUsuario) (*seguridad_domain.CredencialesUsuario, error) {
	if m.errActualizar != nil {
		return nil, m.errActualizar
	}
	m.actualizado = true
	return c, nil
}
func (m *mockCredencialesRepo) ObtenerPorUsuarioID(ctx context.Context, id string) (*seguridad_domain.CredencialesUsuario, error) {
	return m.credenciales, m.err
}
func (m *mockCredencialesRepo) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockCredencialesRepo) Find(ctx context.Context, spec seguridad_domain.EspecificacionCredenciales, pag shared_domain.Paginacion) ([]*seguridad_domain.CredencialesUsuario, error) {
	return nil, nil
}

// mockSesionRepo
type mockSesionRepo struct {
	mock.Mock
	err error
}

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	if m.err != nil {
		return nil, m.err
	}
	return s, nil
}
func (m *mockSesionRepo) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	return s, nil
}
func (m *mockSesionRepo) ObtenerPorID(ctx context.Context, id string) (*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones_domain.Sesion, error) {
	return nil, nil
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
	return nil
}
func (m *mockSesionRepo) Eliminar(ctx context.Context, id string) error { return nil }

// mockEncriptacion
type mockEncriptacion struct{}

func (m *mockEncriptacion) Hashear(password string) (string, error) {
	return "hash:" + password, nil
}
func (m *mockEncriptacion) Verificar(password, hash string) bool {
	return hash == "hash:"+password
}

// mockTokenServicio
type mockTokenServicio struct {
	failAccess  bool
	failRefresh bool
}

func (m *mockTokenServicio) GenerarAccessToken(usuarioID, sesionID string, tenantID string, rol string) (string, time.Time, error) {
	if m.failAccess {
		return "", time.Time{}, errors.New("fallo access token")
	}
	return "access-token", time.Now().Add(15 * time.Minute), nil
}
func (m *mockTokenServicio) GenerarRefreshToken(usuarioID, sesionID string) (string, time.Time, error) {
	if m.failRefresh {
		return "", time.Time{}, errors.New("fallo refresh token")
	}
	return "refresh-token", time.Now().Add(24 * time.Hour), nil
}
func (m *mockTokenServicio) ValidarAccessToken(token string) (*sesiones_domain.TokenClaims, error) {
	return nil, nil
}
func (m *mockTokenServicio) ValidarRefreshToken(token string) (*sesiones_domain.TokenClaims, error) {
	return nil, nil
}
func (m *mockTokenServicio) HashearToken(token string) string { return "hash:" + token }

// mockGeneradorID
type mockGeneradorID struct{}

func (m *mockGeneradorID) NextID(ctx context.Context) (string, error) {
	return "sesion-id-generado", nil
}

// mockIntentoIPBloqueoRepo
type mockIntentoIPBloqueoRepo struct {
	mock.Mock
	intento    *seguridad_domain.IntentoPorIP
	errObtener error
}

func (m *mockIntentoIPBloqueoRepo) ObtenerPorIP(ctx context.Context, ip string) (*seguridad_domain.IntentoPorIP, error) {
	return m.intento, m.errObtener
}
func (m *mockIntentoIPBloqueoRepo) Crear(ctx context.Context, i *seguridad_domain.IntentoPorIP) (*seguridad_domain.IntentoPorIP, error) {
	return i, nil
}
func (m *mockIntentoIPBloqueoRepo) Actualizar(ctx context.Context, i *seguridad_domain.IntentoPorIP) (*seguridad_domain.IntentoPorIP, error) {
	return i, nil
}
func (m *mockIntentoIPBloqueoRepo) Listar(ctx context.Context, spec seguridad_domain.EspecificacionIntentoIP, pag shared_domain.Paginacion) ([]*seguridad_domain.IntentoPorIP, error) {
	args := m.Called(ctx, spec, pag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*seguridad_domain.IntentoPorIP), args.Error(1)
}
func (m *mockIntentoIPBloqueoRepo) EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error {
	return nil
}

// mockGeneradorIDBloqueo
type mockGeneradorIDBloqueo struct{}

func (m *mockGeneradorIDBloqueo) NextID(ctx context.Context) (string, error) {
	return "id-bloqueo", nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func usuarioValido() *usuario_domain.Usuario {
	u, _ := usuario_domain.NuevoUsuario("user-id-1", "test@correo.com", "Juan", "Pérez", "0999999999")
	return u
}

func credencialesValidas() *seguridad_domain.CredencialesUsuario {
	return seguridad_domain.NuevaCredencialesUsuarioDesdeBD(
		"user-id-1", "hash:secreto", true, false, 0, time.Time{},
	)
}

func uowValido(sesionRepo *mockSesionRepo, credRepo *mockCredencialesRepo, usuarioRepo *mockUsuarioRepo, tokenSvc *mockTokenServicio) *mockUnitOfWork {
	return &mockUnitOfWork{
		sesionRepo:    sesionRepo,
		credRepo:      credRepo,
		usuarioRepo:   usuarioRepo,
		encriptacion:  &mockEncriptacion{},
		tokenServicio: tokenSvc,
		generadorID:   &mockGeneradorID{},
	}
}

// ── Tests ────────────────────────────────────────────────────────────────────

// Escenario 1: Login exitoso
func TestLogin_Exitoso(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: credencialesValidas()},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	resp, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "secreto",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Error("tokens vacíos en respuesta exitosa")
	}
	if resp.UsuarioID != "user-id-1" {
		t.Errorf("usuarioID incorrecto: %v", resp.UsuarioID)
	}
}

// Escenario 3: Email vacío
func TestLogin_EmailVacio(t *testing.T) {
	svc := login.NuevoServicioLogin(&mockUnitOfWork{}, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "", Password: "secreto"})
	if !errors.Is(err, login.ErrEmailRequerido) {
		t.Errorf("esperaba ErrEmailRequerido, got %v", err)
	}
}

// Escenario 4: Email mal formado
func TestLogin_EmailInvalido(t *testing.T) {
	svc := login.NuevoServicioLogin(&mockUnitOfWork{}, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "invalido", Password: "secreto"})
	if !errors.Is(err, login.ErrEmailInvalido) {
		t.Errorf("esperaba ErrEmailInvalido, got %v", err)
	}
}

// Escenario 5: Password vacío
func TestLogin_PasswordVacio(t *testing.T) {
	svc := login.NuevoServicioLogin(&mockUnitOfWork{}, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "test@correo.com", Password: ""})
	if !errors.Is(err, login.ErrPasswordRequerido) {
		t.Errorf("esperaba ErrPasswordRequerido, got %v", err)
	}
}

// Escenario 6: Email no registrado
func TestLogin_EmailNoRegistrado(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "noexiste@correo.com", Password: "secreto"})
	if !errors.Is(err, login.ErrCredencialesInvalidas) {
		t.Errorf("esperaba ErrCredencialesInvalidas, got %v", err)
	}
}

// Escenario 7: Cuenta bloqueada
func TestLogin_CuentaBloqueada(t *testing.T) {
	creds := seguridad_domain.NuevaCredencialesUsuarioDesdeBD(
		"user-id-1", "hash:secreto", true, false, 5, time.Now().Add(15*time.Minute),
	)
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: creds},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "test@correo.com", Password: "secreto"})
	if !errors.Is(err, login.ErrCuentaBloqueada) {
		t.Errorf("esperaba ErrCuentaBloqueada, got %v", err)
	}
}

// Escenario 9: Cuenta inactiva
func TestLogin_CuentaInactiva(t *testing.T) {
	creds := seguridad_domain.NuevaCredencialesUsuarioDesdeBD(
		"user-id-1", "hash:secreto", false, false, 0, time.Time{},
	)
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: creds},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "test@correo.com", Password: "secreto"})
	if !errors.Is(err, login.ErrCuentaInactiva) {
		t.Errorf("esperaba ErrCuentaInactiva, got %v", err)
	}
}

// Escenario 11: Password incorrecto incrementa intentos
func TestLogin_PasswordIncorrecto(t *testing.T) {
	credRepo := &mockCredencialesRepo{credenciales: credencialesValidas()}
	uow := uowValido(
		&mockSesionRepo{},
		credRepo,
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "test@correo.com", Password: "incorrecta"})
	if !errors.Is(err, login.ErrCredencialesInvalidas) {
		t.Errorf("esperaba ErrCredencialesInvalidas, got %v", err)
	}
	if !credRepo.actualizado {
		t.Error("esperaba que se actualizaran las credenciales con el intento fallido")
	}
}

// Escenario 14: Fallo al crear sesión
func TestLogin_FalloAlCrearSesion(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{err: errors.New("fallo bd")},
		&mockCredencialesRepo{credenciales: credencialesValidas()},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "test@correo.com", Password: "secreto"})
	if err == nil {
		t.Error("esperaba error al fallar la creación de sesión")
	}
}

// Escenario 17: TokenServicio falla en access token
func TestLogin_FalloAccessToken(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: credencialesValidas()},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{failAccess: true},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "test@correo.com", Password: "secreto"})
	if !errors.Is(err, login.ErrErrorGenerandoTokens) {
		t.Errorf("esperaba ErrErrorGenerandoTokens, got %v", err)
	}
}

// Escenario 18: TokenServicio falla en refresh token
func TestLogin_FalloRefreshToken(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: credencialesValidas()},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{failRefresh: true},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{Email: "test@correo.com", Password: "secreto"})
	if !errors.Is(err, login.ErrErrorGenerandoTokens) {
		t.Errorf("esperaba ErrErrorGenerandoTokens, got %v", err)
	}
}

// Escenario 2: Login tras reintentos previos
func TestLogin_LoginTrasReintentos(t *testing.T) {
	creds := seguridad_domain.NuevaCredencialesUsuarioDesdeBD(
		"user-id-1", "hash:secreto", true, false, 3, time.Time{},
	)
	credRepo := &mockCredencialesRepo{credenciales: creds}
	uow := uowValido(
		&mockSesionRepo{},
		credRepo,
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	resp, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "secreto",
	})
	if err != nil {
		t.Fatalf("esperaba login exitoso tras reintentos, got %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("esperaba access token en respuesta")
	}
	if !credRepo.actualizado {
		t.Error("esperaba que se resetearan los intentos fallidos")
	}
}

// Escenario 8: Bloqueo expirado permite login
func TestLogin_BloqueoExpirado(t *testing.T) {
	creds := seguridad_domain.NuevaCredencialesUsuarioDesdeBD(
		"user-id-1", "hash:secreto", true, false, 5, time.Now().Add(-1*time.Hour),
	)
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: creds},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	resp, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "secreto",
	})
	if err != nil {
		t.Fatalf("esperaba login permitido con bloqueo expirado, got %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("esperaba access token en respuesta")
	}
}

// Escenario 10: Correo no verificado permite login
func TestLogin_CorreoNoVerificado(t *testing.T) {
	creds := seguridad_domain.NuevaCredencialesUsuarioDesdeBD(
		"user-id-1", "hash:secreto", true, false, 0, time.Time{},
	)
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: creds},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	resp, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "secreto",
	})
	if err != nil {
		t.Fatalf("esperaba login permitido con correo no verificado, got %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("esperaba access token en respuesta")
	}
}

// Escenario 12: 5to intento incorrecto bloquea la cuenta
func TestLogin_5toIntentoBloquea(t *testing.T) {
	creds := seguridad_domain.NuevaCredencialesUsuarioDesdeBD(
		"user-id-1", "hash:secreto", true, false, 4, time.Time{},
	)
	credRepo := &mockCredencialesRepo{credenciales: creds}
	uow := uowValido(
		&mockSesionRepo{},
		credRepo,
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "incorrecta",
	})
	if !errors.Is(err, login.ErrCredencialesInvalidas) {
		t.Errorf("esperaba ErrCredencialesInvalidas, got %v", err)
	}
	if !credRepo.actualizado {
		t.Error("esperaba que se actualizaran las credenciales")
	}
	if !creds.EstaBloqueado(time.Now()) {
		t.Error("esperaba que la cuenta quedara bloqueada tras el 5to intento")
	}
}

// Escenario 13: Intento en cuenta bloqueada no incrementa contador
func TestLogin_IntentoEnCuentaBloqueada(t *testing.T) {
	creds := seguridad_domain.NuevaCredencialesUsuarioDesdeBD(
		"user-id-1", "hash:secreto", true, false, 5, time.Now().Add(15*time.Minute),
	)
	intentosAntes := creds.IntentosFallidos()
	credRepo := &mockCredencialesRepo{credenciales: creds}
	uow := uowValido(
		&mockSesionRepo{},
		credRepo,
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "incorrecta",
	})
	if !errors.Is(err, login.ErrCuentaBloqueada) {
		t.Errorf("esperaba ErrCuentaBloqueada, got %v", err)
	}
	if creds.IntentosFallidos() != intentosAntes {
		t.Errorf("contador no debería incrementarse: antes=%d, después=%d",
			intentosAntes, creds.IntentosFallidos())
	}
	if credRepo.actualizado {
		t.Error("no debería actualizarse credenciales si la cuenta ya está bloqueada")
	}
}

// Escenario 15: Fallo al actualizar credenciales → rollback
func TestLogin_FalloAlActualizarCredenciales(t *testing.T) {
	credRepo := &mockCredencialesRepo{
		credenciales:  credencialesValidas(),
		errActualizar: errors.New("fallo bd al actualizar"),
	}
	uow := uowValido(
		&mockSesionRepo{},
		credRepo,
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "secreto",
	})
	if err == nil {
		t.Error("esperaba error al fallar la actualización de credenciales")
	}
}

// Escenario 16: Context cancelado → rollback
func TestLogin_ContextCancelado(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: credencialesValidas()},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, nil, nil)
	_, err := svc.Ejecutar(ctx, login.ComandoLogin{
		Email: "test@correo.com", Password: "secreto",
	})
	if err == nil {
		t.Error("esperaba error con context cancelado")
	}
}

// Escenario IP bloqueada impide login
func TestLogin_IPBloqueada(t *testing.T) {
	ahora := time.Now()
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 20, ahora.Add(-5*time.Minute), ahora.Add(30*time.Minute),
	)
	repo := &mockIntentoIPBloqueoRepo{intento: intento}
	bloqueoSvc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorIDBloqueo{}, bloqueo_ip.ConfigBloqueoIP{
		MaxIntentos: 20,
		Ventana:     15 * time.Minute,
		Duracion:    30 * time.Minute,
	})
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: credencialesValidas()},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, bloqueoSvc, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "secreto", IPOrigen: "10.0.0.1",
	})
	if !errors.Is(err, bloqueo_ip.ErrIPBloqueada) {
		t.Errorf("esperaba ErrIPBloqueada, got %v", err)
	}
}

// Escenario IP con intentos pero login exitoso → contador NO se resetea
func TestLogin_IPConIntentosNoSeReset(t *testing.T) {
	ahora := time.Now()
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 5, ahora, time.Time{},
	)
	repo := &mockIntentoIPBloqueoRepo{intento: intento}
	bloqueoSvc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorIDBloqueo{}, bloqueo_ip.ConfigBloqueoIP{
		MaxIntentos: 20,
		Ventana:     15 * time.Minute,
		Duracion:    30 * time.Minute,
	})
	uow := uowValido(
		&mockSesionRepo{},
		&mockCredencialesRepo{credenciales: credencialesValidas()},
		&mockUsuarioRepo{usuarios: []*usuario_domain.Usuario{usuarioValido()}},
		&mockTokenServicio{},
	)
	svc := login.NuevoServicioLogin(uow, bloqueoSvc, nil)
	_, err := svc.Ejecutar(context.Background(), login.ComandoLogin{
		Email: "test@correo.com", Password: "secreto", IPOrigen: "10.0.0.1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if intento.Contador() != 5 {
		t.Errorf("esperaba contador IP=5 sin reset, got %d", intento.Contador())
	}
}

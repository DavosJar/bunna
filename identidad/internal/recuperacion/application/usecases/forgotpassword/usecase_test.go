package forgotpassword_test

import (
	"context"
	"errors"
	"testing"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	"github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/forgotpassword"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockTokenRepo struct {
	token      *recuperacion.TokenRecuperacion
	err        error
}

func (m *mockTokenRepo) Crear(ctx context.Context, token *recuperacion.TokenRecuperacion) error { return m.err }
func (m *mockTokenRepo) ObtenerPorHash(ctx context.Context, hash string) (*recuperacion.TokenRecuperacion, error) {
	return m.token, m.err
}
func (m *mockTokenRepo) Actualizar(ctx context.Context, token *recuperacion.TokenRecuperacion) error { return m.err }

type mockUsuarioRepo struct {
	usuario *recuperacion.UsuarioRecuperacion
	err     error
}

func (m *mockUsuarioRepo) ObtenerPorCorreo(ctx context.Context, correo string) (*recuperacion.UsuarioRecuperacion, error) {
	return m.usuario, m.err
}
func (m *mockUsuarioRepo) ActualizarPassword(ctx context.Context, usuarioID, nuevoHash string) error { return m.err }

type mockSesionRepo struct {
	errInvalidar error
}

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) { return s, nil }
func (m *mockSesionRepo) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) { return s, nil }
func (m *mockSesionRepo) ObtenerPorID(ctx context.Context, id string) (*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) ListarActivasPorUsuarioID(ctx context.Context, uid string, ahora time.Time) ([]*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) Listar(ctx context.Context, spec sesiones_domain.EspecificacionSesion, pag shareddomain.Paginacion) ([]*sesiones_domain.Sesion, error) { return nil, nil }
func (m *mockSesionRepo) InvalidarTodasPorUsuarioID(ctx context.Context, uid string) error { return m.errInvalidar }
func (m *mockSesionRepo) Eliminar(ctx context.Context, id string) error { return nil }

type mockCredRepo struct {
	credenciales  *seguridad_domain.CredencialesUsuario
	errActualizar error
}

func (m *mockCredRepo) Crear(ctx context.Context, c *seguridad_domain.CredencialesUsuario) (*seguridad_domain.CredencialesUsuario, error) { return c, nil }
func (m *mockCredRepo) Actualizar(ctx context.Context, c *seguridad_domain.CredencialesUsuario) (*seguridad_domain.CredencialesUsuario, error) {
	if m.errActualizar != nil { return nil, m.errActualizar }
	return c, nil
}
func (m *mockCredRepo) ObtenerPorUsuarioID(ctx context.Context, id string) (*seguridad_domain.CredencialesUsuario, error) { return m.credenciales, nil }
func (m *mockCredRepo) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockCredRepo) Find(ctx context.Context, spec seguridad_domain.EspecificacionCredenciales, pag shareddomain.Paginacion) ([]*seguridad_domain.CredencialesUsuario, error) { return nil, nil }

type mockEncriptacion struct{}

func (m *mockEncriptacion) Hashear(password string) (string, error) { return "hash:" + password, nil }
func (m *mockEncriptacion) Verificar(password, hash string) bool { return hash == "hash:"+password }

type mockGenID struct{ id string }

func (m *mockGenID) NextID(ctx context.Context) (string, error) {
	return m.id, nil
}

func configValida() forgotpassword.ConfigRecuperacion {
	return forgotpassword.ConfigRecuperacion{
		TokenExpiracion:     time.Hour,
		RateLimitIPMax:      3,
		RateLimitUsuarioMax: 1,
		RateLimitVentana:    15 * time.Minute,
		FrontendURL:         "http://localhost:5173",
	}
}

func tokenValido() *recuperacion.TokenRecuperacion {
	return recuperacion.NuevoTokenRecuperacion("tok-1", "user-1", "token-plano", time.Now().Add(30*time.Minute))
}

func TestSolicitarRecuperacionExitoso(t *testing.T) {
	tokenRepo := &mockTokenRepo{}
	usuarioRepo := &mockUsuarioRepo{usuario: &recuperacion.UsuarioRecuperacion{
		ID: "user-1", Nombre: "Juan", Correo: "juan@test.com",
	}}
	emailSvc := &notificaciones.MockEmailServicio{}
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		tokenRepo, usuarioRepo, &mockSesionRepo{}, &mockCredRepo{},
		&mockEncriptacion{}, emailSvc, &mockGenID{id: "tok-1"}, configValida(),
	)

	resp, err := uc.Solicitar(context.Background(), forgotpassword.ComandoSolicitarRecuperacion{
		Email: "juan@test.com", IPOrigen: "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.Mensaje == "" {
		t.Error("Mensaje vacío")
	}
}

func TestSolicitarRecuperacionEmailVacio(t *testing.T) {
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		&mockTokenRepo{}, &mockUsuarioRepo{}, &mockSesionRepo{}, &mockCredRepo{},
		&mockEncriptacion{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Solicitar(context.Background(), forgotpassword.ComandoSolicitarRecuperacion{Email: "", IPOrigen: "127.0.0.1"})
	if !errors.Is(err, recuperacion.ErrEmailRequerido) {
		t.Errorf("esperaba ErrEmailRequerido, got %v", err)
	}
}

func TestSolicitarRecuperacionEmailInvalido(t *testing.T) {
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		&mockTokenRepo{}, &mockUsuarioRepo{}, &mockSesionRepo{}, &mockCredRepo{},
		&mockEncriptacion{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Solicitar(context.Background(), forgotpassword.ComandoSolicitarRecuperacion{Email: "not-an-email", IPOrigen: "127.0.0.1"})
	if !errors.Is(err, recuperacion.ErrEmailInvalido) {
		t.Errorf("esperaba ErrEmailInvalido, got %v", err)
	}
}

func TestSolicitarRecuperacionUsuarioNoEncontrado(t *testing.T) {
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		&mockTokenRepo{},
		&mockUsuarioRepo{err: errors.New("no encontrado")},
		&mockSesionRepo{}, &mockCredRepo{}, &mockEncriptacion{},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	resp, err := uc.Solicitar(context.Background(), forgotpassword.ComandoSolicitarRecuperacion{
		Email: "nadie@test.com", IPOrigen: "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	_ = resp
}

func TestValidarTokenExitoso(t *testing.T) {
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		&mockTokenRepo{token: tokenValido()}, &mockUsuarioRepo{}, &mockSesionRepo{},
		&mockCredRepo{}, &mockEncriptacion{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	resp, err := uc.ValidarToken(context.Background(), forgotpassword.ComandoValidarTokenRecuperacion{Token: "token-plano"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !resp.Valido {
		t.Error("esperaba Valido = true")
	}
	if resp.UsuarioID != "user-1" {
		t.Errorf("UsuarioID incorrecto: %v", resp.UsuarioID)
	}
}

func TestValidarTokenVacio(t *testing.T) {
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		&mockTokenRepo{}, &mockUsuarioRepo{}, &mockSesionRepo{},
		&mockCredRepo{}, &mockEncriptacion{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.ValidarToken(context.Background(), forgotpassword.ComandoValidarTokenRecuperacion{Token: ""})
	if !errors.Is(err, recuperacion.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestValidarTokenInvalido(t *testing.T) {
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		&mockTokenRepo{err: errors.New("no encontrado")}, &mockUsuarioRepo{}, &mockSesionRepo{},
		&mockCredRepo{}, &mockEncriptacion{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.ValidarToken(context.Background(), forgotpassword.ComandoValidarTokenRecuperacion{Token: "bad"})
	if !errors.Is(err, recuperacion.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestConfirmarRestablecimientoExitoso(t *testing.T) {
	tokenRepo := &mockTokenRepo{token: tokenValido()}
	usuarioRepo := &mockUsuarioRepo{usuario: &recuperacion.UsuarioRecuperacion{ID: "user-1"}}
	sesionRepo := &mockSesionRepo{}
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		tokenRepo, usuarioRepo, sesionRepo, &mockCredRepo{},
		&mockEncriptacion{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	resp, err := uc.Confirmar(context.Background(), forgotpassword.ComandoConfirmarRestablecimiento{
		Token: "token-plano", NuevaPassword: "NuevoPass1!",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.Mensaje == "" {
		t.Error("Mensaje vacío")
	}
}

func TestConfirmarRestablecimientoPasswordDebil(t *testing.T) {
	uc := forgotpassword.NewRecuperarContrasenaCasoDeUso(
		&mockTokenRepo{}, &mockUsuarioRepo{}, &mockSesionRepo{}, &mockCredRepo{},
		&mockEncriptacion{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Confirmar(context.Background(), forgotpassword.ComandoConfirmarRestablecimiento{
		Token: "tok", NuevaPassword: "short",
	})
	if !errors.Is(err, recuperacion.ErrPasswordDebil) {
		t.Errorf("esperaba ErrPasswordDebil, got %v", err)
	}
}

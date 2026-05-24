package facades_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	uc_sesiones_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	uc_sesiones_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	uc_sesiones_refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	svc_registro "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/services/registro"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockEjecutorRegistro struct {
	respuesta *svc_registro.DtoRespuestaRegistro
	err       error
}

func (m *mockEjecutorRegistro) Ejecutar(ctx context.Context, cmd *svc_registro.ComandoRegistro) (*svc_registro.DtoRespuestaRegistro, error) {
	return m.respuesta, m.err
}

type mockLoginUseCase struct {
	respuesta *uc_sesiones_login.RespuestaIniciarSesion
	err       error
}

func (m *mockLoginUseCase) Ejecutar(ctx context.Context, cmd uc_sesiones_login.ComandoIniciarSesion) (*uc_sesiones_login.RespuestaIniciarSesion, error) {
	return m.respuesta, m.err
}

type mockRefreshUseCase struct {
	respuesta *uc_sesiones_refresh.RespuestaRenovarSesion
	err       error
}

func (m *mockRefreshUseCase) Ejecutar(ctx context.Context, cmd uc_sesiones_refresh.ComandoRenovarSesion) (*uc_sesiones_refresh.RespuestaRenovarSesion, error) {
	return m.respuesta, m.err
}

type mockLogoutUseCase struct {
	respuesta *uc_sesiones_logout.RespuestaCerrarSesion
	err       error
}

func (m *mockLogoutUseCase) Ejecutar(ctx context.Context, cmd uc_sesiones_logout.ComandoCerrarSesion) (*uc_sesiones_logout.RespuestaCerrarSesion, error) {
	return m.respuesta, m.err
}

func (m *mockLogoutUseCase) CerrarTodas(ctx context.Context, cmd uc_sesiones_logout.ComandoCerrarTodasLasSesiones) (*uc_sesiones_logout.RespuestaCerrarSesion, error) {
	return m.respuesta, m.err
}

// newAuthFacadeMock es un helper para tests que crea una AuthFacade con mocks.
func newAuthFacadeMock(reg svc_registro.EjecutorRegistro, login facades.LoginUseCase, refresh facades.RefreshUseCase, logout facades.LogoutUseCase) facades.AuthFacade {
	return facades.NewAuthFacade(reg, login, refresh, logout)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func respuestaRegistroValida() *svc_registro.DtoRespuestaRegistro {
	return &svc_registro.DtoRespuestaRegistro{
		UsuarioID: "usuario-id-1",
		Correo:    "test@correo.com",
		Estado:    "NO_VERIFICADO",
		Timestamp: time.Now(),
	}
}

func respuestaLoginValida() *uc_sesiones_login.RespuestaIniciarSesion {
	return &uc_sesiones_login.RespuestaIniciarSesion{
		AccessToken:       "access-token",
		RefreshToken:      "refresh-token",
		ExpiracionAccess:  time.Now().Add(15 * time.Minute),
		ExpiracionRefresh: time.Now().Add(24 * time.Hour),
		UsuarioID:         "usuario-id-1",
		SesionID:          "sesion-id-1",
	}
}

// ── Capturadores ──────────────────────────────────────────────────────────────

type mockCapturadorRegistro struct {
	cmdRecibido *svc_registro.ComandoRegistro
}

func (m *mockCapturadorRegistro) Ejecutar(ctx context.Context, cmd *svc_registro.ComandoRegistro) (*svc_registro.DtoRespuestaRegistro, error) {
	m.cmdRecibido = cmd
	return respuestaRegistroValida(), nil
}

type mockCapturadorLogin struct {
	cmdRecibido uc_sesiones_login.ComandoIniciarSesion
}

func (m *mockCapturadorLogin) Ejecutar(ctx context.Context, cmd uc_sesiones_login.ComandoIniciarSesion) (*uc_sesiones_login.RespuestaIniciarSesion, error) {
	m.cmdRecibido = cmd
	return respuestaLoginValida(), nil
}

// ── Tests Registrar ───────────────────────────────────────────────────────────

func TestAuthFacade_Registrar_Exitoso(t *testing.T) {
	facade := newAuthFacadeMock(
		&mockEjecutorRegistro{respuesta: respuestaRegistroValida()},
		&mockLoginUseCase{},
		&mockRefreshUseCase{},
		&mockLogoutUseCase{},
	)

	resp, err := facade.Registrar(context.Background(), facades.ComandoRegistro{
		Nombre:   "Juan",
		Apellido: "Pérez",
		Correo:   "test@correo.com",
		Password: "secreto123",
		Telefono: "0999999999",
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.UsuarioID == "" {
		t.Error("esperaba UsuarioID no vacío")
	}
	if resp.Correo != "test@correo.com" {
		t.Errorf("correo incorrecto: %v", resp.Correo)
	}
	if resp.Estado == "" {
		t.Error("esperaba Estado no vacío")
	}
}

func TestAuthFacade_Registrar_TraduccionComando(t *testing.T) {
	captura := &mockCapturadorRegistro{}
	facade := newAuthFacadeMock(captura, &mockLoginUseCase{}, &mockRefreshUseCase{}, &mockLogoutUseCase{})

	_, err := facade.Registrar(context.Background(), facades.ComandoRegistro{
		Nombre:   "Ana",
		Apellido: "García",
		Correo:   "ana@correo.com",
		Password: "pass123",
		Telefono: "0988888888",
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if captura.cmdRecibido.Nombre != "Ana" {
		t.Errorf("Nombre no traducido: %v", captura.cmdRecibido.Nombre)
	}
	if captura.cmdRecibido.Correo != "ana@correo.com" {
		t.Errorf("Correo no traducido: %v", captura.cmdRecibido.Correo)
	}
	if captura.cmdRecibido.Password != "pass123" {
		t.Errorf("Password no traducido: %v", captura.cmdRecibido.Password)
	}
}

func TestAuthFacade_Registrar_PropagaError(t *testing.T) {
	errEsperado := errors.New("correo ya registrado")
	facade := newAuthFacadeMock(
		&mockEjecutorRegistro{err: errEsperado},
		&mockLoginUseCase{},
		&mockRefreshUseCase{},
		&mockLogoutUseCase{},
	)

	resp, err := facade.Registrar(context.Background(), facades.ComandoRegistro{
		Nombre:   "Juan",
		Correo:   "test@correo.com",
		Password: "secreto123",
	})

	if !errors.Is(err, errEsperado) {
		t.Errorf("esperaba error %v, got %v", errEsperado, err)
	}
	if resp != nil {
		t.Error("esperaba respuesta nil cuando hay error")
	}
}

func TestAuthFacade_Registrar_ContextoCancelado(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	errEsperado := errors.New("context cancelado")
	facade := newAuthFacadeMock(
		&mockEjecutorRegistro{err: errEsperado},
		&mockLoginUseCase{},
		&mockRefreshUseCase{},
		&mockLogoutUseCase{},
	)

	_, err := facade.Registrar(ctx, facades.ComandoRegistro{
		Nombre:   "Juan",
		Correo:   "test@correo.com",
		Password: "secreto123",
	})

	if err == nil {
		t.Error("esperaba error con contexto cancelado")
	}
}

// ── Tests Login ───────────────────────────────────────────────────────────────

func TestAuthFacade_Login_Exitoso(t *testing.T) {
	facade := newAuthFacadeMock(
		&mockEjecutorRegistro{},
		&mockLoginUseCase{respuesta: respuestaLoginValida()},
		&mockRefreshUseCase{},
		&mockLogoutUseCase{},
	)

	resp, err := facade.Login(context.Background(), facades.ComandoLogin{
		Email:    "test@correo.com",
		Password: "secreto",
		IPOrigen: "127.0.0.1",
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("esperaba AccessToken no vacío")
	}
	if resp.RefreshToken == "" {
		t.Error("esperaba RefreshToken no vacío")
	}
	if resp.TokenType != "Bearer" {
		t.Errorf("esperaba TokenType=Bearer, got %v", resp.TokenType)
	}
	if resp.ExpiresIn <= 0 {
		t.Errorf("esperaba ExpiresIn > 0, got %v", resp.ExpiresIn)
	}
	if resp.UsuarioID == "" {
		t.Error("esperaba UsuarioID no vacío")
	}
	if resp.SesionID == "" {
		t.Error("esperaba SesionID no vacío")
	}
}

func TestAuthFacade_Login_PropagaError(t *testing.T) {
	errEsperado := errors.New("credenciales inválidas")
	facade := newAuthFacadeMock(
		&mockEjecutorRegistro{},
		&mockLoginUseCase{err: errEsperado},
		&mockRefreshUseCase{},
		&mockLogoutUseCase{},
	)

	resp, err := facade.Login(context.Background(), facades.ComandoLogin{
		Email:    "test@correo.com",
		Password: "incorrecta",
	})

	if !errors.Is(err, errEsperado) {
		t.Errorf("esperaba error %v, got %v", errEsperado, err)
	}
	if resp != nil {
		t.Error("esperaba respuesta nil cuando hay error")
	}
}

func TestAuthFacade_Login_ExpiresInEnSegundos(t *testing.T) {
	expiracion := time.Now().Add(900 * time.Second)
	facade := newAuthFacadeMock(
		&mockEjecutorRegistro{},
		&mockLoginUseCase{respuesta: &uc_sesiones_login.RespuestaIniciarSesion{
			AccessToken:       "token",
			RefreshToken:      "refresh",
			ExpiracionAccess:  expiracion,
			ExpiracionRefresh: time.Now().Add(24 * time.Hour),
			UsuarioID:         "uid",
			SesionID:          "sid",
		}},
		&mockRefreshUseCase{},
		&mockLogoutUseCase{},
	)

	resp, err := facade.Login(context.Background(), facades.ComandoLogin{
		Email:    "test@correo.com",
		Password: "secreto",
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.ExpiresIn < 898 || resp.ExpiresIn > 900 {
		t.Errorf("ExpiresIn fuera de rango (~900s), got %v", resp.ExpiresIn)
	}
}

func TestAuthFacade_Login_TraduccionIPOrigen(t *testing.T) {
	captura := &mockCapturadorLogin{}
	facade := newAuthFacadeMock(&mockEjecutorRegistro{}, captura, &mockRefreshUseCase{}, &mockLogoutUseCase{})

	_, err := facade.Login(context.Background(), facades.ComandoLogin{
		Email:    "test@correo.com",
		Password: "secreto",
		IPOrigen: "192.168.1.1",
	})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if captura.cmdRecibido.IPOrigen != "192.168.1.1" {
		t.Errorf("IPOrigen no traducida: %v", captura.cmdRecibido.IPOrigen)
	}
}

package solicitarrecuperacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	uc_solicitar "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/solicitarrecuperacion"
)

type mockTokenRepo struct {
	err error
}

func (m *mockTokenRepo) Crear(ctx context.Context, token *recuperacion.TokenRecuperacion) error {
	return m.err
}
func (m *mockTokenRepo) ObtenerPorHash(ctx context.Context, hash string) (*recuperacion.TokenRecuperacion, error) {
	return nil, m.err
}
func (m *mockTokenRepo) Actualizar(ctx context.Context, token *recuperacion.TokenRecuperacion) error {
	return m.err
}

type mockUsuarioRepo struct {
	usuario *recuperacion.UsuarioRecuperacion
	err     error
}

func (m *mockUsuarioRepo) ObtenerPorCorreo(ctx context.Context, correo string) (*recuperacion.UsuarioRecuperacion, error) {
	return m.usuario, m.err
}
func (m *mockUsuarioRepo) ActualizarPassword(ctx context.Context, usuarioID, nuevoHash string) error {
	return m.err
}

type mockGenID struct{ id string }

func (m *mockGenID) NextID(ctx context.Context) (string, error) {
	return m.id, nil
}

func configValida() uc_solicitar.ConfigRecuperacion {
	return uc_solicitar.ConfigRecuperacion{
		TokenExpiracion:     time.Hour,
		RateLimitIPMax:      3,
		RateLimitUsuarioMax: 1,
		RateLimitVentana:    15 * time.Minute,
		FrontendURL:         "http://localhost:5173",
	}
}

func TestSolicitarRecuperacionExitoso(t *testing.T) {
	tokenRepo := &mockTokenRepo{}
	usuarioRepo := &mockUsuarioRepo{usuario: &recuperacion.UsuarioRecuperacion{
		ID: "user-1", Nombre: "Juan", Correo: "juan@test.com",
	}}
	emailSvc := &notificaciones.MockEmailServicio{}
	uc := uc_solicitar.NewSolicitarRecuperacionCasoDeUso(
		tokenRepo, usuarioRepo, emailSvc, &mockGenID{id: "tok-1"}, configValida(),
	)

	resp, err := uc.Ejecutar(context.Background(), &uc_solicitar.ComandoSolicitarRecuperacion{
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
	uc := uc_solicitar.NewSolicitarRecuperacionCasoDeUso(
		&mockTokenRepo{}, &mockUsuarioRepo{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Ejecutar(context.Background(), &uc_solicitar.ComandoSolicitarRecuperacion{Email: "", IPOrigen: "127.0.0.1"})
	if !errors.Is(err, recuperacion.ErrEmailRequerido) {
		t.Errorf("esperaba ErrEmailRequerido, got %v", err)
	}
}

func TestSolicitarRecuperacionEmailInvalido(t *testing.T) {
	uc := uc_solicitar.NewSolicitarRecuperacionCasoDeUso(
		&mockTokenRepo{}, &mockUsuarioRepo{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Ejecutar(context.Background(), &uc_solicitar.ComandoSolicitarRecuperacion{Email: "not-an-email", IPOrigen: "127.0.0.1"})
	if !errors.Is(err, recuperacion.ErrEmailInvalido) {
		t.Errorf("esperaba ErrEmailInvalido, got %v", err)
	}
}

func TestSolicitarRecuperacionUsuarioNoEncontrado(t *testing.T) {
	uc := uc_solicitar.NewSolicitarRecuperacionCasoDeUso(
		&mockTokenRepo{},
		&mockUsuarioRepo{err: errors.New("no encontrado")},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	resp, err := uc.Ejecutar(context.Background(), &uc_solicitar.ComandoSolicitarRecuperacion{
		Email: "nadie@test.com", IPOrigen: "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	_ = resp
}

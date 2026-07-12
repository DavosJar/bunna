package solicitarverificacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	uc_solicitar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/solicitarverificacion"
	verificacion "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
)

type mockVerifRepo struct {
	usuario *verificacion.UsuarioVerificacion
	err     error
}

func (m *mockVerifRepo) ObtenerPorHashToken(ctx context.Context, hash string) (*verificacion.UsuarioVerificacion, error) {
	return m.usuario, m.err
}
func (m *mockVerifRepo) ActualizarPrueba(ctx context.Context, usuarioID string, prueba verificacion.PruebaVerificacion) error {
	return nil
}
func (m *mockVerifRepo) MarcarVerificado(ctx context.Context, usuarioID string) error {
	return m.err
}
func (m *mockVerifRepo) ObtenerPorID(ctx context.Context, usuarioID string) (*verificacion.UsuarioVerificacion, error) {
	return m.usuario, m.err
}
func (m *mockVerifRepo) ObtenerPorCorreo(ctx context.Context, correo string) (*verificacion.UsuarioVerificacion, error) {
	return m.usuario, m.err
}

type mockGenID struct{ id string }

func (m *mockGenID) NextID(ctx context.Context) (string, error) {
	return m.id, nil
}

func configValida() uc_solicitar.ConfigVerificacion {
	return uc_solicitar.ConfigVerificacion{
		TokenExpiracion: 24 * time.Hour,
		MaxReenvios:     5,
		VentanaReenvios: 24 * time.Hour,
		FrontendURL:     "http://localhost:5173",
	}
}

func usuarioNoVerificado() *verificacion.UsuarioVerificacion {
	return &verificacion.UsuarioVerificacion{
		ID:                 "user-1",
		Nombre:             "Juan",
		Correo:             "juan@test.com",
		EstadoVerificacion: "PENDIENTE",
		PruebaVerificacion: verificacion.PruebaVerificacionVacia(),
		ContadorReenvios:   0,
	}
}

func TestSolicitarVerificacionExitoso(t *testing.T) {
	repo := &mockVerifRepo{usuario: usuarioNoVerificado()}
	emailSvc := &notificaciones.MockEmailServicio{}
	uc := uc_solicitar.NewSolicitarVerificacionCasoDeUso(repo, emailSvc, &mockGenID{id: "token-1"}, configValida())

	resp, err := uc.Ejecutar(context.Background(), &uc_solicitar.ComandoSolicitarVerificacion{Correo: "juan@test.com"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.Mensaje == "" {
		t.Error("Mensaje vacío")
	}
	time.Sleep(10 * time.Millisecond)
	if len(emailSvc.EmailsEnviados) == 0 {
		t.Error("no se envió email")
	}
}

func TestSolicitarVerificacionUsuarioNoEncontrado(t *testing.T) {
	uc := uc_solicitar.NewSolicitarVerificacionCasoDeUso(
		&mockVerifRepo{err: errors.New("no encontrado")},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Ejecutar(context.Background(), &uc_solicitar.ComandoSolicitarVerificacion{Correo: "user-x@test.com"})
	if !errors.Is(err, verificacion.ErrUsuarioNoEncontrado) {
		t.Errorf("esperaba ErrUsuarioNoEncontrado, got %v", err)
	}
}

func TestSolicitarVerificacionYaVerificado(t *testing.T) {
	usuario := usuarioNoVerificado()
	usuario.EstadoVerificacion = "VERIFICADO"
	uc := uc_solicitar.NewSolicitarVerificacionCasoDeUso(
		&mockVerifRepo{usuario: usuario},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Ejecutar(context.Background(), &uc_solicitar.ComandoSolicitarVerificacion{Correo: "juan@test.com"})
	if !errors.Is(err, verificacion.ErrCorreoYaVerificado) {
		t.Errorf("esperaba ErrCorreoYaVerificado, got %v", err)
	}
}

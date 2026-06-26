package reenviarverificacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	uc_reenviar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/reenviarverificacion"
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

type mockGenID struct{ id string }

func (m *mockGenID) NextID(ctx context.Context) (string, error) {
	return m.id, nil
}

func configValida() uc_reenviar.ConfigVerificacion {
	return uc_reenviar.ConfigVerificacion{
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

func solicitarConfigValida() uc_solicitar.ConfigVerificacion {
	return uc_solicitar.ConfigVerificacion{
		TokenExpiracion: 24 * time.Hour,
		MaxReenvios:     5,
		VentanaReenvios: 24 * time.Hour,
		FrontendURL:     "http://localhost:5173",
	}
}

func TestReenviarVerificacionExitoso(t *testing.T) {
	usuario := usuarioNoVerificado()
	repo := &mockVerifRepo{usuario: usuario}
	emailSvc := &notificaciones.MockEmailServicio{}
	solicitarUC := uc_solicitar.NewSolicitarVerificacionCasoDeUso(repo, emailSvc, &mockGenID{id: "token-2"}, solicitarConfigValida())
	uc := uc_reenviar.NewReenviarVerificacionCasoDeUso(repo, emailSvc, &mockGenID{id: "token-2"}, solicitarUC, configValida())

	resp, err := uc.Ejecutar(context.Background(), &uc_reenviar.ComandoReenviarVerificacion{UsuarioID: "user-1"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.Mensaje == "" {
		t.Error("Mensaje vacío")
	}
	time.Sleep(10 * time.Millisecond)
	if len(emailSvc.EmailsEnviados) == 0 {
		t.Error("no se reenvió email")
	}
}

func TestReenviarVerificacionUsuarioNoEncontrado(t *testing.T) {
	emailSvc := &notificaciones.MockEmailServicio{}
	solicitarUC := uc_solicitar.NewSolicitarVerificacionCasoDeUso(
		&mockVerifRepo{err: errors.New("no encontrado")},
		emailSvc, &mockGenID{}, solicitarConfigValida(),
	)
	uc := uc_reenviar.NewReenviarVerificacionCasoDeUso(
		&mockVerifRepo{err: errors.New("no encontrado")},
		emailSvc, &mockGenID{}, solicitarUC, configValida(),
	)
	_, err := uc.Ejecutar(context.Background(), &uc_reenviar.ComandoReenviarVerificacion{UsuarioID: "user-x"})
	if !errors.Is(err, verificacion.ErrUsuarioNoEncontrado) {
		t.Errorf("esperaba ErrUsuarioNoEncontrado, got %v", err)
	}
}

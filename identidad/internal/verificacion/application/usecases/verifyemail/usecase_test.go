package verifyemail_test

import (
	"context"
	"errors"
	"testing"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/verifyemail"
	verificacion "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
)

type mockVerifRepo struct {
	usuario  *verificacion.UsuarioVerificacion
	err      error
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

func configValida() verifyemail.ConfigVerificacion {
	return verifyemail.ConfigVerificacion{
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
	uc := verifyemail.NewVerificarCorreoCasoDeUso(repo, emailSvc, &mockGenID{id: "token-1"}, configValida())

	resp, err := uc.Solicitar(context.Background(), verifyemail.ComandoSolicitarVerificacion{UsuarioID: "user-1"})
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
	uc := verifyemail.NewVerificarCorreoCasoDeUso(
		&mockVerifRepo{err: errors.New("no encontrado")},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Solicitar(context.Background(), verifyemail.ComandoSolicitarVerificacion{UsuarioID: "user-x"})
	if !errors.Is(err, verificacion.ErrUsuarioNoEncontrado) {
		t.Errorf("esperaba ErrUsuarioNoEncontrado, got %v", err)
	}
}

func TestSolicitarVerificacionYaVerificado(t *testing.T) {
	usuario := usuarioNoVerificado()
	usuario.EstadoVerificacion = "VERIFICADO"
	uc := verifyemail.NewVerificarCorreoCasoDeUso(
		&mockVerifRepo{usuario: usuario},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Solicitar(context.Background(), verifyemail.ComandoSolicitarVerificacion{UsuarioID: "user-1"})
	if !errors.Is(err, verificacion.ErrCorreoYaVerificado) {
		t.Errorf("esperaba ErrCorreoYaVerificado, got %v", err)
	}
}

func TestConfirmarVerificacionExitoso(t *testing.T) {
	usuario := usuarioNoVerificado()
	tokenPlano := "token-valido"
	usuario.PruebaVerificacion = verificacion.NuevaPruebaVerificacion(tokenPlano, time.Now().Add(24*time.Hour))
	repo := &mockVerifRepo{usuario: usuario}
	uc := verifyemail.NewVerificarCorreoCasoDeUso(repo, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida())

	resp, err := uc.Confirmar(context.Background(), verifyemail.ComandoConfirmarVerificacion{Token: tokenPlano})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.Mensaje == "" {
		t.Error("Mensaje vacío")
	}
}

func TestConfirmarVerificacionTokenVacio(t *testing.T) {
	uc := verifyemail.NewVerificarCorreoCasoDeUso(&mockVerifRepo{}, &notificaciones.MockEmailServicio{}, &mockGenID{}, configValida())
	_, err := uc.Confirmar(context.Background(), verifyemail.ComandoConfirmarVerificacion{Token: ""})
	if !errors.Is(err, verificacion.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestConfirmarVerificacionTokenInvalido(t *testing.T) {
	uc := verifyemail.NewVerificarCorreoCasoDeUso(
		&mockVerifRepo{err: errors.New("no encontrado")},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Confirmar(context.Background(), verifyemail.ComandoConfirmarVerificacion{Token: "bad-token"})
	if !errors.Is(err, verificacion.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestConfirmarVerificacionTokenExpirado(t *testing.T) {
	usuario := usuarioNoVerificado()
	usuario.PruebaVerificacion = verificacion.NuevaPruebaVerificacion("expirado", time.Now().Add(-48*time.Hour))
	uc := verifyemail.NewVerificarCorreoCasoDeUso(
		&mockVerifRepo{usuario: usuario},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Confirmar(context.Background(), verifyemail.ComandoConfirmarVerificacion{Token: "expirado"})
	if !errors.Is(err, verificacion.ErrEnlaceExpirado) {
		t.Errorf("esperaba ErrEnlaceExpirado, got %v", err)
	}
}

func TestConfirmarVerificacionFalloMarcar(t *testing.T) {
	usuario := usuarioNoVerificado()
	usuario.PruebaVerificacion = verificacion.NuevaPruebaVerificacion("token", time.Now().Add(24*time.Hour))
	uc := verifyemail.NewVerificarCorreoCasoDeUso(
		&mockVerifRepo{usuario: usuario, err: errors.New("fallo bd")},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Confirmar(context.Background(), verifyemail.ComandoConfirmarVerificacion{Token: "token"})
	if err == nil {
		t.Fatal("esperaba error al marcar verificado")
	}
}

func TestReenviarVerificacionExitoso(t *testing.T) {
	usuario := usuarioNoVerificado()
	repo := &mockVerifRepo{usuario: usuario}
	emailSvc := &notificaciones.MockEmailServicio{}
	uc := verifyemail.NewVerificarCorreoCasoDeUso(repo, emailSvc, &mockGenID{id: "token-2"}, configValida())

	resp, err := uc.Reenviar(context.Background(), verifyemail.ComandoReenviarVerificacion{UsuarioID: "user-1"})
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
	uc := verifyemail.NewVerificarCorreoCasoDeUso(
		&mockVerifRepo{err: errors.New("no encontrado")},
		&notificaciones.MockEmailServicio{}, &mockGenID{}, configValida(),
	)
	_, err := uc.Reenviar(context.Background(), verifyemail.ComandoReenviarVerificacion{UsuarioID: "user-x"})
	if !errors.Is(err, verificacion.ErrUsuarioNoEncontrado) {
		t.Errorf("esperaba ErrUsuarioNoEncontrado, got %v", err)
	}
}

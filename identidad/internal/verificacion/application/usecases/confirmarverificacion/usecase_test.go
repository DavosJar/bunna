package confirmarverificacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	uc_confirmar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/confirmarverificacion"
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

func configValida() uc_confirmar.ConfigVerificacion {
	return uc_confirmar.ConfigVerificacion{
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

func TestConfirmarVerificacionExitoso(t *testing.T) {
	usuario := usuarioNoVerificado()
	tokenPlano := "token-valido"
	usuario.PruebaVerificacion = verificacion.NuevaPruebaVerificacion(tokenPlano, time.Now().Add(24*time.Hour))
	repo := &mockVerifRepo{usuario: usuario}
	uc := uc_confirmar.NewConfirmarVerificacionCasoDeUso(repo, configValida())

	resp, err := uc.Ejecutar(context.Background(), &uc_confirmar.ComandoConfirmarVerificacion{Token: tokenPlano})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.Mensaje == "" {
		t.Error("Mensaje vacío")
	}
}

func TestConfirmarVerificacionTokenVacio(t *testing.T) {
	uc := uc_confirmar.NewConfirmarVerificacionCasoDeUso(&mockVerifRepo{}, configValida())
	_, err := uc.Ejecutar(context.Background(), &uc_confirmar.ComandoConfirmarVerificacion{Token: ""})
	if !errors.Is(err, verificacion.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestConfirmarVerificacionTokenInvalido(t *testing.T) {
	uc := uc_confirmar.NewConfirmarVerificacionCasoDeUso(
		&mockVerifRepo{err: errors.New("no encontrado")},
		configValida(),
	)
	_, err := uc.Ejecutar(context.Background(), &uc_confirmar.ComandoConfirmarVerificacion{Token: "bad-token"})
	if !errors.Is(err, verificacion.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestConfirmarVerificacionTokenExpirado(t *testing.T) {
	usuario := usuarioNoVerificado()
	usuario.PruebaVerificacion = verificacion.NuevaPruebaVerificacion("expirado", time.Now().Add(-48*time.Hour))
	uc := uc_confirmar.NewConfirmarVerificacionCasoDeUso(
		&mockVerifRepo{usuario: usuario},
		configValida(),
	)
	_, err := uc.Ejecutar(context.Background(), &uc_confirmar.ComandoConfirmarVerificacion{Token: "expirado"})
	if !errors.Is(err, verificacion.ErrEnlaceExpirado) {
		t.Errorf("esperaba ErrEnlaceExpirado, got %v", err)
	}
}

func TestConfirmarVerificacionFalloMarcar(t *testing.T) {
	usuario := usuarioNoVerificado()
	usuario.PruebaVerificacion = verificacion.NuevaPruebaVerificacion("token", time.Now().Add(24*time.Hour))
	uc := uc_confirmar.NewConfirmarVerificacionCasoDeUso(
		&mockVerifRepo{usuario: usuario, err: errors.New("fallo bd")},
		configValida(),
	)
	_, err := uc.Ejecutar(context.Background(), &uc_confirmar.ComandoConfirmarVerificacion{Token: "token"})
	if err == nil {
		t.Fatal("esperaba error al marcar verificado")
	}
}

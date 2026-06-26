package validartokenrecuperacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	uc_validar "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/validartokenrecuperacion"
)

type mockTokenRepo struct {
	token *recuperacion.TokenRecuperacion
	err   error
}

func (m *mockTokenRepo) Crear(ctx context.Context, token *recuperacion.TokenRecuperacion) error { return m.err }
func (m *mockTokenRepo) ObtenerPorHash(ctx context.Context, hash string) (*recuperacion.TokenRecuperacion, error) {
	return m.token, m.err
}
func (m *mockTokenRepo) Actualizar(ctx context.Context, token *recuperacion.TokenRecuperacion) error { return m.err }

func tokenValido() *recuperacion.TokenRecuperacion {
	return recuperacion.NuevoTokenRecuperacion("tok-1", "user-1", "token-plano", time.Now().Add(30*time.Minute))
}

func TestValidarTokenExitoso(t *testing.T) {
	uc := uc_validar.NewValidarTokenRecuperacionCasoDeUso(
		&mockTokenRepo{token: tokenValido()},
	)
	resp, err := uc.Ejecutar(context.Background(), &uc_validar.ComandoValidarTokenRecuperacion{Token: "token-plano"})
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
	uc := uc_validar.NewValidarTokenRecuperacionCasoDeUso(&mockTokenRepo{})
	_, err := uc.Ejecutar(context.Background(), &uc_validar.ComandoValidarTokenRecuperacion{Token: ""})
	if !errors.Is(err, recuperacion.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

func TestValidarTokenInvalido(t *testing.T) {
	uc := uc_validar.NewValidarTokenRecuperacionCasoDeUso(
		&mockTokenRepo{err: errors.New("no encontrado")},
	)
	_, err := uc.Ejecutar(context.Background(), &uc_validar.ComandoValidarTokenRecuperacion{Token: "bad"})
	if !errors.Is(err, recuperacion.ErrEnlaceInvalido) {
		t.Errorf("esperaba ErrEnlaceInvalido, got %v", err)
	}
}

package confirmarrecuperacion_test

import (
	"context"
	"errors"
	"testing"
	"time"

	recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	uc_confirmar "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/confirmarrecuperacion"
	uc_validar "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/validartokenrecuperacion"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
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

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
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
func (m *mockSesionRepo) ListarActivasPorUsuarioID(ctx context.Context, uid string, ahora time.Time) ([]*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) Listar(ctx context.Context, spec sesiones_domain.EspecificacionSesion, pag shareddomain.Paginacion) ([]*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) InvalidarTodasPorUsuarioID(ctx context.Context, uid string) error { return m.errInvalidar }
func (m *mockSesionRepo) Eliminar(ctx context.Context, id string) error                     { return nil }

type mockEncriptacion struct{}

func (m *mockEncriptacion) Hashear(password string) (string, error) { return "hash:" + password, nil }
func (m *mockEncriptacion) Verificar(password, hash string) bool    { return hash == "hash:"+password }

func tokenValido() *recuperacion.TokenRecuperacion {
	return recuperacion.NuevoTokenRecuperacion("tok-1", "user-1", "token-plano", time.Now().Add(30*time.Minute))
}

func TestConfirmarRecuperacionExitoso(t *testing.T) {
	tokenRepo := &mockTokenRepo{token: tokenValido()}
	usuarioRepo := &mockUsuarioRepo{usuario: &recuperacion.UsuarioRecuperacion{ID: "user-1"}}
	sesionRepo := &mockSesionRepo{}
	validarUC := uc_validar.NewValidarTokenRecuperacionCasoDeUso(tokenRepo)
	uc := uc_confirmar.NewConfirmarRecuperacionCasoDeUso(
		tokenRepo, usuarioRepo, sesionRepo, &mockEncriptacion{}, validarUC,
	)
	resp, err := uc.Ejecutar(context.Background(), &uc_confirmar.ComandoConfirmarRecuperacion{
		Token: "token-plano", NuevaPassword: "NuevoPass1!",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.Mensaje == "" {
		t.Error("Mensaje vacío")
	}
}

func TestConfirmarRecuperacionPasswordDebil(t *testing.T) {
	validarUC := uc_validar.NewValidarTokenRecuperacionCasoDeUso(&mockTokenRepo{})
	uc := uc_confirmar.NewConfirmarRecuperacionCasoDeUso(
		&mockTokenRepo{}, &mockUsuarioRepo{}, &mockSesionRepo{}, &mockEncriptacion{}, validarUC,
	)
	_, err := uc.Ejecutar(context.Background(), &uc_confirmar.ComandoConfirmarRecuperacion{
		Token: "tok", NuevaPassword: "short",
	})
	if !errors.Is(err, recuperacion.ErrPasswordDebil) {
		t.Errorf("esperaba ErrPasswordDebil, got %v", err)
	}
}

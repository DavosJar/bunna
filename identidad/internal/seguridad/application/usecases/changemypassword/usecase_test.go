package changemypassword_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/changemypassword"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockCredRepo struct {
	credenciales *seguridad.CredencialesUsuario
	errObtener   error
	errActualizar error
	actualizado   bool
}

func (m *mockCredRepo) Crear(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error) {
	return c, nil
}
func (m *mockCredRepo) Actualizar(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error) {
	if m.errActualizar != nil {
		return nil, m.errActualizar
	}
	m.actualizado = true
	return c, nil
}
func (m *mockCredRepo) ObtenerPorUsuarioID(ctx context.Context, id string) (*seguridad.CredencialesUsuario, error) {
	return m.credenciales, m.errObtener
}
func (m *mockCredRepo) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockCredRepo) Find(ctx context.Context, spec seguridad.EspecificacionCredenciales, pag shareddomain.Paginacion) ([]*seguridad.CredencialesUsuario, error) {
	return nil, nil
}

type mockEncriptacion struct {
	hash string
}

func (m *mockEncriptacion) Hashear(password string) (string, error) {
	return m.hash, nil
}
func (m *mockEncriptacion) Verificar(password, hash string) bool {
	return hash == "hash:"+password
}

func credencialesValidas() *seguridad.CredencialesUsuario {
	return seguridad.NuevaCredencialesUsuarioDesdeBD("user-id-1", "hash:PasswordActual1!", true, false, 0, time.Time{})
}

func TestCambiarMiContrasenaExitoso(t *testing.T) {
	credRepo := &mockCredRepo{credenciales: credencialesValidas()}
	uc := changemypassword.NewCambiarMiContrasenaCasoDeUso(credRepo, &mockEncriptacion{hash: "hash:NuevoPass1!"})

	resp, err := uc.Ejecutar(context.Background(), &changemypassword.ComandoCambiarMiContrasena{
		EjecutorID:     "user-id-1",
		PasswordActual: "PasswordActual1!",
		NuevaPassword:  "NuevoPass1!",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.EjecutorID != "user-id-1" {
		t.Errorf("EjecutorID incorrecto: %v", resp.EjecutorID)
	}
	if resp.ModificadoEn == "" {
		t.Error("ModificadoEn no debería estar vacío")
	}
	if !credRepo.actualizado {
		t.Error("credenciales no fueron actualizadas")
	}
}

func TestCambiarMiContrasenaNuevaPasswordVacia(t *testing.T) {
	uc := changemypassword.NewCambiarMiContrasenaCasoDeUso(&mockCredRepo{}, &mockEncriptacion{})
	_, err := uc.Ejecutar(context.Background(), &changemypassword.ComandoCambiarMiContrasena{
		EjecutorID: "user-id-1", PasswordActual: "old", NuevaPassword: "",
	})
	if err == nil {
		t.Fatal("esperaba error por password vacío")
	}
}

func TestCambiarMiContrasenaNuevaPasswordCorta(t *testing.T) {
	uc := changemypassword.NewCambiarMiContrasenaCasoDeUso(&mockCredRepo{}, &mockEncriptacion{})
	_, err := uc.Ejecutar(context.Background(), &changemypassword.ComandoCambiarMiContrasena{
		EjecutorID: "user-id-1", PasswordActual: "old", NuevaPassword: "Ab1!",
	})
	if err == nil {
		t.Fatal("esperaba error por password corta")
	}
}

func TestCambiarMiContrasenaPasswordActualIncorrecta(t *testing.T) {
	uc := changemypassword.NewCambiarMiContrasenaCasoDeUso(
		&mockCredRepo{credenciales: credencialesValidas()},
		&mockEncriptacion{hash: "hash:x"},
	)
	_, err := uc.Ejecutar(context.Background(), &changemypassword.ComandoCambiarMiContrasena{
		EjecutorID: "user-id-1", PasswordActual: "wrong", NuevaPassword: "NuevoPass1!",
	})
	if err == nil {
		t.Fatal("esperaba error por password actual incorrecta")
	}
}

func TestCambiarMiContrasenaUsuarioNoEncontrado(t *testing.T) {
	uc := changemypassword.NewCambiarMiContrasenaCasoDeUso(
		&mockCredRepo{errObtener: errors.New("no encontrado")},
		&mockEncriptacion{},
	)
	_, err := uc.Ejecutar(context.Background(), &changemypassword.ComandoCambiarMiContrasena{
		EjecutorID: "user-id-x", PasswordActual: "old", NuevaPassword: "NuevoPass1!",
	})
	if err == nil {
		t.Fatal("esperaba error por usuario no encontrado")
	}
}

func TestCambiarMiContrasenaFalloAlActualizar(t *testing.T) {
	uc := changemypassword.NewCambiarMiContrasenaCasoDeUso(
		&mockCredRepo{credenciales: credencialesValidas(), errActualizar: errors.New("fallo bd")},
		&mockEncriptacion{hash: "hash:NuevoPass1!"},
	)
	_, err := uc.Ejecutar(context.Background(), &changemypassword.ComandoCambiarMiContrasena{
		EjecutorID: "user-id-1", PasswordActual: "PasswordActual1!", NuevaPassword: "NuevoPass1!",
	})
	if err == nil {
		t.Fatal("esperaba error al actualizar")
	}
}

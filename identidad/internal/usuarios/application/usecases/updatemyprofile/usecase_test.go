package updatemyprofile_test

import (
	"context"
	"errors"
	"testing"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/updatemyprofile"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepoMyProfile struct {
	obtenerPorID func(ctx context.Context, id string) (*usuariodomain.Usuario, error)
	actualizar   func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepoMyProfile) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoMyProfile) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	if m.actualizar != nil {
		return m.actualizar(ctx, u)
	}
	return u, nil
}
func (m *mockUsuarioRepoMyProfile) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepoMyProfile) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	if m.obtenerPorID != nil {
		return m.obtenerPorID(ctx, id)
	}
	return nil, nil
}
func (m *mockUsuarioRepoMyProfile) Listar(ctx context.Context, _ usuariodomain.EspecificacionUsuario, _ shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	return nil, nil
}

func miembroPerfil(id, nombre, apellido string) *usuariodomain.Usuario {
	u, _ := usuariodomain.NuevoUsuario(id, "test@example.com", nombre, apellido, "555-0000")
	_ = u.Activar()
	return u
}

func TestModificarMiPerfilExitoso(t *testing.T) {
	user := miembroPerfil("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoMyProfile{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
	}
	uc := updatemyprofile.NewModificarMiPerfilCasoDeUso(repo)
	resp, err := uc.Ejecutar(context.Background(), &updatemyprofile.ComandoModificarMiPerfil{
		EjecutorID: "user-1", Nombre: "Juan Carlos", Apellido: "Perez Lopez",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.ID != "user-1" {
		t.Errorf("ID incorrecto: %s", resp.ID)
	}
	if resp.ModificadoEn == "" {
		t.Error("ModificadoEn no debe estar vacío")
	}
}

func TestModificarMiPerfilNoEncontrado(t *testing.T) {
	repo := &mockUsuarioRepoMyProfile{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return nil, errors.New("not found")
		},
	}
	uc := updatemyprofile.NewModificarMiPerfilCasoDeUso(repo)
	_, err := uc.Ejecutar(context.Background(), &updatemyprofile.ComandoModificarMiPerfil{
		EjecutorID: "no-existe",
	})
	if err == nil {
		t.Fatal("esperaba error de usuario no encontrado")
	}
}

func TestModificarMiPerfilActualizarError(t *testing.T) {
	user := miembroPerfil("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoMyProfile{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
		actualizar: func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
			return nil, errors.New("db error")
		},
	}
	uc := updatemyprofile.NewModificarMiPerfilCasoDeUso(repo)
	_, err := uc.Ejecutar(context.Background(), &updatemyprofile.ComandoModificarMiPerfil{
		EjecutorID: "user-1", Nombre: "Juan Carlos", Apellido: "Perez Lopez",
	})
	if err == nil {
		t.Fatal("esperaba error al actualizar")
	}
}

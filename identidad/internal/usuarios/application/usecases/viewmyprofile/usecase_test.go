package viewmyprofile_test

import (
	"context"
	"errors"
	"testing"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/viewmyprofile"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepoView struct {
	obtenerPorID func(ctx context.Context, id string) (*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepoView) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoView) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoView) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepoView) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	if m.obtenerPorID != nil {
		return m.obtenerPorID(ctx, id)
	}
	return nil, nil
}
func (m *mockUsuarioRepoView) Listar(ctx context.Context, _ usuariodomain.EspecificacionUsuario, _ shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	return nil, nil
}

func perfilUsuario(id, nombre, apellido string) *usuariodomain.Usuario {
	u, _ := usuariodomain.NuevoUsuario(id, "test@example.com", nombre, apellido, "555-0000")
	return u
}

func TestVerMiPerfilExitoso(t *testing.T) {
	user := perfilUsuario("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoView{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
	}
	uc := viewmyprofile.NewVerMiPerfilCasoDeUso(repo)
	resp, err := uc.Ejecutar(context.Background(), &viewmyprofile.ComandoVerMiPerfil{EjecutorID: "user-1"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.ID != "user-1" {
		t.Errorf("ID incorrecto: %s", resp.ID)
	}
	if resp.Correo != "test@example.com" {
		t.Errorf("Correo incorrecto: %s", resp.Correo)
	}
	if resp.Nombre != "Juan" {
		t.Errorf("Nombre incorrecto: %s", resp.Nombre)
	}
	if resp.Apellido != "Perez" {
		t.Errorf("Apellido incorrecto: %s", resp.Apellido)
	}
	if resp.Telefono != "555-0000" {
		t.Errorf("Telefono incorrecto: %s", resp.Telefono)
	}
	if resp.Estado != string(usuariodomain.NO_VERIFICADO) {
		t.Errorf("Estado incorrecto: %s", resp.Estado)
	}
	if resp.CreadoEn == "" {
		t.Error("CreadoEn no debe estar vacío")
	}
}

func TestVerMiPerfilNoEncontrado(t *testing.T) {
	repo := &mockUsuarioRepoView{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return nil, errors.New("not found")
		},
	}
	uc := viewmyprofile.NewVerMiPerfilCasoDeUso(repo)
	_, err := uc.Ejecutar(context.Background(), &viewmyprofile.ComandoVerMiPerfil{EjecutorID: "no-existe"})
	if err == nil {
		t.Fatal("esperaba error de usuario no encontrado")
	}
}

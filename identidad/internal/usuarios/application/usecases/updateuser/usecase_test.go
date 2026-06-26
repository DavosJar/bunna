package updateuser_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/updateuser"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepoUpdate struct {
	obtenerPorID func(ctx context.Context, id string) (*usuariodomain.Usuario, error)
	actualizar   func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepoUpdate) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoUpdate) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	if m.actualizar != nil {
		return m.actualizar(ctx, u)
	}
	return u, nil
}
func (m *mockUsuarioRepoUpdate) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepoUpdate) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	if m.obtenerPorID != nil {
		return m.obtenerPorID(ctx, id)
	}
	return nil, nil
}
func (m *mockUsuarioRepoUpdate) Listar(ctx context.Context, _ usuariodomain.EspecificacionUsuario, _ shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	return nil, nil
}

type mockAuthSvcUpdate struct {
	ok  bool
	err error
}

func (m *mockAuthSvcUpdate) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func usuarioModificable(id, nombre, apellido string) *usuariodomain.Usuario {
	u, _ := usuariodomain.NuevoUsuario(id, "test@example.com", nombre, apellido, "555-0000")
	_ = u.Activar()
	return u
}

func TestModificarUsuarioExitoso(t *testing.T) {
	user := usuarioModificable("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoUpdate{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
	}
	uc := updateuser.NewModificarUsuarioCasoDeUso(repo, &mockAuthSvcUpdate{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &updateuser.ComandoModificarUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
		Nombre: "Juan Modificado", Apellido: "Perez Modificado",
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

func TestModificarUsuarioPermisoDenegado(t *testing.T) {
	uc := updateuser.NewModificarUsuarioCasoDeUso(&mockUsuarioRepoUpdate{}, &mockAuthSvcUpdate{ok: false})
	_, err := uc.Ejecutar(context.Background(), &updateuser.ComandoModificarUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestModificarUsuarioAuthError(t *testing.T) {
	uc := updateuser.NewModificarUsuarioCasoDeUso(&mockUsuarioRepoUpdate{}, &mockAuthSvcUpdate{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &updateuser.ComandoModificarUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestModificarUsuarioNoEncontrado(t *testing.T) {
	repo := &mockUsuarioRepoUpdate{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return nil, errors.New("no encontrado")
		},
	}
	uc := updateuser.NewModificarUsuarioCasoDeUso(repo, &mockAuthSvcUpdate{ok: true})
	_, err := uc.Ejecutar(context.Background(), &updateuser.ComandoModificarUsuario{
		UsuarioID: "no-existe", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de usuario no encontrado")
	}
}

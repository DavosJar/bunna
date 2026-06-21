package deleteuser_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/deleteuser"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepoDelete struct {
	obtenerPorID func(ctx context.Context, id string) (*usuariodomain.Usuario, error)
	actualizar   func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepoDelete) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoDelete) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	if m.actualizar != nil {
		return m.actualizar(ctx, u)
	}
	return u, nil
}
func (m *mockUsuarioRepoDelete) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepoDelete) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	if m.obtenerPorID != nil {
		return m.obtenerPorID(ctx, id)
	}
	return nil, nil
}
func (m *mockUsuarioRepoDelete) Listar(ctx context.Context, _ usuariodomain.EspecificacionUsuario, _ shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	return nil, nil
}

type mockAuthSvcDelete struct {
	ok  bool
	err error
}

func (m *mockAuthSvcDelete) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func usuarioActivo(id, nombre, apellido string) *usuariodomain.Usuario {
	u, _ := usuariodomain.NuevoUsuario(id, "test@example.com", nombre, apellido, "555-0000")
	_ = u.Activar()
	return u
}

func usuarioEnPendiente(id string) *usuariodomain.Usuario {
	u, _ := usuariodomain.NuevoUsuario(id, "test@example.com", "Test", "User", "")
	_ = u.Activar()
	_ = u.CambiarEstado(usuariodomain.PENDIENTE_DE_ELIMINACION)
	return u
}

func TestDarDeBajaUsuarioExitoso(t *testing.T) {
	user := usuarioActivo("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoDelete{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
	}
	uc := deleteuser.NewDarDeBajaUsuarioCasoDeUso(repo, &mockAuthSvcDelete{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &deleteuser.ComandoDarDeBajaUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.UsuarioID != "user-1" {
		t.Errorf("UsuarioID incorrecto: %s", resp.UsuarioID)
	}
	if resp.Estado != string(usuariodomain.PENDIENTE_DE_ELIMINACION) {
		t.Errorf("Estado incorrecto: %s", resp.Estado)
	}
	if resp.BajaEn == "" {
		t.Error("BajaEn no debe estar vacío")
	}
}

func TestDarDeBajaUsuarioPermisoDenegado(t *testing.T) {
	uc := deleteuser.NewDarDeBajaUsuarioCasoDeUso(&mockUsuarioRepoDelete{}, &mockAuthSvcDelete{ok: false})
	_, err := uc.Ejecutar(context.Background(), &deleteuser.ComandoDarDeBajaUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestDarDeBajaUsuarioAuthError(t *testing.T) {
	uc := deleteuser.NewDarDeBajaUsuarioCasoDeUso(&mockUsuarioRepoDelete{}, &mockAuthSvcDelete{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &deleteuser.ComandoDarDeBajaUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestDarDeBajaUsuarioAutoBaja(t *testing.T) {
	uc := deleteuser.NewDarDeBajaUsuarioCasoDeUso(&mockUsuarioRepoDelete{}, &mockAuthSvcDelete{ok: true})
	_, err := uc.Ejecutar(context.Background(), &deleteuser.ComandoDarDeBajaUsuario{
		UsuarioID: "admin-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil || err.Error() != "no puedes darte de baja a ti mismo" {
		t.Errorf("esperaba error de auto-baja, got %v", err)
	}
}

func TestDarDeBajaUsuarioNoEncontrado(t *testing.T) {
	repo := &mockUsuarioRepoDelete{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return nil, errors.New("not found")
		},
	}
	uc := deleteuser.NewDarDeBajaUsuarioCasoDeUso(repo, &mockAuthSvcDelete{ok: true})
	_, err := uc.Ejecutar(context.Background(), &deleteuser.ComandoDarDeBajaUsuario{
		UsuarioID: "no-existe", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de usuario no encontrado")
	}
}

func TestDarDeBajaUsuarioYaEnPendiente(t *testing.T) {
	user := usuarioEnPendiente("user-1")
	repo := &mockUsuarioRepoDelete{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
	}
	uc := deleteuser.NewDarDeBajaUsuarioCasoDeUso(repo, &mockAuthSvcDelete{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &deleteuser.ComandoDarDeBajaUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("no debería fallar para usuario ya pendiente: %v", err)
	}
	if resp.Estado != string(usuariodomain.PENDIENTE_DE_ELIMINACION) {
		t.Errorf("Estado esperado PENDIENTE_DE_ELIMINACION, got %s", resp.Estado)
	}
}

func TestDarDeBajaUsuarioActualizarError(t *testing.T) {
	user := usuarioActivo("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoDelete{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
		actualizar: func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
			return nil, errors.New("db error")
		},
	}
	uc := deleteuser.NewDarDeBajaUsuarioCasoDeUso(repo, &mockAuthSvcDelete{ok: true})
	_, err := uc.Ejecutar(context.Background(), &deleteuser.ComandoDarDeBajaUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error al persistir")
	}
}

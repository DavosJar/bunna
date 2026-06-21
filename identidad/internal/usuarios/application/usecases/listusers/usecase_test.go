package listusers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/listusers"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepoList struct {
	listar func(ctx context.Context, spec usuariodomain.EspecificacionUsuario, pag shareddomain.Paginacion) ([]*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepoList) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoList) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoList) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepoList) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	return nil, nil
}
func (m *mockUsuarioRepoList) Listar(ctx context.Context, spec usuariodomain.EspecificacionUsuario, pag shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	if m.listar != nil {
		return m.listar(ctx, spec, pag)
	}
	return nil, nil
}

type mockAuthSvcList struct {
	ok  bool
	err error
}

func (m *mockAuthSvcList) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func usuarioListable(id, nombre, apellido string) *usuariodomain.Usuario {
	u, _ := usuariodomain.NuevoUsuario(id, "test@example.com", nombre, apellido, "")
	return u
}

func TestListarUsuariosExitoso(t *testing.T) {
	users := []*usuariodomain.Usuario{
		usuarioListable("u1", "Juan", "Perez"),
		usuarioListable("u2", "Maria", "Gomez"),
	}
	repo := &mockUsuarioRepoList{
		listar: func(ctx context.Context, spec usuariodomain.EspecificacionUsuario, pag shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
			return users, nil
		},
	}
	uc := listusers.NewListarUsuariosCasoDeUso(repo, &mockAuthSvcList{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &listusers.ComandoListarUsuarios{
		TenantID: "tenant-1", EjecutorID: "admin-1",
		Paginacion: shareddomain.Paginacion{Pagina: 1, TamanoPagina: 10},
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(resp.Usuarios) != 2 {
		t.Errorf("esperaba 2 usuarios, got %d", len(resp.Usuarios))
	}
	if resp.Total != 2 {
		t.Errorf("Total incorrecto: %d", resp.Total)
	}
	if resp.Pagina != 1 {
		t.Errorf("Pagina incorrecta: %d", resp.Pagina)
	}
}

func TestListarUsuariosVacio(t *testing.T) {
	repo := &mockUsuarioRepoList{
		listar: func(ctx context.Context, spec usuariodomain.EspecificacionUsuario, pag shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
			return []*usuariodomain.Usuario{}, nil
		},
	}
	uc := listusers.NewListarUsuariosCasoDeUso(repo, &mockAuthSvcList{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &listusers.ComandoListarUsuarios{
		TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(resp.Usuarios) != 0 {
		t.Errorf("esperaba 0 usuarios, got %d", len(resp.Usuarios))
	}
}

func TestListarUsuariosPermisoDenegado(t *testing.T) {
	uc := listusers.NewListarUsuariosCasoDeUso(&mockUsuarioRepoList{}, &mockAuthSvcList{ok: false})
	_, err := uc.Ejecutar(context.Background(), &listusers.ComandoListarUsuarios{
		TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestListarUsuariosAuthError(t *testing.T) {
	uc := listusers.NewListarUsuariosCasoDeUso(&mockUsuarioRepoList{}, &mockAuthSvcList{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &listusers.ComandoListarUsuarios{
		TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestListarUsuariosRepoError(t *testing.T) {
	repo := &mockUsuarioRepoList{
		listar: func(ctx context.Context, spec usuariodomain.EspecificacionUsuario, pag shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
			return nil, errors.New("db error")
		},
	}
	uc := listusers.NewListarUsuariosCasoDeUso(repo, &mockAuthSvcList{ok: true})
	_, err := uc.Ejecutar(context.Background(), &listusers.ComandoListarUsuarios{
		TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

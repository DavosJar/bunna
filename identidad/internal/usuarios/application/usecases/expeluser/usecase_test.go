package expeluser_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	sesiones "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/expeluser"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepoExpel struct {
	obtenerPorID func(ctx context.Context, id string) (*usuariodomain.Usuario, error)
	actualizar   func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepoExpel) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoExpel) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	if m.actualizar != nil {
		return m.actualizar(ctx, u)
	}
	return u, nil
}
func (m *mockUsuarioRepoExpel) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepoExpel) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	if m.obtenerPorID != nil {
		return m.obtenerPorID(ctx, id)
	}
	return nil, nil
}
func (m *mockUsuarioRepoExpel) Listar(ctx context.Context, _ usuariodomain.EspecificacionUsuario, _ shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	return nil, nil
}

type mockSesionRepoExpel struct {
	errInvalidar error
}

func (m *mockSesionRepoExpel) Crear(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) { return s, nil }
func (m *mockSesionRepoExpel) Actualizar(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) { return s, nil }
func (m *mockSesionRepoExpel) ObtenerPorID(ctx context.Context, id string) (*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoExpel) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoExpel) ListarActivasPorUsuarioID(ctx context.Context, uid string, ahora time.Time) ([]*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoExpel) Listar(ctx context.Context, spec sesiones.EspecificacionSesion, pag shareddomain.Paginacion) ([]*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoExpel) InvalidarTodasPorUsuarioID(ctx context.Context, uid string) error { return m.errInvalidar }
func (m *mockSesionRepoExpel) Eliminar(ctx context.Context, id string) error { return nil }

type mockAuthSvcExpel struct {
	ok  bool
	err error
}

func (m *mockAuthSvcExpel) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func usuarioActivoExpel(id, nombre, apellido string) *usuariodomain.Usuario {
	u, _ := usuariodomain.NuevoUsuario(id, "test@example.com", nombre, apellido, "")
	_ = u.Activar()
	return u
}

func TestExpulsarUsuarioExitoso(t *testing.T) {
	user := usuarioActivoExpel("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoExpel{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
	}
	uc := expeluser.NewExpulsarUsuarioCasoDeUso(repo, &mockSesionRepoExpel{}, &mockAuthSvcExpel{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &expeluser.ComandoExpulsarUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.UsuarioID != "user-1" {
		t.Errorf("UsuarioID incorrecto: %s", resp.UsuarioID)
	}
	if resp.Estado != string(usuariodomain.BLOQUEADO) {
		t.Errorf("Estado esperado BLOQUEADO, got %s", resp.Estado)
	}
	if resp.ExpulsadoEn == "" {
		t.Error("ExpulsadoEn no debe estar vacío")
	}
}

func TestExpulsarUsuarioPermisoDenegado(t *testing.T) {
	uc := expeluser.NewExpulsarUsuarioCasoDeUso(&mockUsuarioRepoExpel{}, &mockSesionRepoExpel{}, &mockAuthSvcExpel{ok: false})
	_, err := uc.Ejecutar(context.Background(), &expeluser.ComandoExpulsarUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestExpulsarUsuarioAuthError(t *testing.T) {
	uc := expeluser.NewExpulsarUsuarioCasoDeUso(&mockUsuarioRepoExpel{}, &mockSesionRepoExpel{}, &mockAuthSvcExpel{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &expeluser.ComandoExpulsarUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestExpulsarUsuarioAutoExpulsion(t *testing.T) {
	uc := expeluser.NewExpulsarUsuarioCasoDeUso(&mockUsuarioRepoExpel{}, &mockSesionRepoExpel{}, &mockAuthSvcExpel{ok: true})
	_, err := uc.Ejecutar(context.Background(), &expeluser.ComandoExpulsarUsuario{
		UsuarioID: "admin-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil || err.Error() != "no puedes expulsarte a ti mismo" {
		t.Errorf("esperaba error de auto-expulsión, got %v", err)
	}
}

func TestExpulsarUsuarioNoEncontrado(t *testing.T) {
	repo := &mockUsuarioRepoExpel{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return nil, errors.New("not found")
		},
	}
	uc := expeluser.NewExpulsarUsuarioCasoDeUso(repo, &mockSesionRepoExpel{}, &mockAuthSvcExpel{ok: true})
	_, err := uc.Ejecutar(context.Background(), &expeluser.ComandoExpulsarUsuario{
		UsuarioID: "no-existe", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de usuario no encontrado")
	}
}

func TestExpulsarUsuarioSesionRepoError(t *testing.T) {
	user := usuarioActivoExpel("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoExpel{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
	}
	sessionRepo := &mockSesionRepoExpel{errInvalidar: errors.New("db error")}
	uc := expeluser.NewExpulsarUsuarioCasoDeUso(repo, sessionRepo, &mockAuthSvcExpel{ok: true})
	_, err := uc.Ejecutar(context.Background(), &expeluser.ComandoExpulsarUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error al invalidar sesiones")
	}
}

func TestExpulsarUsuarioActualizarError(t *testing.T) {
	user := usuarioActivoExpel("user-1", "Juan", "Perez")
	repo := &mockUsuarioRepoExpel{
		obtenerPorID: func(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
			return user, nil
		},
		actualizar: func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
			return nil, errors.New("db error")
		},
	}
	uc := expeluser.NewExpulsarUsuarioCasoDeUso(repo, &mockSesionRepoExpel{}, &mockAuthSvcExpel{ok: true})
	_, err := uc.Ejecutar(context.Background(), &expeluser.ComandoExpulsarUsuario{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error al persistir")
	}
}

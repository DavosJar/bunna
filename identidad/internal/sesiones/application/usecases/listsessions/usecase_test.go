package listsessions_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/listsessions"
	sesiones "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockSesionRepoList struct {
	listarActivas func(ctx context.Context, uid string, ahora time.Time) ([]*sesiones.Sesion, error)
}

func (m *mockSesionRepoList) Crear(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) { return s, nil }
func (m *mockSesionRepoList) Actualizar(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) { return s, nil }
func (m *mockSesionRepoList) ObtenerPorID(ctx context.Context, id string) (*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoList) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoList) ListarActivasPorUsuarioID(ctx context.Context, uid string, ahora time.Time) ([]*sesiones.Sesion, error) {
	if m.listarActivas != nil {
		return m.listarActivas(ctx, uid, ahora)
	}
	return nil, nil
}
func (m *mockSesionRepoList) Listar(ctx context.Context, spec sesiones.EspecificacionSesion, pag shareddomain.Paginacion) ([]*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoList) InvalidarTodasPorUsuarioID(ctx context.Context, uid string) error { return nil }
func (m *mockSesionRepoList) Eliminar(ctx context.Context, id string) error { return nil }

type mockAuthSvcListSes struct {
	ok  bool
	err error
}

func (m *mockAuthSvcListSes) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func sesionListable(id, usuarioID, ip string) *sesiones.Sesion {
	now := time.Now()
	s, _ := sesiones.NuevaSesion(id, usuarioID, "ah", "rh", ip,
		now, now.Add(15*time.Minute), now.Add(24*time.Hour))
	return s
}

func TestListarSesionesExitoso(t *testing.T) {
	sessions := []*sesiones.Sesion{
		sesionListable("s1", "user-1", "10.0.0.1"),
		sesionListable("s2", "user-1", "10.0.0.2"),
	}
	repo := &mockSesionRepoList{
		listarActivas: func(ctx context.Context, uid string, ahora time.Time) ([]*sesiones.Sesion, error) {
			return sessions, nil
		},
	}
	uc := listsessions.NewListarSesionesCasoDeUso(repo, &mockAuthSvcListSes{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &listsessions.ComandoListarSesiones{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(resp.Sesiones) != 2 {
		t.Errorf("esperaba 2 sesiones, got %d", len(resp.Sesiones))
	}
	if resp.Total != 2 {
		t.Errorf("Total incorrecto: %d", resp.Total)
	}
}

func TestListarSesionesVacio(t *testing.T) {
	repo := &mockSesionRepoList{
		listarActivas: func(ctx context.Context, uid string, ahora time.Time) ([]*sesiones.Sesion, error) {
			return []*sesiones.Sesion{}, nil
		},
	}
	uc := listsessions.NewListarSesionesCasoDeUso(repo, &mockAuthSvcListSes{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &listsessions.ComandoListarSesiones{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(resp.Sesiones) != 0 {
		t.Errorf("esperaba 0 sesiones, got %d", len(resp.Sesiones))
	}
}

func TestListarSesionesPermisoDenegado(t *testing.T) {
	uc := listsessions.NewListarSesionesCasoDeUso(&mockSesionRepoList{}, &mockAuthSvcListSes{ok: false})
	_, err := uc.Ejecutar(context.Background(), &listsessions.ComandoListarSesiones{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestListarSesionesAuthError(t *testing.T) {
	uc := listsessions.NewListarSesionesCasoDeUso(&mockSesionRepoList{}, &mockAuthSvcListSes{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &listsessions.ComandoListarSesiones{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestListarSesionesRepoError(t *testing.T) {
	repo := &mockSesionRepoList{
		listarActivas: func(ctx context.Context, uid string, ahora time.Time) ([]*sesiones.Sesion, error) {
			return nil, errors.New("db error")
		},
	}
	uc := listsessions.NewListarSesionesCasoDeUso(repo, &mockAuthSvcListSes{ok: true})
	_, err := uc.Ejecutar(context.Background(), &listsessions.ComandoListarSesiones{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

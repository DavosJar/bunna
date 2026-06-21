package terminatesession_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/terminatesession"
	sesiones "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockSesionRepoTerminate struct {
	obtenerPorID  func(ctx context.Context, id string) (*sesiones.Sesion, error)
	actualizar    func(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error)
}

func (m *mockSesionRepoTerminate) Crear(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) { return s, nil }
func (m *mockSesionRepoTerminate) Actualizar(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) {
	if m.actualizar != nil {
		return m.actualizar(ctx, s)
	}
	return s, nil
}
func (m *mockSesionRepoTerminate) ObtenerPorID(ctx context.Context, id string) (*sesiones.Sesion, error) {
	if m.obtenerPorID != nil {
		return m.obtenerPorID(ctx, id)
	}
	return nil, nil
}
func (m *mockSesionRepoTerminate) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoTerminate) ListarActivasPorUsuarioID(ctx context.Context, uid string, ahora time.Time) ([]*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoTerminate) Listar(ctx context.Context, spec sesiones.EspecificacionSesion, pag shareddomain.Paginacion) ([]*sesiones.Sesion, error) { return nil, nil }
func (m *mockSesionRepoTerminate) InvalidarTodasPorUsuarioID(ctx context.Context, uid string) error { return nil }
func (m *mockSesionRepoTerminate) Eliminar(ctx context.Context, id string) error { return nil }

type mockAuthSvcTerminate struct {
	ok  bool
	err error
}

func (m *mockAuthSvcTerminate) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func sesionTerminable(id, usuarioID string) *sesiones.Sesion {
	now := time.Now()
	s, _ := sesiones.NuevaSesion(id, usuarioID, "ah", "rh", "10.0.0.1",
		now, now.Add(15*time.Minute), now.Add(24*time.Hour))
	return s
}

func TestForzarCierreSesionExitoso(t *testing.T) {
	session := sesionTerminable("s1", "user-1")
	repo := &mockSesionRepoTerminate{
		obtenerPorID: func(ctx context.Context, id string) (*sesiones.Sesion, error) {
			return session, nil
		},
	}
	uc := terminatesession.NewForzarCierreSesionCasoDeUso(repo, &mockAuthSvcTerminate{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &terminatesession.ComandoForzarCierreSesion{
		SesionID: "s1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.SesionID != "s1" {
		t.Errorf("SesionID incorrecto: %s", resp.SesionID)
	}
	if resp.Estado != string(sesiones.EstadoRevocada) {
		t.Errorf("Estado esperado REVOCADA, got %s", resp.Estado)
	}
	if resp.RevocadoEn == "" {
		t.Error("RevocadoEn no debe estar vacío")
	}
}

func TestForzarCierreSesionPermisoDenegado(t *testing.T) {
	uc := terminatesession.NewForzarCierreSesionCasoDeUso(&mockSesionRepoTerminate{}, &mockAuthSvcTerminate{ok: false})
	_, err := uc.Ejecutar(context.Background(), &terminatesession.ComandoForzarCierreSesion{
		SesionID: "s1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestForzarCierreSesionAuthError(t *testing.T) {
	uc := terminatesession.NewForzarCierreSesionCasoDeUso(&mockSesionRepoTerminate{}, &mockAuthSvcTerminate{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &terminatesession.ComandoForzarCierreSesion{
		SesionID: "s1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestForzarCierreSesionNoEncontrada(t *testing.T) {
	repo := &mockSesionRepoTerminate{
		obtenerPorID: func(ctx context.Context, id string) (*sesiones.Sesion, error) {
			return nil, errors.New("not found")
		},
	}
	uc := terminatesession.NewForzarCierreSesionCasoDeUso(repo, &mockAuthSvcTerminate{ok: true})
	_, err := uc.Ejecutar(context.Background(), &terminatesession.ComandoForzarCierreSesion{
		SesionID: "no-existe", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de sesión no encontrada")
	}
}

func TestForzarCierreSesionActualizarError(t *testing.T) {
	session := sesionTerminable("s1", "user-1")
	repo := &mockSesionRepoTerminate{
		obtenerPorID: func(ctx context.Context, id string) (*sesiones.Sesion, error) {
			return session, nil
		},
		actualizar: func(ctx context.Context, s *sesiones.Sesion) (*sesiones.Sesion, error) {
			return nil, errors.New("db error")
		},
	}
	uc := terminatesession.NewForzarCierreSesionCasoDeUso(repo, &mockAuthSvcTerminate{ok: true})
	_, err := uc.Ejecutar(context.Background(), &terminatesession.ComandoForzarCierreSesion{
		SesionID: "s1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error al persistir")
	}
}

package viewcredentials_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/viewcredentials"
	seguridadDomain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockCredRepoView struct {
	obtenerPorUsuarioID func(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error)
}

func (m *mockCredRepoView) Crear(ctx context.Context, c *seguridadDomain.CredencialesUsuario) (*seguridadDomain.CredencialesUsuario, error) {
	return c, nil
}
func (m *mockCredRepoView) Actualizar(ctx context.Context, c *seguridadDomain.CredencialesUsuario) (*seguridadDomain.CredencialesUsuario, error) {
	return c, nil
}
func (m *mockCredRepoView) ObtenerPorUsuarioID(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error) {
	if m.obtenerPorUsuarioID != nil {
		return m.obtenerPorUsuarioID(ctx, uid)
	}
	return nil, nil
}
func (m *mockCredRepoView) Eliminar(ctx context.Context, uid string) error { return nil }
func (m *mockCredRepoView) Find(ctx context.Context, _ seguridadDomain.EspecificacionCredenciales, _ shareddomain.Paginacion) ([]*seguridadDomain.CredencialesUsuario, error) {
	return nil, nil
}

type mockAuthSvcView struct {
	ok  bool
	err error
}

func (m *mockAuthSvcView) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func TestConsultarCredencialesExitoso(t *testing.T) {
	creds := seguridadDomain.NuevaCredencialesUsuarioDesdeBD("user-1", "hash", true, true, 3, time.Time{})
	repo := &mockCredRepoView{
		obtenerPorUsuarioID: func(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error) {
			return creds, nil
		},
	}
	uc := viewcredentials.NewConsultarCredencialesCasoDeUso(repo, &mockAuthSvcView{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &viewcredentials.ComandoConsultarCredenciales{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.UsuarioID != "user-1" {
		t.Errorf("UsuarioID incorrecto: %s", resp.UsuarioID)
	}
	if !resp.Activo {
		t.Error("Activo debe ser true")
	}
	if !resp.CorreoVerificado {
		t.Error("CorreoVerificado debe ser true")
	}
	if resp.IntentosFallidos != 3 {
		t.Errorf("IntentosFallidos incorrecto: %d", resp.IntentosFallidos)
	}
	if resp.BloqueadoHasta != "" {
		t.Errorf("BloqueadoHasta debe ser vacío, got %s", resp.BloqueadoHasta)
	}
}

func TestConsultarCredencialesConBloqueo(t *testing.T) {
	bloqueadoHasta := time.Now().Add(30 * time.Minute)
	creds := seguridadDomain.NuevaCredencialesUsuarioDesdeBD("user-1", "hash", false, false, 5, bloqueadoHasta)
	repo := &mockCredRepoView{
		obtenerPorUsuarioID: func(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error) {
			return creds, nil
		},
	}
	uc := viewcredentials.NewConsultarCredencialesCasoDeUso(repo, &mockAuthSvcView{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &viewcredentials.ComandoConsultarCredenciales{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.BloqueadoHasta == "" {
		t.Error("BloqueadoHasta no debe estar vacío para cuenta bloqueada")
	}
	if resp.Activo {
		t.Error("Activo debe ser false")
	}
	if resp.CorreoVerificado {
		t.Error("CorreoVerificado debe ser false")
	}
}

func TestConsultarCredencialesPermisoDenegado(t *testing.T) {
	uc := viewcredentials.NewConsultarCredencialesCasoDeUso(&mockCredRepoView{}, &mockAuthSvcView{ok: false})
	_, err := uc.Ejecutar(context.Background(), &viewcredentials.ComandoConsultarCredenciales{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestConsultarCredencialesAuthError(t *testing.T) {
	uc := viewcredentials.NewConsultarCredencialesCasoDeUso(&mockCredRepoView{}, &mockAuthSvcView{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &viewcredentials.ComandoConsultarCredenciales{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestConsultarCredencialesNoEncontradas(t *testing.T) {
	repo := &mockCredRepoView{
		obtenerPorUsuarioID: func(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error) {
			return nil, errors.New("not found")
		},
	}
	uc := viewcredentials.NewConsultarCredencialesCasoDeUso(repo, &mockAuthSvcView{ok: true})
	_, err := uc.Ejecutar(context.Background(), &viewcredentials.ComandoConsultarCredenciales{
		UsuarioID: "no-existe", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de credenciales no encontradas")
	}
}

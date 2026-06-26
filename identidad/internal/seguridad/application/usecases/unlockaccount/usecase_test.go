package unlockaccount_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/unlockaccount"
	seguridadDomain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockCredRepoUnlock struct {
	obtenerPorUsuarioID func(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error)
}

func (m *mockCredRepoUnlock) Crear(ctx context.Context, c *seguridadDomain.CredencialesUsuario) (*seguridadDomain.CredencialesUsuario, error) {
	return c, nil
}
func (m *mockCredRepoUnlock) Actualizar(ctx context.Context, c *seguridadDomain.CredencialesUsuario) (*seguridadDomain.CredencialesUsuario, error) {
	return c, nil
}
func (m *mockCredRepoUnlock) ObtenerPorUsuarioID(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error) {
	if m.obtenerPorUsuarioID != nil {
		return m.obtenerPorUsuarioID(ctx, uid)
	}
	return nil, nil
}
func (m *mockCredRepoUnlock) Eliminar(ctx context.Context, uid string) error { return nil }
func (m *mockCredRepoUnlock) Find(ctx context.Context, _ seguridadDomain.EspecificacionCredenciales, _ shareddomain.Paginacion) ([]*seguridadDomain.CredencialesUsuario, error) {
	return nil, nil
}

type mockAuthSvcUnlock struct {
	ok  bool
	err error
}

func (m *mockAuthSvcUnlock) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func TestDesbloquearCuentaExitoso(t *testing.T) {
	creds := seguridadDomain.NuevaCredencialesUsuario("user-1", "hash")
	repo := &mockCredRepoUnlock{
		obtenerPorUsuarioID: func(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error) {
			return creds, nil
		},
	}
	uc := unlockaccount.NewDesbloquearCuentaCasoDeUso(repo, &mockAuthSvcUnlock{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &unlockaccount.ComandoDesbloquearCuenta{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.UsuarioID != "user-1" {
		t.Errorf("UsuarioID incorrecto: %s", resp.UsuarioID)
	}
	if resp.DesbloqueadoEn == "" {
		t.Error("DesbloqueadoEn no debe estar vacío")
	}
}

func TestDesbloquearCuentaPermisoDenegado(t *testing.T) {
	uc := unlockaccount.NewDesbloquearCuentaCasoDeUso(&mockCredRepoUnlock{}, &mockAuthSvcUnlock{ok: false})
	_, err := uc.Ejecutar(context.Background(), &unlockaccount.ComandoDesbloquearCuenta{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestDesbloquearCuentaAuthError(t *testing.T) {
	uc := unlockaccount.NewDesbloquearCuentaCasoDeUso(&mockCredRepoUnlock{}, &mockAuthSvcUnlock{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &unlockaccount.ComandoDesbloquearCuenta{
		UsuarioID: "user-1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestDesbloquearCuentaNoEncontrada(t *testing.T) {
	repo := &mockCredRepoUnlock{
		obtenerPorUsuarioID: func(ctx context.Context, uid string) (*seguridadDomain.CredencialesUsuario, error) {
			return nil, errors.New("no encontrado")
		},
	}
	uc := unlockaccount.NewDesbloquearCuentaCasoDeUso(repo, &mockAuthSvcUnlock{ok: true})
	_, err := uc.Ejecutar(context.Background(), &unlockaccount.ComandoDesbloquearCuenta{
		UsuarioID: "no-existe", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de credenciales no encontradas")
	}
}

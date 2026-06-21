package listblockedips_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/listblockedips"
	seguridadDomain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockIntentoRepoList struct {
	listar func(ctx context.Context, spec seguridadDomain.EspecificacionIntentoIP, pag shareddomain.Paginacion) ([]*seguridadDomain.IntentoPorIP, error)
}

func (m *mockIntentoRepoList) ObtenerPorIP(ctx context.Context, ip string) (*seguridadDomain.IntentoPorIP, error) {
	return nil, nil
}
func (m *mockIntentoRepoList) Crear(ctx context.Context, intento *seguridadDomain.IntentoPorIP) (*seguridadDomain.IntentoPorIP, error) {
	return intento, nil
}
func (m *mockIntentoRepoList) Actualizar(ctx context.Context, intento *seguridadDomain.IntentoPorIP) (*seguridadDomain.IntentoPorIP, error) {
	return intento, nil
}
func (m *mockIntentoRepoList) Listar(ctx context.Context, spec seguridadDomain.EspecificacionIntentoIP, pag shareddomain.Paginacion) ([]*seguridadDomain.IntentoPorIP, error) {
	if m.listar != nil {
		return m.listar(ctx, spec, pag)
	}
	return nil, nil
}
func (m *mockIntentoRepoList) EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error { return nil }

type mockAuthSvcListIPs struct {
	ok  bool
	err error
}

func (m *mockAuthSvcListIPs) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func TestListarIPsBloqueadasExitoso(t *testing.T) {
	uc := seguridad.NewListarIPsBloqueadasCasoDeUso(
		&mockIntentoRepoList{}, &mockAuthSvcListIPs{ok: true},
	)
	resp, err := uc.Ejecutar(context.Background(), &seguridad.ComandoListarIPsBloqueadas{
		TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp == nil {
		t.Fatal("respuesta no debe ser nil")
	}
	if len(resp.IPs) != 0 {
		t.Errorf("esperaba 0 IPs, got %d", len(resp.IPs))
	}
}

func TestListarIPsBloqueadasPermisoDenegado(t *testing.T) {
	uc := seguridad.NewListarIPsBloqueadasCasoDeUso(
		&mockIntentoRepoList{}, &mockAuthSvcListIPs{ok: false},
	)
	_, err := uc.Ejecutar(context.Background(), &seguridad.ComandoListarIPsBloqueadas{
		TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestListarIPsBloqueadasAuthError(t *testing.T) {
	uc := seguridad.NewListarIPsBloqueadasCasoDeUso(
		&mockIntentoRepoList{}, &mockAuthSvcListIPs{err: errors.New("fallo")},
	)
	_, err := uc.Ejecutar(context.Background(), &seguridad.ComandoListarIPsBloqueadas{
		TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

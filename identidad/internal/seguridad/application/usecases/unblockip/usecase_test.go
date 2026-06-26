package unblockip_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/unblockip"
	seguridadDomain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type mockIntentoRepoUnblock struct {
	obtenerPorIP  func(ctx context.Context, ip string) (*seguridadDomain.IntentoPorIP, error)
	actualizar    func(ctx context.Context, intento *seguridadDomain.IntentoPorIP) (*seguridadDomain.IntentoPorIP, error)
}

func (m *mockIntentoRepoUnblock) ObtenerPorIP(ctx context.Context, ip string) (*seguridadDomain.IntentoPorIP, error) {
	if m.obtenerPorIP != nil {
		return m.obtenerPorIP(ctx, ip)
	}
	return nil, nil
}
func (m *mockIntentoRepoUnblock) Crear(ctx context.Context, intento *seguridadDomain.IntentoPorIP) (*seguridadDomain.IntentoPorIP, error) {
	return intento, nil
}
func (m *mockIntentoRepoUnblock) Actualizar(ctx context.Context, intento *seguridadDomain.IntentoPorIP) (*seguridadDomain.IntentoPorIP, error) {
	if m.actualizar != nil {
		return m.actualizar(ctx, intento)
	}
	return intento, nil
}
func (m *mockIntentoRepoUnblock) Listar(ctx context.Context, spec seguridadDomain.EspecificacionIntentoIP, pag shareddomain.Paginacion) ([]*seguridadDomain.IntentoPorIP, error) {
	return nil, nil
}
func (m *mockIntentoRepoUnblock) EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error { return nil }

type mockAuthSvcUnblockIP struct {
	ok  bool
	err error
}

func (m *mockAuthSvcUnblockIP) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

func intentoBloqueado(ip string) *seguridadDomain.IntentoPorIP {
	now := time.Now()
	intento := seguridadDomain.NuevoIntentoPorIP("id-1", ip, now.Add(-10*time.Minute))
	intento.Bloquear(now.Add(30 * time.Minute))
	return intento
}

func TestDesbloquearIPExitoso(t *testing.T) {
	intento := intentoBloqueado("10.0.0.1")
	repo := &mockIntentoRepoUnblock{
		obtenerPorIP: func(ctx context.Context, ip string) (*seguridadDomain.IntentoPorIP, error) {
			return intento, nil
		},
	}
	uc := unblockip.NewDesbloquearIPCasoDeUso(repo, &mockAuthSvcUnblockIP{ok: true})
	resp, err := uc.Ejecutar(context.Background(), &unblockip.ComandoDesbloquearIP{
		IP: "10.0.0.1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.IP != "10.0.0.1" {
		t.Errorf("IP incorrecta: %s", resp.IP)
	}
	if resp.DesbloqueadoEn == "" {
		t.Error("DesbloqueadoEn no debe estar vacío")
	}
	if intento.EstaBloqueada(time.Now()) {
		t.Error("la IP debió ser desbloqueada")
	}
}

func TestDesbloquearIPPermisoDenegado(t *testing.T) {
	uc := unblockip.NewDesbloquearIPCasoDeUso(&mockIntentoRepoUnblock{}, &mockAuthSvcUnblockIP{ok: false})
	_, err := uc.Ejecutar(context.Background(), &unblockip.ComandoDesbloquearIP{
		IP: "10.0.0.1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestDesbloquearIPAuthError(t *testing.T) {
	uc := unblockip.NewDesbloquearIPCasoDeUso(&mockIntentoRepoUnblock{}, &mockAuthSvcUnblockIP{err: errors.New("fallo")})
	_, err := uc.Ejecutar(context.Background(), &unblockip.ComandoDesbloquearIP{
		IP: "10.0.0.1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de auth")
	}
}

func TestDesbloquearIPNoEncontrada(t *testing.T) {
	repo := &mockIntentoRepoUnblock{
		obtenerPorIP: func(ctx context.Context, ip string) (*seguridadDomain.IntentoPorIP, error) {
			return nil, errors.New("no encontrado")
		},
	}
	uc := unblockip.NewDesbloquearIPCasoDeUso(repo, &mockAuthSvcUnblockIP{ok: true})
	_, err := uc.Ejecutar(context.Background(), &unblockip.ComandoDesbloquearIP{
		IP: "no-existe", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de IP no encontrada")
	}
}

func TestDesbloquearIPNoBloqueada(t *testing.T) {
	now := time.Now()
	intento := seguridadDomain.NuevoIntentoPorIP("id-1", "10.0.0.1", now.Add(-10*time.Minute))
	repo := &mockIntentoRepoUnblock{
		obtenerPorIP: func(ctx context.Context, ip string) (*seguridadDomain.IntentoPorIP, error) {
			return intento, nil
		},
	}
	uc := unblockip.NewDesbloquearIPCasoDeUso(repo, &mockAuthSvcUnblockIP{ok: true})
	_, err := uc.Ejecutar(context.Background(), &unblockip.ComandoDesbloquearIP{
		IP: "10.0.0.1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error de IP no bloqueada")
	}
}

func TestDesbloquearIPActualizarError(t *testing.T) {
	intento := intentoBloqueado("10.0.0.1")
	repo := &mockIntentoRepoUnblock{
		obtenerPorIP: func(ctx context.Context, ip string) (*seguridadDomain.IntentoPorIP, error) {
			return intento, nil
		},
		actualizar: func(ctx context.Context, i *seguridadDomain.IntentoPorIP) (*seguridadDomain.IntentoPorIP, error) {
			return nil, errors.New("error de BD")
		},
	}
	uc := unblockip.NewDesbloquearIPCasoDeUso(repo, &mockAuthSvcUnblockIP{ok: true})
	_, err := uc.Ejecutar(context.Background(), &unblockip.ComandoDesbloquearIP{
		IP: "10.0.0.1", TenantID: "tenant-1", EjecutorID: "admin-1",
	})
	if err == nil {
		t.Fatal("esperaba error al persistir")
	}
}

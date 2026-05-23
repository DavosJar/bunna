
package rate_limiter_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/rate_limiter"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/stretchr/testify/mock"
)

// ── Mocks ────────────────────────────────────────────────────────────────────

type mockIntentoIPRepo struct {
	mock.Mock
	intento     *seguridad_domain.IntentoPorIP
	errObtener  error
	creado      *seguridad_domain.IntentoPorIP
	actualizado *seguridad_domain.IntentoPorIP
}

func (m *mockIntentoIPRepo) ObtenerPorIP(ctx context.Context, ip string) (*seguridad_domain.IntentoPorIP, error) {
	return m.intento, m.errObtener
}
func (m *mockIntentoIPRepo) Crear(ctx context.Context, i *seguridad_domain.IntentoPorIP) (*seguridad_domain.IntentoPorIP, error) {
	m.creado = i
	return i, nil
}
func (m *mockIntentoIPRepo) Actualizar(ctx context.Context, i *seguridad_domain.IntentoPorIP) (*seguridad_domain.IntentoPorIP, error) {
	m.actualizado = i
	return i, nil
}
func (m *mockIntentoIPRepo) Listar(ctx context.Context, spec seguridad_domain.EspecificacionIntentoIP, pag shared_domain.Paginacion) ([]*seguridad_domain.IntentoPorIP, error) {
	args := m.Called(ctx, spec, pag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*seguridad_domain.IntentoPorIP), args.Error(1)
}
func (m *mockIntentoIPRepo) EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error {
	return nil
}

type mockGeneradorID struct{}

func (m *mockGeneradorID) NextID(ctx context.Context) (string, error) {
	return "id-generado", nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func configDefault() rate_limiter.ConfigRateLimit {
	return rate_limiter.ConfigRateLimit{
		MaxRequests: 10,
		Ventana:     1 * time.Minute,
	}
}

// ── Tests ────────────────────────────────────────────────────────────────────

// Escenario 7: Rate limit dentro del límite → permitir
func TestRateLimit_DentroDelLimite(t *testing.T) {
	ahora := time.Now()
	// 5 requests previos, límite 10
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 5, ahora, time.Time{},
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := rate_limiter.NuevoServicioRateLimit(repo, &mockGeneradorID{}, configDefault())

	err := svc.Verificar(context.Background(), "10.0.0.1")
	if err != nil {
		t.Errorf("esperaba request permitido dentro del límite, got %v", err)
	}
	if repo.actualizado == nil {
		t.Error("esperaba que se actualizara el contador")
	}
	if repo.actualizado.Contador() != 6 {
		t.Errorf("esperaba contador=6, got %d", repo.actualizado.Contador())
	}
}

// Escenario 8: Rate limit excedido → error
func TestRateLimit_LimiteExcedido(t *testing.T) {
	ahora := time.Now()
	// Ya en el límite exacto (10)
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 10, ahora, time.Time{},
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := rate_limiter.NuevoServicioRateLimit(repo, &mockGeneradorID{}, configDefault())

	err := svc.Verificar(context.Background(), "10.0.0.1")
	if !errors.Is(err, rate_limiter.ErrRateLimitExcedido) {
		t.Errorf("esperaba ErrRateLimitExcedido, got %v", err)
	}
}

// Escenario 9: Ventana deslizante — ventana expirada reinicia contador
func TestRateLimit_VentanaDeslizante(t *testing.T) {
	// ventanaInicio hace 2 minutos, ventana es 1 minuto → expirada
	hace2min := time.Now().Add(-2 * time.Minute)
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 10, hace2min, time.Time{},
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := rate_limiter.NuevoServicioRateLimit(repo, &mockGeneradorID{}, configDefault())

	err := svc.Verificar(context.Background(), "10.0.0.1")
	if err != nil {
		t.Errorf("esperaba request permitido con ventana expirada, got %v", err)
	}
	// debe haber creado un nuevo registro con contador reiniciado
	if repo.creado == nil {
		t.Error("esperaba nuevo registro al expirar la ventana")
	}
	if repo.creado.Contador() != 1 {
		t.Errorf("esperaba contador=1 tras reinicio de ventana, got %d", repo.creado.Contador())
	}
}

// Escenario 10: Reset después de ventana → permitir
func TestRateLimit_ResetDespuesDeVentana(t *testing.T) {
	// IP sin registro previo → primera request
	repo := &mockIntentoIPRepo{errObtener: errors.New("no encontrado")}
	svc := rate_limiter.NuevoServicioRateLimit(repo, &mockGeneradorID{}, configDefault())

	err := svc.Verificar(context.Background(), "10.0.0.1")
	if err != nil {
		t.Errorf("esperaba primera request permitida, got %v", err)
	}
	if repo.creado == nil {
		t.Error("esperaba que se creara registro para primera request")
	}
	if repo.creado.Contador() != 1 {
		t.Errorf("esperaba contador=1 para primera request, got %d", repo.creado.Contador())
	}
}

// Validación: IP vacía
func TestRateLimit_IPVacia(t *testing.T) {
	svc := rate_limiter.NuevoServicioRateLimit(&mockIntentoIPRepo{}, &mockGeneradorID{}, configDefault())
	err := svc.Verificar(context.Background(), "")
	if !errors.Is(err, rate_limiter.ErrIPRequerida) {
		t.Errorf("esperaba ErrIPRequerida, got %v", err)
	}
}

// Límite exactamente en el umbral con un request más
func TestRateLimit_11RequestsExcedeLimite(t *testing.T) {
	ahora := time.Now()
	// 11 requests en la ventana actual
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 11, ahora, time.Time{},
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := rate_limiter.NuevoServicioRateLimit(repo, &mockGeneradorID{}, configDefault())

	err := svc.Verificar(context.Background(), "10.0.0.1")
	if !errors.Is(err, rate_limiter.ErrRateLimitExcedido) {
		t.Errorf("esperaba ErrRateLimitExcedido, got %v", err)
	}
}
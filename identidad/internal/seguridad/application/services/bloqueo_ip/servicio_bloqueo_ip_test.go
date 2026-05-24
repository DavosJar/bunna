package bloqueo_ip_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
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

func configDefault() bloqueo_ip.ConfigBloqueoIP {
	return bloqueo_ip.ConfigBloqueoIP{
		MaxIntentos: 20,
		Ventana:     15 * time.Minute,
		Duracion:    30 * time.Minute,
	}
}

// ── Tests ────────────────────────────────────────────────────────────────────

// Escenario 1: IP no bloqueada
func TestBloqueoIP_IPNoRegistrada_Permitida(t *testing.T) {
	repo := &mockIntentoIPRepo{errObtener: errors.New("no encontrada")}
	svc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorID{}, configDefault())

	err := svc.Verificar(context.Background(), "192.168.1.1")
	if err != nil {
		t.Errorf("esperaba IP permitida, got %v", err)
	}
}

// Escenario 2: IP bloqueada por umbral
func TestBloqueoIP_IPBloqueada(t *testing.T) {
	ahora := time.Now()
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "192.168.1.1", 20, ahora.Add(-5*time.Minute), ahora.Add(30*time.Minute),
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorID{}, configDefault())

	err := svc.Verificar(context.Background(), "192.168.1.1")
	if !errors.Is(err, bloqueo_ip.ErrIPBloqueada) {
		t.Errorf("esperaba ErrIPBloqueada, got %v", err)
	}
}

// Escenario 5: Bloqueo de IP expirado → permitir
func TestBloqueoIP_BloqueoExpirado_Permitida(t *testing.T) {
	ahora := time.Now()
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "192.168.1.1", 20, ahora.Add(-1*time.Hour), ahora.Add(-1*time.Minute),
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorID{}, configDefault())

	err := svc.Verificar(context.Background(), "192.168.1.1")
	if err != nil {
		t.Errorf("esperaba IP permitida con bloqueo expirado, got %v", err)
	}
}

// Escenario: Registrar intento fallido crea registro nuevo
func TestBloqueoIP_RegistrarIntento_NuevoRegistro(t *testing.T) {
	repo := &mockIntentoIPRepo{errObtener: errors.New("no encontrada")}
	svc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorID{}, configDefault())

	err := svc.RegistrarIntentoFallido(context.Background(), "10.0.0.1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.creado == nil {
		t.Error("esperaba que se creara un nuevo registro de intento")
	}
}

// Escenario: Registrar intento fallido incrementa contador existente
func TestBloqueoIP_RegistrarIntento_IncrementaContador(t *testing.T) {
	ahora := time.Now()
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 5, ahora, time.Time{},
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorID{}, configDefault())

	err := svc.RegistrarIntentoFallido(context.Background(), "10.0.0.1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.actualizado == nil {
		t.Error("esperaba que se actualizara el registro")
	}
	if repo.actualizado.Contador() != 6 {
		t.Errorf("esperaba contador=6, got %d", repo.actualizado.Contador())
	}
}

// Escenario 2: Al llegar al umbral bloquea la IP
func TestBloqueoIP_AlcanzarUmbral_BloquearIP(t *testing.T) {
	ahora := time.Now()
	// 19 intentos previos, el siguiente llega al umbral (20)
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 19, ahora, time.Time{},
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorID{}, configDefault())

	err := svc.RegistrarIntentoFallido(context.Background(), "10.0.0.1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !repo.actualizado.EstaBloqueada(time.Now()) {
		t.Error("esperaba que la IP quedara bloqueada al alcanzar el umbral")
	}
}

// Escenario: Ventana expirada reinicia contador
func TestBloqueoIP_VentanaExpirada_ReiniciaContador(t *testing.T) {
	// ventanaInicio hace 20 min, ventana es 15 min → expirada
	hace20min := time.Now().Add(-20 * time.Minute)
	intento := seguridad_domain.NuevoIntentoPorIPDesdeBD(
		"id-1", "10.0.0.1", 15, hace20min, time.Time{},
	)
	repo := &mockIntentoIPRepo{intento: intento}
	svc := bloqueo_ip.NuevoServicioBloqueoIP(repo, &mockGeneradorID{}, configDefault())

	err := svc.RegistrarIntentoFallido(context.Background(), "10.0.0.1")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	// debe haber creado un nuevo registro (ventana reiniciada)
	if repo.creado == nil {
		t.Error("esperaba nuevo registro al expirar la ventana")
	}
	if repo.creado.Contador() != 1 {
		t.Errorf("esperaba contador=1 tras reinicio, got %d", repo.creado.Contador())
	}
}

// Validación: IP vacía
func TestBloqueoIP_IPVacia(t *testing.T) {
	svc := bloqueo_ip.NuevoServicioBloqueoIP(&mockIntentoIPRepo{}, &mockGeneradorID{}, configDefault())
	err := svc.Verificar(context.Background(), "")
	if !errors.Is(err, bloqueo_ip.ErrIPRequerida) {
		t.Errorf("esperaba ErrIPRequerida, got %v", err)
	}
}

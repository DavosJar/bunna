// Package rate_limiter implementa el caso de uso de limitación de requests por IP.
package rate_limiter

import (
	"context"
	"errors"
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

var (
	// ErrRateLimitExcedido se retorna cuando la IP excede el límite de requests por ventana.
	ErrRateLimitExcedido = errors.New("demasiados intentos, intente más tarde")

	// ErrIPRequerida se retorna cuando la IP está vacía.
	ErrIPRequerida = errors.New("la IP es requerida")
)

// ConfigRateLimit contiene los parámetros configurables del rate limiter.
type ConfigRateLimit struct {
	// MaxRequests es el número máximo de requests permitidos por ventana.
	MaxRequests int

	// Ventana es el periodo de tiempo de la ventana deslizante.
	Ventana time.Duration
}

// ServicioRateLimit implementa la limitación de requests por IP con ventana deslizante.
type ServicioRateLimit struct {
	repo        seguridad_domain.IntentoIPRepositorio
	generadorID shared_domain.GeneradorID
	config      ConfigRateLimit
}

// NuevoServicioRateLimit crea una nueva instancia de ServicioRateLimit.
func NuevoServicioRateLimit(
	repo seguridad_domain.IntentoIPRepositorio,
	generadorID shared_domain.GeneradorID,
	config ConfigRateLimit,
) *ServicioRateLimit {
	return &ServicioRateLimit{
		repo:        repo,
		generadorID: generadorID,
		config:      config,
	}
}

// Verificar comprueba si una IP ha excedido el límite de requests en la ventana actual.
// Es preventivo: debe llamarse ANTES de procesar el request.
// Retorna ErrRateLimitExcedido si se superó el límite, nil si está permitido.
func (s *ServicioRateLimit) Verificar(ctx context.Context, ip string) error {
	if ip == "" {
		return ErrIPRequerida
	}

	ahora := time.Now()
	intento, err := s.repo.ObtenerPorIP(ctx, ip)
	if err != nil {
		// No existe registro → primera request, registrar y permitir
		id, err := s.generadorID.NextID(ctx)
		if err != nil {
			return err
		}
		nuevo := seguridad_domain.NuevoIntentoPorIP(id, ip, ahora)
		_, err = s.repo.Crear(ctx, nuevo)
		return err
	}

	// Si la ventana expiró → reiniciar y permitir
	if intento.VentanaExpirada(ahora, s.config.Ventana) {
		id, err := s.generadorID.NextID(ctx)
		if err != nil {
			return err
		}
		nuevo := seguridad_domain.NuevoIntentoPorIP(id, ip, ahora)
		_, err = s.repo.Crear(ctx, nuevo)
		return err
	}

	// Verificar límite en ventana activa
	if intento.Contador() >= s.config.MaxRequests {
		return ErrRateLimitExcedido
	}

	// Incrementar contador y permitir
	intento.IncrementarContador()
	_, err = s.repo.Actualizar(ctx, intento)
	return err
}

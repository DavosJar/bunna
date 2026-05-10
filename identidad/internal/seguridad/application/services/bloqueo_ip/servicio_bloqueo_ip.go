// Package bloqueo_ip implementa el caso de uso de bloqueo de IPs por intentos fallidos.
package bloqueo_ip

import (
	"context"
	"errors"
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

var (
	// ErrIPBloqueada se retorna cuando la IP está bloqueada por exceso de intentos.
	ErrIPBloqueada = errors.New("IP bloqueada temporalmente por exceso de intentos")

	// ErrIPRequerida se retorna cuando la IP está vacía.
	ErrIPRequerida = errors.New("la IP es requerida")
)

// ConfigBloqueoIP contiene los parámetros configurables del servicio de bloqueo por IP.
type ConfigBloqueoIP struct {
	// MaxIntentos es el número de intentos fallidos antes de bloquear la IP.
	MaxIntentos int

	// Ventana es el periodo de tiempo en el que se cuentan los intentos.
	Ventana time.Duration

	// Duracion es el tiempo que la IP permanece bloqueada.
	Duracion time.Duration
}

// ServicioBloqueoIP implementa la verificación y registro de intentos fallidos por IP.
type ServicioBloqueoIP struct {
	repo        seguridad_domain.IntentoIPRepositorio
	generadorID shared_domain.GeneradorID
	config      ConfigBloqueoIP
}

// NuevoServicioBloqueoIP crea una nueva instancia de ServicioBloqueoIP.
func NuevoServicioBloqueoIP(
	repo seguridad_domain.IntentoIPRepositorio,
	generadorID shared_domain.GeneradorID,
	config ConfigBloqueoIP,
) *ServicioBloqueoIP {
	return &ServicioBloqueoIP{
		repo:        repo,
		generadorID: generadorID,
		config:      config,
	}
}

// Verificar comprueba si una IP está bloqueada.
// Debe llamarse ANTES de procesar cualquier intento de login.
// Retorna ErrIPBloqueada si la IP está bloqueada, nil si está permitida.
func (s *ServicioBloqueoIP) Verificar(ctx context.Context, ip string) error {
	if ip == "" {
		return ErrIPRequerida
	}

	ahora := time.Now()
	intento, err := s.repo.ObtenerPorIP(ctx, ip)
	if err != nil {
		// No existe registro → IP limpia
		return nil
	}

	// Verificar si el bloqueo está activo
	if intento.EstaBloqueada(ahora) {
		return ErrIPBloqueada
	}

	return nil
}

// RegistrarIntentoFallido incrementa el contador de intentos fallidos para una IP.
// Si se supera el umbral, bloquea la IP por el tiempo configurado.
// Debe llamarse DESPUÉS de un intento de login fallido.
func (s *ServicioBloqueoIP) RegistrarIntentoFallido(ctx context.Context, ip string) error {
	if ip == "" {
		return ErrIPRequerida
	}

	ahora := time.Now()
	intento, err := s.repo.ObtenerPorIP(ctx, ip)
	if err != nil {
		// No existe registro → crear uno nuevo
		id, err := s.generadorID.NextID(ctx)
		if err != nil {
			return err
		}
		nuevo := seguridad_domain.NuevoIntentoPorIP(id, ip, ahora)
		_, err = s.repo.Crear(ctx, nuevo)
		return err
	}

	// Si la ventana expiró, reiniciar contador
	if intento.VentanaExpirada(ahora, s.config.Ventana) {
		id, err := s.generadorID.NextID(ctx)
		if err != nil {
			return err
		}
		nuevo := seguridad_domain.NuevoIntentoPorIP(id, ip, ahora)
		_, err = s.repo.Crear(ctx, nuevo)
		return err
	}

	// Incrementar contador en ventana activa
	intento.IncrementarContador()

	// Bloquear si se supera el umbral
	if intento.Contador() >= s.config.MaxIntentos {
		intento.Bloquear(ahora.Add(s.config.Duracion))
	}

	_, err = s.repo.Actualizar(ctx, intento)
	return err
}
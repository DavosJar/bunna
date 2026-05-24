package refresh

import (
	"context"
	"errors"
	"time"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

var (
	// ErrRefreshTokenRequerido se retorna cuando el refresh token está vacío.
	ErrRefreshTokenRequerido = errors.New("el refresh token es requerido")

	// ErrTokenInvalido se retorna cuando el token no puede ser validado
	// (expirado, mal formado, firma inválida) o cuando su hash no existe en BD
	// (posible robo de token).
	ErrTokenInvalido = errors.New("token inválido o expirado")

	// ErrSesionNoValida se retorna cuando la sesión asociada al token
	// está en estado REVOCADA o EXPIRADA.
	ErrSesionNoValida = errors.New("sesión no válida")

	// ErrRefreshTokenExpirado se retorna cuando el refresh token de la sesión
	// ya superó su fecha de expiración.
	ErrRefreshTokenExpirado = errors.New("refresh token expirado, inicie sesión nuevamente")

	// ErrLimiteRefrescosAlcanzado se retorna cuando la sesión superó el máximo
	// de refrescos configurado.
	ErrLimiteRefrescosAlcanzado = errors.New("límite de refrescos alcanzado, inicie sesión nuevamente")

	// ErrSesionAbsolutaExpirada se retorna cuando la sesión superó el tiempo
	// máximo absoluto de vida configurado.
	ErrSesionAbsolutaExpirada = errors.New("sesión expirada, inicie sesión nuevamente")

	// ErrErrorGenerandoTokens se retorna cuando el TokenServicio falla
	// al generar el nuevo par de tokens.
	ErrErrorGenerandoTokens = errors.New("error al generar tokens")
)

// ConfigRefresh contiene los parámetros configurables del servicio de refresh.
type ConfigRefresh struct {
	// MaxRefrescos es el número máximo de refrescos permitidos por sesión.
	// 0 significa sin límite.
	MaxRefrescos int

	// TimeoutAbsoluto es el tiempo máximo de vida de una sesión desde su creación,
	// independientemente de la actividad. Por ejemplo: 168h (7 días).
	TimeoutAbsoluto time.Duration
}

// ServicioRefresh implementa el caso de uso de renovación de sesión.
// Aplica rotación obligatoria de refresh token y detección de robo.
type ServicioRefresh struct {
	uow    sesiones_domain.UnitOfWork
	config ConfigRefresh
}

// NuevoServicioRefresh crea una nueva instancia de ServicioRefresh.
func NuevoServicioRefresh(uow sesiones_domain.UnitOfWork, config ConfigRefresh) *ServicioRefresh {
	return &ServicioRefresh{uow: uow, config: config}
}

// Ejecutar renueva la sesión asociada al refresh token recibido.
//
// Flujo:
//  1. Valida que el token no esté vacío.
//  2. Valida la firma y expiración del JWT (fuera de transacción).
//  3. Dentro de una transacción:
//     a. Busca la sesión por hash del refresh token.
//     b. Si no existe → detección de robo: invalida todas las sesiones del usuario.
//     c. Si existe pero está inactiva → error sesión no válida.
//     d. Si existe y activa → rota tokens, actualiza sesión.
func (s *ServicioRefresh) Ejecutar(ctx context.Context, cmd ComandoRefresh) (*RespuestaRefresh, error) {
	// 1. Validar comando
	if cmd.RefreshToken == "" {
		return nil, ErrRefreshTokenRequerido
	}

	// 2. Validar JWT fuera de transacción
	claims, err := s.uow.TokenServicio().ValidarRefreshToken(cmd.RefreshToken)
	if err != nil {
		return nil, ErrTokenInvalido
	}

	// Calcular hash del token recibido para buscar en BD
	refreshTokenHash := s.uow.TokenServicio().HashearToken(cmd.RefreshToken)

	var respuesta *RespuestaRefresh

	err = s.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

		// 3a. Buscar sesión por hash del refresh token
		sesion, err := tx.SesionRepositorio().ObtenerPorRefreshTokenHash(ctx, refreshTokenHash)
		if err != nil {
			// 3b. Hash no existe en BD → token fue rotado → posible robo
			// Invalidar todas las sesiones activas del usuario
			_ = tx.SesionRepositorio().InvalidarTodasPorUsuarioID(ctx, claims.UsuarioID)
			return ErrTokenInvalido
		}

		// 3c. Sesión existe pero está inactiva
		if sesion.Estado() == sesiones_domain.EstadoRevocada || sesion.Estado() == sesiones_domain.EstadoExpirada {
			return ErrSesionNoValida
		}

		// 3d. Sesión activa: verificar timeout absoluto
		if s.config.TimeoutAbsoluto > 0 {
			if ahora.After(sesion.FechaCreacion().Add(s.config.TimeoutAbsoluto)) {
				_ = sesion.MarcarExpirada()
				_, _ = tx.SesionRepositorio().Actualizar(ctx, sesion)
				return ErrSesionAbsolutaExpirada
			}
		}

		// Verificar que el refresh token de la sesión no haya expirado
		if !sesion.RefreshTokenValido(ahora) {
			_ = sesion.MarcarExpirada()
			_, _ = tx.SesionRepositorio().Actualizar(ctx, sesion)
			return ErrRefreshTokenExpirado
		}

		// Verificar límite de refrescos
		if s.config.MaxRefrescos > 0 && sesion.ContadorRefrescos() >= s.config.MaxRefrescos {
			return ErrLimiteRefrescosAlcanzado
		}

		// Generar nuevo par de tokens
		nuevoAccessToken, expiracionAccess, err := tx.TokenServicio().GenerarAccessToken(claims.UsuarioID, sesion.ID(), nil)
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		nuevoRefreshToken, expiracionRefresh, err := tx.TokenServicio().GenerarRefreshToken(claims.UsuarioID, sesion.ID())
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		// Hashear nuevos tokens
		nuevoAccessHash := tx.TokenServicio().HashearToken(nuevoAccessToken)
		nuevoRefreshHash := tx.TokenServicio().HashearToken(nuevoRefreshToken)

		// Rotar tokens en la entidad
		if err := sesion.RotarTokens(nuevoAccessHash, nuevoRefreshHash, expiracionAccess, expiracionRefresh, ahora); err != nil {
			return err
		}

		// Persistir sesión actualizada
		if _, err := tx.SesionRepositorio().Actualizar(ctx, sesion); err != nil {
			return err
		}

		respuesta = &RespuestaRefresh{
			AccessToken:       nuevoAccessToken,
			RefreshToken:      nuevoRefreshToken,
			ExpiracionAccess:  expiracionAccess,
			ExpiracionRefresh: expiracionRefresh,
			SesionID:          sesion.ID(),
			UsuarioID:         claims.UsuarioID,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

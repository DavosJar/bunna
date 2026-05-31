package login

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/rate_limiter"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

var (
	ErrCredencialesInvalidas = errors.New("credenciales inválidas")
	ErrCuentaBloqueada       = errors.New("cuenta temporalmente bloqueada")
	ErrCuentaInactiva        = errors.New("cuenta inactiva")
	ErrEmailRequerido        = errors.New("el email es requerido")
	ErrEmailInvalido         = errors.New("el email no tiene un formato válido")
	ErrPasswordRequerido     = errors.New("la contraseña es requerida")
	ErrErrorGenerandoTokens  = errors.New("error al generar tokens")
)

// ServicioLogin implementa el caso de uso de inicio de sesión.
// Integra verificación de rate limiting y bloqueo por IP antes de procesar credenciales.
type ServicioLogin struct {
	uow         sesiones_domain.UnitOfWork
	bloqueoIP   *bloqueo_ip.ServicioBloqueoIP
	rateLimiter *rate_limiter.ServicioRateLimit
}

// NuevoServicioLogin crea una nueva instancia de ServicioLogin.
// bloqueoIP y rateLimiter son opcionales — si son nil se omite la verificación.
func NuevoServicioLogin(
	uow sesiones_domain.UnitOfWork,
	bloqueoIP *bloqueo_ip.ServicioBloqueoIP,
	rateLimiter *rate_limiter.ServicioRateLimit,
) *ServicioLogin {
	return &ServicioLogin{
		uow:         uow,
		bloqueoIP:   bloqueoIP,
		rateLimiter: rateLimiter,
	}
}

// Ejecutar procesa el intento de login con el siguiente flujo:
//  1. Validar comando (email y password).
//  2. Verificar rate limiting por IP (preventivo).
//  3. Verificar bloqueo por IP.
//  4. Transacción: buscar usuario → credenciales → verificar password → crear sesión.
//  5. Si falla por password incorrecto: registrar intento fallido por IP.
func (s *ServicioLogin) Ejecutar(ctx context.Context, cmd ComandoLogin) (*RespuestaLogin, error) {
	// 1. Validar comando fuera de transacción
	if err := validarComando(cmd); err != nil {
		return nil, err
	}

	// 2. Rate limiting preventivo (antes de cualquier procesamiento)
	if s.rateLimiter != nil && cmd.IPOrigen != "" {
		if err := s.rateLimiter.Verificar(ctx, cmd.IPOrigen); err != nil {
			return nil, err
		}
	}

	// 3. Verificar bloqueo por IP
	if s.bloqueoIP != nil && cmd.IPOrigen != "" {
		if err := s.bloqueoIP.Verificar(ctx, cmd.IPOrigen); err != nil {
			return nil, err
		}
	}

	var respuesta *RespuestaLogin
	var passwordIncorrecto bool

	err := s.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

		// Resolver email → usuarioID
		usuarios, err := tx.UsuarioRepositorio().Listar(ctx,
			usuario_domain.EspecificacionUsuario{
				ListaLiltros: []domain.CriterioFiltro{
					{Campo: "correo", Operador: "=", Valor: cmd.Email},
				},
			},
			domain.Paginacion{Pagina: 1, TamanoPagina: 1},
		)
		if err != nil || len(usuarios) == 0 {
			return ErrCredencialesInvalidas
		}
		usuarioID := usuarios[0].ID()

		credenciales, err := tx.CredencialesRepositorio().ObtenerPorUsuarioID(ctx, usuarioID)
		if err != nil {
			return ErrCredencialesInvalidas
		}

		if credenciales.EstaBloqueado(ahora) {
			return ErrCuentaBloqueada
		}

		if !credenciales.Activo() {
			return ErrCuentaInactiva
		}

		if !tx.EncriptacionServicio().Verificar(cmd.Password, credenciales.PasswordHash()) {
			credenciales.MarcarIntentoFallido(ahora)
			if _, err := tx.CredencialesRepositorio().Actualizar(ctx, credenciales); err != nil {
				return err
			}
			passwordIncorrecto = true
			return ErrCredencialesInvalidas
		}

		sesionID, err := tx.GeneradorID().NextID(ctx)
		if err != nil {
			return err
		}

		accessToken, expiracionAccess, err := tx.TokenServicio().GenerarAccessToken(usuarioID, sesionID, "", "")
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		refreshToken, expiracionRefresh, err := tx.TokenServicio().GenerarRefreshToken(usuarioID, sesionID)
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		accessTokenHash := tx.TokenServicio().HashearToken(accessToken)
		refreshTokenHash := tx.TokenServicio().HashearToken(refreshToken)

		sesion, err := sesiones_domain.NuevaSesion(
			sesionID,
			usuarioID,
			accessTokenHash,
			refreshTokenHash,
			cmd.IPOrigen,
			ahora,
			expiracionAccess,
			expiracionRefresh,
		)
		if err != nil {
			return err
		}

		if _, err := tx.SesionRepositorio().Crear(ctx, sesion); err != nil {
			return err
		}

		credenciales.ResetearIntentos()
		if _, err := tx.CredencialesRepositorio().Actualizar(ctx, credenciales); err != nil {
			return err
		}

		respuesta = &RespuestaLogin{
			AccessToken:       accessToken,
			RefreshToken:      refreshToken,
			ExpiracionAccess:  expiracionAccess,
			ExpiracionRefresh: expiracionRefresh,
			UsuarioID:         usuarioID,
			SesionID:          sesionID,
		}
		return nil
	})

	// 5. Registrar intento fallido por IP si el password fue incorrecto
	if passwordIncorrecto && s.bloqueoIP != nil && cmd.IPOrigen != "" {
		_ = s.bloqueoIP.RegistrarIntentoFallido(ctx, cmd.IPOrigen)
	}

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

func validarComando(cmd ComandoLogin) error {
	if cmd.Email == "" {
		return ErrEmailRequerido
	}
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return ErrEmailInvalido
	}
	if cmd.Password == "" {
		return ErrPasswordRequerido
	}
	return nil
}

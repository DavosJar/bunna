package login

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/rate_limiter"
	"github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
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

type IniciarSesionCasoDeUso struct {
	uow         sesiones_domain.UnitOfWork
	bloqueoIP   *bloqueo_ip.ServicioBloqueoIP
	rateLimiter *rate_limiter.ServicioRateLimit
}

func NewIniciarSesionCasoDeUso(
	uow sesiones_domain.UnitOfWork,
	bloqueoIP *bloqueo_ip.ServicioBloqueoIP,
	rateLimiter *rate_limiter.ServicioRateLimit,
) *IniciarSesionCasoDeUso {
	return &IniciarSesionCasoDeUso{
		uow:         uow,
		bloqueoIP:   bloqueoIP,
		rateLimiter: rateLimiter,
	}
}

func (uc *IniciarSesionCasoDeUso) Ejecutar(ctx context.Context, cmd ComandoIniciarSesion) (*RespuestaIniciarSesion, error) {
	if err := validarComando(cmd); err != nil {
		return nil, err
	}

	if uc.rateLimiter != nil && cmd.IPOrigen != "" {
		if err := uc.rateLimiter.Verificar(ctx, cmd.IPOrigen); err != nil {
			return nil, err
		}
	}

	if uc.bloqueoIP != nil && cmd.IPOrigen != "" {
		if err := uc.bloqueoIP.Verificar(ctx, cmd.IPOrigen); err != nil {
			return nil, err
		}
	}

	var respuesta *RespuestaIniciarSesion
	var passwordIncorrecto bool

	err := uc.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

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

		accessToken, expiracionAccess, err := tx.TokenServicio().GenerarAccessToken(usuarioID, sesionID, nil)
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

		respuesta = &RespuestaIniciarSesion{
			AccessToken:       accessToken,
			RefreshToken:      refreshToken,
			ExpiracionAccess:  expiracionAccess,
			ExpiracionRefresh: expiracionRefresh,
			UsuarioID:         usuarioID,
			SesionID:          sesionID,
		}
		return nil
	})

	if passwordIncorrecto && uc.bloqueoIP != nil && cmd.IPOrigen != "" {
		_ = uc.bloqueoIP.RegistrarIntentoFallido(ctx, cmd.IPOrigen)
	}

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

func validarComando(cmd ComandoIniciarSesion) error {
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

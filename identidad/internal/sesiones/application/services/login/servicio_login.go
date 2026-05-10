package login

import (
	"context"
	"errors"
	"net/mail"
	"time"

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

type ServicioLogin struct {
	uow sesiones_domain.UnitOfWork
}

func NuevoServicioLogin(uow sesiones_domain.UnitOfWork) *ServicioLogin {
	return &ServicioLogin{uow: uow}
}

func (s *ServicioLogin) Ejecutar(ctx context.Context, cmd ComandoLogin) (*RespuestaLogin, error) {
	if err := validarComando(cmd); err != nil {
		return nil, err
	}

	var respuesta *RespuestaLogin

	err := s.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

		usuarios, err := tx.UsuarioRepositorio().Listar(ctx,
			usuario_domain.EspecificacionUsuario{
				ListaLiltros: []usuario_domain.CriterioFiltro{
					{Campo: "correo", Operador: "=", Valor: cmd.Email},
				},
			},
			usuario_domain.Paginacion{Pagina: 1, TamanoPagina: 1},
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
			return ErrCredencialesInvalidas
		}

		sesionID, err := tx.GeneradorID().NextID(ctx)
		if err != nil {
			return err
		}

		accessToken, expiracionAccess, err := tx.TokenServicio().GenerarAccessToken(usuarioID, sesionID)
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
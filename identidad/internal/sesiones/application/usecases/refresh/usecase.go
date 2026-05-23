package refresh

import (
	"context"
	"errors"
	"time"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

var (
	ErrRefreshTokenRequerido     = errors.New("el refresh token es requerido")
	ErrTokenInvalido             = errors.New("token inválido o expirado")
	ErrSesionNoValida            = errors.New("sesión no válida")
	ErrRefreshTokenExpirado      = errors.New("refresh token expirado, inicie sesión nuevamente")
	ErrLimiteRefrescosAlcanzado  = errors.New("límite de refrescos alcanzado, inicie sesión nuevamente")
	ErrSesionAbsolutaExpirada    = errors.New("sesión expirada, inicie sesión nuevamente")
	ErrErrorGenerandoTokens      = errors.New("error al generar tokens")
)

type ConfigRefresh struct {
	MaxRefrescos    int
	TimeoutAbsoluto time.Duration
}

type RenovarSesionCasoDeUso struct {
	uow    sesiones_domain.UnitOfWork
	config ConfigRefresh
}

func NewRenovarSesionCasoDeUso(uow sesiones_domain.UnitOfWork, config ConfigRefresh) *RenovarSesionCasoDeUso {
	return &RenovarSesionCasoDeUso{uow: uow, config: config}
}

func (uc *RenovarSesionCasoDeUso) Ejecutar(ctx context.Context, cmd ComandoRenovarSesion) (*RespuestaRenovarSesion, error) {
	if cmd.RefreshToken == "" {
		return nil, ErrRefreshTokenRequerido
	}

	claims, err := uc.uow.TokenServicio().ValidarRefreshToken(cmd.RefreshToken)
	if err != nil {
		return nil, ErrTokenInvalido
	}

	refreshTokenHash := uc.uow.TokenServicio().HashearToken(cmd.RefreshToken)

	var respuesta *RespuestaRenovarSesion

	err = uc.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

		sesion, err := tx.SesionRepositorio().ObtenerPorRefreshTokenHash(ctx, refreshTokenHash)
		if err != nil {
			_ = tx.SesionRepositorio().InvalidarTodasPorUsuarioID(ctx, claims.UsuarioID)
			return ErrTokenInvalido
		}

		if sesion.Estado() == sesiones_domain.EstadoRevocada || sesion.Estado() == sesiones_domain.EstadoExpirada {
			return ErrSesionNoValida
		}

		if uc.config.TimeoutAbsoluto > 0 {
			if ahora.After(sesion.FechaCreacion().Add(uc.config.TimeoutAbsoluto)) {
				_ = sesion.MarcarExpirada()
				_, _ = tx.SesionRepositorio().Actualizar(ctx, sesion)
				return ErrSesionAbsolutaExpirada
			}
		}

		if !sesion.RefreshTokenValido(ahora) {
			_ = sesion.MarcarExpirada()
			_, _ = tx.SesionRepositorio().Actualizar(ctx, sesion)
			return ErrRefreshTokenExpirado
		}

		if uc.config.MaxRefrescos > 0 && sesion.ContadorRefrescos() >= uc.config.MaxRefrescos {
			return ErrLimiteRefrescosAlcanzado
		}

		nuevoAccessToken, expiracionAccess, err := tx.TokenServicio().GenerarAccessToken(claims.UsuarioID, sesion.ID(), nil)
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		nuevoRefreshToken, expiracionRefresh, err := tx.TokenServicio().GenerarRefreshToken(claims.UsuarioID, sesion.ID())
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		nuevoAccessHash := tx.TokenServicio().HashearToken(nuevoAccessToken)
		nuevoRefreshHash := tx.TokenServicio().HashearToken(nuevoRefreshToken)

		if err := sesion.RotarTokens(nuevoAccessHash, nuevoRefreshHash, expiracionAccess, expiracionRefresh, ahora); err != nil {
			return err
		}

		if _, err := tx.SesionRepositorio().Actualizar(ctx, sesion); err != nil {
			return err
		}

		respuesta = &RespuestaRenovarSesion{
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

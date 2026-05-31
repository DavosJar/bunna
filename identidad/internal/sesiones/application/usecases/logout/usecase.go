package logout

import (
	"context"
	"errors"
	"time"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

var (
	ErrSesionIDRequerido  = errors.New("el ID de sesión es requerido")
	ErrUsuarioIDRequerido = errors.New("el ID de usuario es requerido")
	ErrSesionNoEncontrada = errors.New("sesión no encontrada")
	ErrNoAutorizado       = errors.New("no autorizado para cerrar esta sesión")
)

type CerrarSesionCasoDeUso struct {
	uow sesiones_domain.UnitOfWork
}

func NewCerrarSesionCasoDeUso(uow sesiones_domain.UnitOfWork) *CerrarSesionCasoDeUso {
	return &CerrarSesionCasoDeUso{uow: uow}
}

func (uc *CerrarSesionCasoDeUso) Ejecutar(ctx context.Context, cmd ComandoCerrarSesion) (*RespuestaCerrarSesion, error) {
	if cmd.SesionID == "" {
		return nil, ErrSesionIDRequerido
	}
	if cmd.UsuarioID == "" {
		return nil, ErrUsuarioIDRequerido
	}

	var respuesta *RespuestaCerrarSesion

	err := uc.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		sesion, err := tx.SesionRepositorio().ObtenerPorID(ctx, cmd.SesionID)
		if err != nil {
			return ErrSesionNoEncontrada
		}

		if sesion.UsuarioID() != cmd.UsuarioID {
			return ErrNoAutorizado
		}

		if sesion.Estado() == sesiones_domain.EstadoRevocada ||
			sesion.Estado() == sesiones_domain.EstadoExpirada {
			respuesta = &RespuestaCerrarSesion{SesionesRevocadas: 0}
			return nil
		}

		sesion.Revocar()
		if _, err := tx.SesionRepositorio().Actualizar(ctx, sesion); err != nil {
			return err
		}

		respuesta = &RespuestaCerrarSesion{SesionesRevocadas: 1}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

func (uc *CerrarSesionCasoDeUso) CerrarTodas(ctx context.Context, cmd ComandoCerrarTodasLasSesiones) (*RespuestaCerrarSesion, error) {
	if cmd.UsuarioID == "" {
		return nil, ErrUsuarioIDRequerido
	}

	var respuesta *RespuestaCerrarSesion

	err := uc.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

		sesiones, err := tx.SesionRepositorio().ListarActivasPorUsuarioID(ctx, cmd.UsuarioID, ahora)
		if err != nil {
			return err
		}

		revocadas := 0
		for _, sesion := range sesiones {
			sesion.Revocar()
			if _, err := tx.SesionRepositorio().Actualizar(ctx, sesion); err != nil {
				return err
			}
			revocadas++
		}

		respuesta = &RespuestaCerrarSesion{SesionesRevocadas: revocadas}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

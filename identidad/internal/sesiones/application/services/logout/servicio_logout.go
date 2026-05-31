package logout

import (
	"context"
	"errors"
	"time"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

var (
	// ErrSesionIDRequerido se retorna cuando el sesionID está vacío.
	ErrSesionIDRequerido = errors.New("el ID de sesión es requerido")

	// ErrUsuarioIDRequerido se retorna cuando el usuarioID está vacío.
	ErrUsuarioIDRequerido = errors.New("el ID de usuario es requerido")

	// ErrSesionNoEncontrada se retorna cuando el sesionID no existe en BD.
	ErrSesionNoEncontrada = errors.New("sesión no encontrada")

	// ErrNoAutorizado se retorna cuando la sesión no pertenece al usuario autenticado.
	ErrNoAutorizado = errors.New("no autorizado para cerrar esta sesión")
)

// ServicioLogout implementa el caso de uso de cierre de sesión.
// Soporta cierre de una sesión específica y cierre masivo de todas las sesiones.
type ServicioLogout struct {
	uow sesiones_domain.UnitOfWork
}

// NuevoServicioLogout crea una nueva instancia de ServicioLogout.
func NuevoServicioLogout(uow sesiones_domain.UnitOfWork) *ServicioLogout {
	return &ServicioLogout{uow: uow}
}

// Ejecutar cierra una sesión específica del usuario autenticado.
//
// Flujo:
//  1. Valida que sesionID y usuarioID no estén vacíos.
//  2. Obtiene la sesión por ID.
//  3. Verifica que la sesión pertenezca al usuario.
//  4. Si la sesión ya está revocada → no-op, sin error.
//  5. Si la sesión está expirada → no-op, sin error.
//  6. Si está activa → la revoca y persiste.
func (s *ServicioLogout) Ejecutar(ctx context.Context, cmd ComandoLogout) (*RespuestaLogout, error) {
	if cmd.SesionID == "" {
		return nil, ErrSesionIDRequerido
	}
	if cmd.UsuarioID == "" {
		return nil, ErrUsuarioIDRequerido
	}

	var respuesta *RespuestaLogout

	err := s.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		sesion, err := tx.SesionRepositorio().ObtenerPorID(ctx, cmd.SesionID)
		if err != nil {
			return ErrSesionNoEncontrada
		}

		// Verificar que la sesión pertenece al usuario autenticado
		if sesion.UsuarioID() != cmd.UsuarioID {
			return ErrNoAutorizado
		}

		// No-op si ya está revocada o expirada
		if sesion.Estado() == sesiones_domain.EstadoRevocada ||
			sesion.Estado() == sesiones_domain.EstadoExpirada {
			respuesta = &RespuestaLogout{SesionesRevocadas: 0}
			return nil
		}

		// Revocar sesión activa
		sesion.Revocar()
		if _, err := tx.SesionRepositorio().Actualizar(ctx, sesion); err != nil {
			return err
		}

		respuesta = &RespuestaLogout{SesionesRevocadas: 1}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

// CerrarTodas revoca todas las sesiones activas del usuario.
//
// Flujo:
//  1. Valida que usuarioID no esté vacío.
//  2. Lista todas las sesiones activas del usuario.
//  3. Revoca cada una y las persiste.
func (s *ServicioLogout) CerrarTodas(ctx context.Context, cmd ComandoCerrarTodas) (*RespuestaLogout, error) {
	if cmd.UsuarioID == "" {
		return nil, ErrUsuarioIDRequerido
	}

	var respuesta *RespuestaLogout

	err := s.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
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

		respuesta = &RespuestaLogout{SesionesRevocadas: revocadas}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

// VerificarTimeout verifica si una sesión ha excedido el timeout de inactividad.
// Si lo excedió, marca la sesión como expirada y la persiste.
//
// Esta función es invocada desde la capa de aplicación antes de cada operación
// que requiera sesión activa, no desde el middleware.
func (s *ServicioLogout) VerificarTimeout(ctx context.Context, sesionID string, timeout time.Duration) error {
	return s.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

		sesion, err := tx.SesionRepositorio().ObtenerPorID(ctx, sesionID)
		if err != nil {
			return ErrSesionNoEncontrada
		}

		if sesion.TimeoutExcedido(ahora, timeout) {
			_ = sesion.MarcarExpirada()
			if _, err := tx.SesionRepositorio().Actualizar(ctx, sesion); err != nil {
				return err
			}
		}

		return nil
	})
}

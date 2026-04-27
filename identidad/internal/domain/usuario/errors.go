package usuario

import "errors"

// Errores de validación de usuario
var (
	ErrIDRequerido                       = errors.New("id requerido")
	ErrNombreRequerido                   = errors.New("nombre requerido")
	ErrApellidoRequerido                 = errors.New("apellido requerido")
	ErrCorreoRequerido                   = errors.New("correo requerido")
	ErrTelefonoRequerido                 = errors.New("teléfono requerido")
	ErrTransicionNoPermitida             = errors.New("transición de estado no permitida")
	ErrTransicionVerificacionNoPermitida = errors.New("transición de estado de verificación no permitida")
)

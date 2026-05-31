package verificacion

import "errors"

var (
	ErrEnlaceInvalido        = errors.New("enlace de verificación inválido")
	ErrEnlaceExpirado        = errors.New("enlace de verificación expirado")
	ErrCorreoYaVerificado    = errors.New("el correo ya está verificado")
	ErrDemasiadosReenvios    = errors.New("demasiados intentos de reenvío, intente más tarde")
	ErrUsuarioNoEncontrado   = errors.New("usuario no encontrado")
	ErrVerificacionPendiente = errors.New("no hay verificación pendiente")
)

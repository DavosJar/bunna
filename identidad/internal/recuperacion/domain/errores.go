package recuperacion

import "errors"

var (
	ErrEnlaceInvalido        = errors.New("enlace de recuperación inválido")
	ErrEnlaceExpirado        = errors.New("enlace de recuperación expirado")
	ErrEnlaceYaUtilizado     = errors.New("enlace ya utilizado")
	ErrDemasiadasSolicitudes = errors.New("demasiadas solicitudes, intente más tarde")
	ErrPasswordDebil         = errors.New("la contraseña no cumple los requisitos mínimos")
	ErrEmailRequerido        = errors.New("email requerido")
	ErrEmailInvalido         = errors.New("formato de email inválido")
	ErrUsuarioNoEncontrado   = errors.New("usuario no encontrado")
)

package domain

import "errors"

var (
	// Sesion - constructor
	ErrUsuarioIDRequerido        = errors.New("el usuarioID es requerido")
	ErrAccessTokenHashRequerido  = errors.New("el hash del access token es requerido")
	ErrRefreshTokenHashRequerido = errors.New("el hash del refresh token es requerido")
	ErrFechaExpiracionInvalida   = errors.New("la fecha de expiración no puede ser anterior a la fecha de creación")

	// Sesion - transiciones de estado
	ErrTransicionEstadoInvalida = errors.New("transición de estado no permitida")

	// TokenPair - constructor
	ErrAccessTokenRequerido  = errors.New("el access token es requerido")
	ErrRefreshTokenRequerido = errors.New("el refresh token es requerido")
)

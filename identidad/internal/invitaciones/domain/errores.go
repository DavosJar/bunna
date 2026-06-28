package invitaciones

import "errors"

var (
	ErrEnlaceInvalido       = errors.New("enlace de invitación inválido")
	ErrEnlaceExpirado       = errors.New("el enlace de invitación ha expirado")
	ErrYaAceptada           = errors.New("la invitación ya fue aceptada")
	ErrNoEncontrada         = errors.New("invitación no encontrada")
	ErrEmailRequerido       = errors.New("correo requerido")
	ErrRolRequerido         = errors.New("rol requerido")
	ErrInvitacionNoValida   = errors.New("la invitación no es válida")
	ErrUsuarioNoRegistrado  = errors.New("usuario no registrado: debes crear una cuenta primero")
	ErrYaEliminada          = errors.New("la invitación ya fue eliminada")
)

package dto

type SolicitarRecuperacionRequest struct {
	Correo string `json:"correo" doc:"Correo electrónico registrado" example:"juan@correo.com"`
}

type SolicitarRecuperacionResponse struct {
	Mensaje string `json:"mensaje" doc:"Mensaje informativo" example:"Si el correo está registrado, recibirás un enlace de recuperación"`
}

type ValidarTokenRecuperacionRequest struct {
	Token string `json:"token" doc:"Token de recuperación" example:"abc123..."`
}

type ValidarTokenRecuperacionResponse struct {
	UsuarioID string `json:"usuario_id" doc:"ID del usuario asociado al token"`
	Valido    bool   `json:"valido"     doc:"El token es válido"`
}

type ConfirmarRecuperacionRequest struct {
	Token         string `json:"token"          doc:"Token de recuperación"     example:"abc123..."`
	NuevaPassword string `json:"nueva_password"  doc:"Nueva contraseña"          example:"nueva123"`
}

type ConfirmarRecuperacionResponse struct {
	Mensaje string `json:"mensaje" doc:"Mensaje informativo" example:"Contraseña restablecida exitosamente"`
}

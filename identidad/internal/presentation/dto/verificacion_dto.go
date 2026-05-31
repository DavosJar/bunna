package dto

type SolicitarVerificacionResponse struct {
	Mensaje string `json:"mensaje" doc:"Mensaje informativo" example:"Se ha enviado un enlace de verificación al correo registrado"`
}

type ConfirmarVerificacionRequest struct {
	Token string `json:"token" doc:"Token de verificación recibido por correo" example:"abc123..."`
}

type ConfirmarVerificacionResponse struct {
	Mensaje string `json:"mensaje" doc:"Mensaje informativo" example:"Correo verificado exitosamente"`
}

type ReenviarVerificacionResponse struct {
	Mensaje string `json:"mensaje" doc:"Mensaje informativo" example:"Se ha reenviado el enlace de verificación"`
}

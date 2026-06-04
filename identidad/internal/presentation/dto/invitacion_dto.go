package dto

type CrearInvitacionRequest struct {
	Correo string `json:"correo" doc:"Correo electrónico del invitado" example:"invitado@correo.com"`
	RolID  string `json:"rol_id" doc:"ID del rol a asignar" example:"01926b1e-dead-beef-cafe-000000000001"`
}

type CrearInvitacionResponse struct {
	Mensaje string `json:"mensaje" doc:"Mensaje de confirmación" example:"Invitación enviada exitosamente"`
}

type AceptarInvitacionRequest struct {
	Token string `json:"token" doc:"Token de invitación" example:"01926b1e-dead-beef-cafe-000000000001"`
}

type AceptarInvitacionResponse struct {
	TenantID string `json:"tenant_id" doc:"ID del tenant al que se unió" example:"01926b1e-dead-beef-cafe-000000000002"`
	RolID    string `json:"rol_id" doc:"ID del rol asignado" example:"01926b1e-dead-beef-cafe-000000000003"`
}

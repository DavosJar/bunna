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

type ObtenerInvitacionResponse struct {
	ID           string `json:"id" doc:"ID de la invitación" example:"01926b1e-dead-beef-cafe-000000000001"`
	TenantID     string `json:"tenant_id" doc:"ID del tenant" example:"01926b1e-dead-beef-cafe-000000000002"`
	TenantNombre string `json:"tenant_nombre" doc:"Nombre del tenant" example:"Mi Empresa"`
	RolID        string `json:"rol_id" doc:"ID del rol" example:"01926b1e-dead-beef-cafe-000000000003"`
	RolNombre    string `json:"rol_nombre" doc:"Nombre del rol" example:"Administrador"`
	Email        string `json:"email" doc:"Correo del invitado" example:"invitado@correo.com"`
	Estado       string `json:"estado" doc:"Estado de la invitación (pendiente|aceptada|expirada)" example:"pendiente"`
	Expiracion   string `json:"expiracion" doc:"Fecha de expiración" example:"2026-07-01T00:00:00Z"`
}

type InvitacionItem struct {
	ID            string `json:"id" doc:"ID de la invitación" example:"01926b1e-dead-beef-cafe-000000000001"`
	Email         string `json:"email" doc:"Correo del invitado" example:"invitado@correo.com"`
	Nombre        string `json:"nombre" doc:"Nombre del invitado" example:"Juan Pérez"`
	RolID         string `json:"rol_id" doc:"ID del rol asignado" example:"01926b1e-dead-beef-cafe-000000000003"`
	RolNombre     string `json:"rol_nombre" doc:"Nombre del rol" example:"Administrador"`
	Estado        string `json:"estado" doc:"Estado de la invitación (pendiente|aceptada|expirada)" example:"pendiente"`
	FechaCreacion string `json:"fecha_creacion" doc:"Fecha de creación" example:"2026-06-28T10:00:00Z"`
	Expiracion    string `json:"expiracion" doc:"Fecha de expiración" example:"2026-07-01T00:00:00Z"`
}

type ListarInvitacionesResponse struct {
	Invitaciones []InvitacionItem `json:"invitaciones" doc:"Lista de invitaciones"`
	Total        int              `json:"total" doc:"Número total de invitaciones (sin paginación)"`
}

type ReenviarInvitacionResponse struct {
	Mensaje string `json:"mensaje" doc:"Mensaje de confirmación" example:"Invitación reenviada exitosamente"`
}

type EliminarInvitacionResponse struct {
	Mensaje string `json:"mensaje" doc:"Mensaje de confirmación" example:"Invitación eliminada exitosamente"`
}

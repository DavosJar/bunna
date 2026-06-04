package aceptarInvitacion

type ComandoAceptarInvitacion struct {
	Token     string
	UsuarioID string
}

type RespuestaAceptarInvitacion struct {
	TenantID string
	RolID    string
}

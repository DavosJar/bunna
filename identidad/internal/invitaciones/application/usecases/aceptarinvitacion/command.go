package aceptarInvitacion

type ComandoAceptarInvitacion struct {
	Token string
}

// ToLog returns a safe representation excluding the secret token.
func (c ComandoAceptarInvitacion) ToLog() map[string]any {
	return map[string]any{
		"token_presente": c.Token != "",
	}
}

type RespuestaAceptarInvitacion struct {
	TenantID string
	RolID    string
}

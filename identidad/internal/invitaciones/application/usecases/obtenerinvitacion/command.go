package obtenerinvitacion

type ComandoObtenerInvitacion struct {
	Token string
}

func (c ComandoObtenerInvitacion) ToLog() map[string]any {
	return map[string]any{
		"token_presente": c.Token != "",
	}
}

type RespuestaObtenerInvitacion struct {
	ID           string
	TenantID     string
	TenantNombre string
	RolID        string
	RolNombre    string
	Email        string
	Estado       string // pendiente, aceptada, expirada
	Expiracion   string
}

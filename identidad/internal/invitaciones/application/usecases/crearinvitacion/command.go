package crearInvitacion

type ComandoCrearInvitacion struct {
	TenantID  string
	RolID     string
	Correo    string
	Nombre    string
	CreadoPor string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoCrearInvitacion) ToLog() map[string]any {
	return map[string]any{
		"tenant_id":  c.TenantID,
		"rol_id":     c.RolID,
		"correo":     c.Correo,
		"nombre":     c.Nombre,
		"creado_por": c.CreadoPor,
	}
}

type RespuestaCrearInvitacion struct {
	ID    string
	Token string
}

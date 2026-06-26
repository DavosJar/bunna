package viewcredentials

type ComandoConsultarCredenciales struct {
	UsuarioID  string
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoConsultarCredenciales) ToLog() map[string]any {
	return map[string]any{
		"usuario_id":  c.UsuarioID,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

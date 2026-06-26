package deleteuser

type ComandoDarDeBajaUsuario struct {
	UsuarioID  string
	Motivo     string
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoDarDeBajaUsuario) ToLog() map[string]any {
	return map[string]any{
		"usuario_id":  c.UsuarioID,
		"motivo":      c.Motivo,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

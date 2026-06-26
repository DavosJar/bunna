package revokerole

type ComandoRevocarRol struct {
	UsuarioID  string
	RolID      string
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoRevocarRol) ToLog() map[string]any {
	return map[string]any{
		"usuario_id":  c.UsuarioID,
		"rol_id":      c.RolID,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

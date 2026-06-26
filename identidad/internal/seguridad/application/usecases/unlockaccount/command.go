package unlockaccount

type ComandoDesbloquearCuenta struct {
	UsuarioID  string
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoDesbloquearCuenta) ToLog() map[string]any {
	return map[string]any{
		"usuario_id":  c.UsuarioID,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

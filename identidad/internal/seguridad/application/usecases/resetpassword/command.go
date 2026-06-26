package resetpassword

type ComandoResetearContrasena struct {
	UsuarioID     string
	NuevaPassword string
	TenantID      string
	EjecutorID    string
}

// ToLog returns a safe representation excluding the new password.
func (c ComandoResetearContrasena) ToLog() map[string]any {
	return map[string]any{
		"usuario_id":  c.UsuarioID,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

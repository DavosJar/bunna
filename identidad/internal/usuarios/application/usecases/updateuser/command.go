package updateuser

type ComandoModificarUsuario struct {
	UsuarioID  string
	Nombre     string
	Apellido   string
	Telefono   string
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoModificarUsuario) ToLog() map[string]any {
	return map[string]any{
		"usuario_id":  c.UsuarioID,
		"nombre":      c.Nombre,
		"apellido":    c.Apellido,
		"telefono":    c.Telefono,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

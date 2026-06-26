package switchtenant

// ComandoCambiarTenant contiene los datos necesarios para cambiar de tenant activo.
type ComandoCambiarTenant struct {
	UsuarioID string
	SesionID  string
	TenantID  string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoCambiarTenant) ToLog() map[string]any {
	return map[string]any{
		"usuario_id": c.UsuarioID,
		"sesion_id":  c.SesionID,
		"tenant_id":  c.TenantID,
	}
}

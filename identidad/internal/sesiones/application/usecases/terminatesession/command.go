package terminatesession

type ComandoForzarCierreSesion struct {
	SesionID   string
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoForzarCierreSesion) ToLog() map[string]any {
	return map[string]any{
		"sesion_id":   c.SesionID,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

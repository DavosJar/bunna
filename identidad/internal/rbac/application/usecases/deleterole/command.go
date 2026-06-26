package deleterole

type ComandoEliminarRol struct {
	RolID      string
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoEliminarRol) ToLog() map[string]any {
	return map[string]any{
		"rol_id":      c.RolID,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

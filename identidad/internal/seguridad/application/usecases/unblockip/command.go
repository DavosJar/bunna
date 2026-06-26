package unblockip

type ComandoDesbloquearIP struct {
	IP         string
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoDesbloquearIP) ToLog() map[string]any {
	return map[string]any{
		"ip":          c.IP,
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

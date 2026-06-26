package solicitarrecuperacion

type ComandoSolicitarRecuperacion struct {
	Email    string
	IPOrigen string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoSolicitarRecuperacion) ToLog() map[string]any {
	return map[string]any{
		"email":     c.Email,
		"ip_origen": c.IPOrigen,
	}
}

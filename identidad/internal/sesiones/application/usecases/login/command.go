package login

type ComandoIniciarSesion struct {
	Email    string
	Password string
	IPOrigen string
}

// ToLog returns a safe representation excluding sensitive fields.
func (c ComandoIniciarSesion) ToLog() map[string]any {
	return map[string]any{
		"email":     c.Email,
		"ip_origen": c.IPOrigen,
	}
}

package reenviarverificacion

type ComandoReenviarVerificacion struct {
	UsuarioID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoReenviarVerificacion) ToLog() map[string]any {
	return map[string]any{
		"usuario_id": c.UsuarioID,
	}
}

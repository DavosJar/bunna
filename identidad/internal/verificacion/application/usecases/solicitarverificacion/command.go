package solicitarverificacion

type ComandoSolicitarVerificacion struct {
	UsuarioID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoSolicitarVerificacion) ToLog() map[string]any {
	return map[string]any{
		"usuario_id": c.UsuarioID,
	}
}

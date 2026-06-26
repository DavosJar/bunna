package viewmyprofile

type ComandoVerMiPerfil struct {
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoVerMiPerfil) ToLog() map[string]any {
	return map[string]any{
		"ejecutor_id": c.EjecutorID,
	}
}

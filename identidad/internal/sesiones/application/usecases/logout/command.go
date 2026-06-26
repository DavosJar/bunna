package logout

type ComandoCerrarSesion struct {
	SesionID  string
	UsuarioID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoCerrarSesion) ToLog() map[string]any {
	return map[string]any{
		"sesion_id":  c.SesionID,
		"usuario_id": c.UsuarioID,
	}
}

type ComandoCerrarTodasLasSesiones struct {
	UsuarioID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoCerrarTodasLasSesiones) ToLog() map[string]any {
	return map[string]any{
		"usuario_id": c.UsuarioID,
	}
}

package reenviarverificacion

type ComandoReenviarVerificacion struct {
	Correo string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoReenviarVerificacion) ToLog() map[string]any {
	return map[string]any{
		"correo": c.Correo,
	}
}

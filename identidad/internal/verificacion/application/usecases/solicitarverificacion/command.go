package solicitarverificacion

type ComandoSolicitarVerificacion struct {
	Correo string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoSolicitarVerificacion) ToLog() map[string]any {
	return map[string]any{
		"correo": c.Correo,
	}
}

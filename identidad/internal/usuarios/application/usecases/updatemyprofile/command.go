package updatemyprofile

type ComandoModificarMiPerfil struct {
	EjecutorID string
	Nombre     string
	Apellido   string
	Telefono   string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoModificarMiPerfil) ToLog() map[string]any {
	return map[string]any{
		"ejecutor_id": c.EjecutorID,
		"nombre":      c.Nombre,
		"apellido":    c.Apellido,
		"telefono":    c.Telefono,
	}
}

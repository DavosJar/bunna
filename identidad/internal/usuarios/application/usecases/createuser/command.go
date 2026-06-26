package createuser

type ComandoCrearUsuario struct {
	Correo     string
	Nombre     string
	Apellido   string
	Password   string
	CreatedBy  string
	EjecutorID string
}

// ToLog returns a safe representation excluding the password.
func (c ComandoCrearUsuario) ToLog() map[string]any {
	return map[string]any{
		"correo":      c.Correo,
		"nombre":      c.Nombre,
		"apellido":    c.Apellido,
		"created_by":  c.CreatedBy,
		"ejecutor_id": c.EjecutorID,
	}
}

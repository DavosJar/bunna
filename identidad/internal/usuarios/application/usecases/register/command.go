package register

type ComandoRegistrarUsuario struct {
	Correo   string
	Password string
	Nombre   string
	Apellido string
	Telefono string
}

// ToLog returns a safe representation excluding the password.
func (c ComandoRegistrarUsuario) ToLog() map[string]any {
	return map[string]any{
		"correo":   c.Correo,
		"nombre":   c.Nombre,
		"apellido": c.Apellido,
		"telefono": c.Telefono,
	}
}

package confirmarrecuperacion

type ComandoConfirmarRecuperacion struct {
	Token         string
	NuevaPassword string
}

// ToLog returns a safe representation excluding token and password.
func (c ComandoConfirmarRecuperacion) ToLog() map[string]any {
	return map[string]any{
		"token_presente": c.Token != "",
	}
}

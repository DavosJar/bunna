package validartokenrecuperacion

type ComandoValidarTokenRecuperacion struct {
	Token string
}

// ToLog returns a safe representation excluding the secret token.
func (c ComandoValidarTokenRecuperacion) ToLog() map[string]any {
	return map[string]any{
		"token_presente": c.Token != "",
	}
}

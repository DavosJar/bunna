package confirmarverificacion

type ComandoConfirmarVerificacion struct {
	Token string
}

// ToLog returns a safe representation excluding the secret token.
func (c ComandoConfirmarVerificacion) ToLog() map[string]any {
	return map[string]any{
		"token_presente": c.Token != "",
	}
}

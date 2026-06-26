package refresh

type ComandoRenovarSesion struct {
	RefreshToken string
}

// ToLog returns a safe representation excluding the token itself.
func (c ComandoRenovarSesion) ToLog() map[string]any {
	return map[string]any{
		"token_presente": c.RefreshToken != "",
	}
}

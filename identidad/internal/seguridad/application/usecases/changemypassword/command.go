package changemypassword

type ComandoCambiarMiContrasena struct {
	EjecutorID     string
	PasswordActual string
	NuevaPassword  string
}

// ToLog returns a safe representation excluding passwords.
func (c ComandoCambiarMiContrasena) ToLog() map[string]any {
	return map[string]any{
		"ejecutor_id": c.EjecutorID,
	}
}

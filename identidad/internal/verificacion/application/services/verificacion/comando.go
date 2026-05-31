package verificacion

// ComandoSolicitarVerificacion solicita verificación de correo
type ComandoSolicitarVerificacion struct {
	UsuarioID string
}

// ComandoConfirmarVerificacion confirma verificación con token
type ComandoConfirmarVerificacion struct {
	Token string
}

// ComandoReenviarVerificacion reenvía el email de verificación
type ComandoReenviarVerificacion struct {
	UsuarioID string
}

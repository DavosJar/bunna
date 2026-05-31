package verifyemail

type ComandoSolicitarVerificacion struct {
	UsuarioID string
}

type ComandoConfirmarVerificacion struct {
	Token string
}

type ComandoReenviarVerificacion struct {
	UsuarioID string
}

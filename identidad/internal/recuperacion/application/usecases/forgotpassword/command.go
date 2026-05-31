package forgotpassword

type ComandoSolicitarRecuperacion struct {
	Email    string
	IPOrigen string
}

type ComandoValidarTokenRecuperacion struct {
	Token string
}

type ComandoConfirmarRestablecimiento struct {
	Token         string
	NuevaPassword string
}

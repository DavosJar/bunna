package forgotpassword

type RespuestaSolicitarRecuperacion struct {
	Mensaje string
}

type RespuestaValidarTokenRecuperacion struct {
	UsuarioID string
	Valido    bool
}

type RespuestaConfirmarRestablecimiento struct {
	Mensaje string
}

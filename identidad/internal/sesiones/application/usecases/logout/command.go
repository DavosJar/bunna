package logout

type ComandoCerrarSesion struct {
	SesionID  string
	UsuarioID string
}

type ComandoCerrarTodasLasSesiones struct {
	UsuarioID string
}

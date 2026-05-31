package recuperacion

// ComandoSolicitarRecuperacion solicita restablecimiento de contraseña
type ComandoSolicitarRecuperacion struct {
	Email    string
	IPOrigen string
}

// ComandoValidarToken valida que un token sea válido
type ComandoValidarToken struct {
	Token string
}

// ComandoConfirmarRestablecimiento confirma el restablecimiento con nuevo password
type ComandoConfirmarRestablecimiento struct {
	Token         string
	NuevoPassword string
}

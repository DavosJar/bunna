package recuperacion

// RespuestaSolicitar respuesta genérica al solicitar recuperación
type RespuestaSolicitar struct {
	Mensaje string `json:"mensaje"`
}

// RespuestaValidar respuesta al validar token
type RespuestaValidar struct {
	UsuarioID string `json:"usuario_id"`
	Valido    bool   `json:"valido"`
}

// RespuestaConfirmar respuesta al confirmar restablecimiento
type RespuestaConfirmar struct {
	Mensaje string `json:"mensaje"`
}

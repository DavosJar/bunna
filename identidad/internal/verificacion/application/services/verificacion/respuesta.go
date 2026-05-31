package verificacion

// RespuestaSolicitarVerificacion respuesta al solicitar verificación
type RespuestaSolicitarVerificacion struct {
	Mensaje string `json:"mensaje"`
}

// RespuestaConfirmarVerificacion respuesta al confirmar verificación
type RespuestaConfirmarVerificacion struct {
	Mensaje string `json:"mensaje"`
}

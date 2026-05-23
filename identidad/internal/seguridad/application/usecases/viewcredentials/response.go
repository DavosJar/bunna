package viewcredentials

type RespuestaConsultarCredenciales struct {
	UsuarioID        string
	Activo           bool
	CorreoVerificado bool
	IntentosFallidos int
	BloqueadoHasta   string
}

package expeluser

type RespuestaExpulsarUsuario struct {
	UsuarioID         string
	Estado            string
	SesionesRevocadas int
	ExpulsadoEn       string
}

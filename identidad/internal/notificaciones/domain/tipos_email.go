package notificaciones

// TipoEmail representa el tipo de template de email
type TipoEmail string

const (
	TipoVerificacionCorreo    TipoEmail = "VERIFICACION_CORREO"
	TipoRecuperacionContrasena TipoEmail = "RECUPERACION_CONTRASENA"
)
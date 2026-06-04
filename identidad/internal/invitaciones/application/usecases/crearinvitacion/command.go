package crearInvitacion

type ComandoCrearInvitacion struct {
	TenantID  string
	RolID     string
	Correo    string
	Nombre    string
	CreadoPor string
}

type RespuestaCrearInvitacion struct {
	ID    string
	Token string
}

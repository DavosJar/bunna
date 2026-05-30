package register

import "time"

type RespuestaRegistrarUsuario struct {
	UsuarioID string
	TenantID  string
	Correo    string
	Estado    string
	CreadoEn  time.Time
}

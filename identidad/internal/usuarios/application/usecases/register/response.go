package register

import "time"

type RespuestaRegistrarUsuario struct {
	UsuarioID string
	Correo    string
	Estado    string
	CreadoEn  time.Time
}

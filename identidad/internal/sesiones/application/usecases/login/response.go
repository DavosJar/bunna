package login

import "time"

type RespuestaIniciarSesion struct {
	AccessToken       string
	RefreshToken      string
	ExpiracionAccess  time.Time
	ExpiracionRefresh time.Time
	UsuarioID         string
	SesionID          string
}

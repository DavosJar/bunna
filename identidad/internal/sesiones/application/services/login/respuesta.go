package login

import "time"

type RespuestaLogin struct {
	AccessToken       string
	RefreshToken      string
	ExpiracionAccess  time.Time
	ExpiracionRefresh time.Time
	UsuarioID         string
	SesionID          string
}

package refresh

import "time"

type RespuestaRenovarSesion struct {
	AccessToken       string
	RefreshToken      string
	ExpiracionAccess  time.Time
	ExpiracionRefresh time.Time
	SesionID          string
	UsuarioID         string
}

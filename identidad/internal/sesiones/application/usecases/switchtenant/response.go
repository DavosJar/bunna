package switchtenant

import "time"

// RespuestaCambiarTenant contiene los tokens y datos actualizados tras cambiar de tenant.
type RespuestaCambiarTenant struct {
	AccessToken       string
	RefreshToken      string
	ExpiracionAccess  time.Time
	ExpiracionRefresh time.Time
	UsuarioID         string
	SesionID          string
	TenantID          string
	Rol               string
}

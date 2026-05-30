package domain

import (
	"time"
)

// TokenServicio define el contrato para generar y validar tokens JWT.
type TokenServicio interface {
	GenerarAccessToken(usuarioID, sesionID string, tenantID string, rol string) (tokenString string, expira time.Time, err error)
	GenerarRefreshToken(usuarioID, sesionID string) (tokenString string, expira time.Time, err error)
	ValidarAccessToken(tokenString string) (*TokenClaims, error)
	ValidarRefreshToken(tokenString string) (*TokenClaims, error)
	HashearToken(tokenString string) string
}

type TokenClaims struct {
	UsuarioID string
	SesionID  string
	Tipo      string
	Expira    time.Time
	TenantID  string
	Rol       string
}

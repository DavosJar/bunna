package domain

import (
	"time"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

// TokenServicio define el contrato para generar y validar tokens JWT.
type TokenServicio interface {
	GenerarAccessToken(usuarioID, sesionID string, claims *rbac.UsuarioClaims) (tokenString string, expira time.Time, err error)
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
	Global    bool
	Tenants   map[string]rbac.TenantClaims
}
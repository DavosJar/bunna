package jwt

import (
	"crypto/sha256"
	"fmt"
	"time"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/golang-jwt/jwt/v5"
)

type ConfigJWT struct {
	Secret            string
	Issuer            string
	ExpiracionAccess  time.Duration
	ExpiracionRefresh time.Duration
}

type claimsJWT struct {
	SesionID string                       `json:"sid"`
	Tipo     string                       `json:"typ"`
	Global   bool                         `json:"global"`
	Tenants  map[string]rbac.TenantClaims `json:"tenants,omitempty"`
	jwt.RegisteredClaims
}

type JWTTokenServicio struct {
	config ConfigJWT
}

func NewJWTTokenServicio(config ConfigJWT) sesiones_domain.TokenServicio {
	return &JWTTokenServicio{config: config}
}

func (s *JWTTokenServicio) GenerarAccessToken(usuarioID, sesionID string, claims *rbac.UsuarioClaims) (string, time.Time, error) {
	expira := time.Now().Add(s.config.ExpiracionAccess)

	global := false
	var tenants map[string]rbac.TenantClaims

	if claims != nil {
		global = claims.Global
		tenants = claims.Tenants
	}

	jwtClaims := claimsJWT{
		SesionID: sesionID,
		Tipo:     "access",
		Global:   global,
		Tenants:  tenants,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usuarioID,
			Issuer:    s.config.Issuer,
			ExpiresAt: jwt.NewNumericDate(expira),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error al firmar access token: %w", err)
	}
	return tokenString, expira, nil
}

func (s *JWTTokenServicio) GenerarRefreshToken(usuarioID, sesionID string) (string, time.Time, error) {
	expira := time.Now().Add(s.config.ExpiracionRefresh)
	claims := claimsJWT{
		SesionID: sesionID,
		Tipo:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usuarioID,
			Issuer:    s.config.Issuer,
			ExpiresAt: jwt.NewNumericDate(expira),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error al firmar refresh token: %w", err)
	}
	return tokenString, expira, nil
}

func (s *JWTTokenServicio) ValidarAccessToken(tokenString string) (*sesiones_domain.TokenClaims, error) {
	return s.validarToken(tokenString, "access")
}

func (s *JWTTokenServicio) ValidarRefreshToken(tokenString string) (*sesiones_domain.TokenClaims, error) {
	return s.validarToken(tokenString, "refresh")
}

func (s *JWTTokenServicio) HashearToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return fmt.Sprintf("%x", hash)
}

func (s *JWTTokenServicio) validarToken(tokenString, tipoEsperado string) (*sesiones_domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claimsJWT{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(*claimsJWT)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("claims inválidos")
	}

	if claims.Tipo != tipoEsperado {
		return nil, fmt.Errorf("tipo de token incorrecto: esperaba %s, got %s", tipoEsperado, claims.Tipo)
	}

	return &sesiones_domain.TokenClaims{
		UsuarioID: claims.Subject,
		SesionID:  claims.SesionID,
		Tipo:      claims.Tipo,
		Expira:    claims.ExpiresAt.Time,
		Global:    claims.Global,
		Tenants:   claims.Tenants,
	}, nil
}

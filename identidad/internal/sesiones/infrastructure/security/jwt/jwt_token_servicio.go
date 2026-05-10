// Package jwt implementa TokenServicio usando HMAC-SHA256 (HS256).
package jwt

import (
	"crypto/sha256"
	"fmt"
	"time"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/golang-jwt/jwt/v5"
)

// ConfigJWT contiene la configuración para generar y validar tokens JWT.
type ConfigJWT struct {
	// Secret es la clave HMAC-SHA256 para firmar los tokens.
	Secret string

	// ExpiracionAccess es la duración del access token.
	ExpiracionAccess time.Duration

	// ExpiracionRefresh es la duración del refresh token.
	ExpiracionRefresh time.Duration
}

// claimsJWT define la estructura de claims del JWT.
type claimsJWT struct {
	SesionID string `json:"sid"`
	Tipo     string `json:"typ"`
	jwt.RegisteredClaims
}

// JWTTokenServicio implementa la interfaz TokenServicio usando JWT HS256.
type JWTTokenServicio struct {
	config ConfigJWT
}

// NewJWTTokenServicio crea una nueva instancia de JWTTokenServicio.
func NewJWTTokenServicio(config ConfigJWT) sesiones_domain.TokenServicio {
	return &JWTTokenServicio{config: config}
}

// GenerarAccessToken genera un JWT de tipo access para el usuario y sesión dados.
func (s *JWTTokenServicio) GenerarAccessToken(usuarioID, sesionID string) (string, time.Time, error) {
	expira := time.Now().Add(s.config.ExpiracionAccess)
	claims := claimsJWT{
		SesionID: sesionID,
		Tipo:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usuarioID,
			ExpiresAt: jwt.NewNumericDate(expira),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error al firmar access token: %w", err)
	}
	return tokenString, expira, nil
}

// GenerarRefreshToken genera un JWT de tipo refresh para el usuario y sesión dados.
func (s *JWTTokenServicio) GenerarRefreshToken(usuarioID, sesionID string) (string, time.Time, error) {
	expira := time.Now().Add(s.config.ExpiracionRefresh)
	claims := claimsJWT{
		SesionID: sesionID,
		Tipo:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usuarioID,
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

// ValidarAccessToken valida un access token y retorna sus claims si es válido.
func (s *JWTTokenServicio) ValidarAccessToken(tokenString string) (*sesiones_domain.TokenClaims, error) {
	return s.validarToken(tokenString, "access")
}

// ValidarRefreshToken valida un refresh token y retorna sus claims si es válido.
func (s *JWTTokenServicio) ValidarRefreshToken(tokenString string) (*sesiones_domain.TokenClaims, error) {
	return s.validarToken(tokenString, "refresh")
}

// HashearToken genera un hash SHA-256 del token para almacenamiento seguro.
// El token en plano nunca se persiste en la base de datos.
func (s *JWTTokenServicio) HashearToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return fmt.Sprintf("%x", hash)
}

// validarToken parsea y valida un JWT verificando firma, expiración y tipo.
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
	}, nil
}
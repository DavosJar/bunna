package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Config contiene la configuración para validar tokens JWT emitidos por identidad.
type Config struct {
	Secret string // clave HMAC-SHA256 para verificar firma
	Issuer string // issuer esperado (iss claim), opcional
}

// ClaimsJWT extiende los claims estándar con campos específicos de identidad.
type ClaimsJWT struct {
	SesionID string `json:"sid"`
	Tipo     string `json:"typ"`
	TenantID string `json:"tenant_id"`
	Rol      string `json:"rol"`
	jwt.RegisteredClaims
}

// TokenValidator valida tokens JWT firmados por el servicio de identidad.
type TokenValidator struct {
	config Config
}

// NewTokenValidator crea una nueva instancia de TokenValidator.
func NewTokenValidator(config Config) *TokenValidator {
	return &TokenValidator{config: config}
}

// ValidarAccessToken valida un access token JWT y retorna (usuarioID, sesionID, tenantID, rol, issuer, error).
// Valida: firma, expiración, issuer (si está configurado) y tipo (debe ser "access").
func (v *TokenValidator) ValidarAccessToken(tokenString string) (usuarioID, sesionID, tenantID, rol, issuer string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClaimsJWT{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
		}
		return []byte(v.config.Secret), nil
	})
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(*ClaimsJWT)
	if !ok || !token.Valid {
		return "", "", "", "", "", fmt.Errorf("claims inválidos")
	}

	if claims.Tipo != "access" {
		return "", "", "", "", "", fmt.Errorf("tipo de token incorrecto: se esperaba access, got %s", claims.Tipo)
	}

	// Validar issuer si está configurado
	if v.config.Issuer != "" && claims.Issuer != v.config.Issuer {
		return "", "", "", "", "", fmt.Errorf("issuer inválido: esperado %s, got %s", v.config.Issuer, claims.Issuer)
	}

	return claims.Subject, claims.SesionID, claims.TenantID, claims.Rol, claims.Issuer, nil
}

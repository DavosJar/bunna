package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry"
	jwtvalidator "github.com/davosjar/bunna/services/fincas/internal/infrastructure/security/jwt"
)

// Claves para el contexto Gin
const (
	ClaveAuthContext = "authContext"
)

// AuthMiddleware valida el token Bearer del header Authorization
// y construye el AuthContext para los casos de uso.
type AuthMiddleware struct {
	validator *jwtvalidator.TokenValidator
}

// NewAuthMiddleware crea una nueva instancia de AuthMiddleware.
func NewAuthMiddleware(validator *jwtvalidator.TokenValidator) *AuthMiddleware {
	return &AuthMiddleware{validator: validator}
}

// RequireAuth retorna un handler Gin que exige un token JWT válido.
// Si el token falta o es inválido, responde 401 y aborta.
// En caso exitoso, inyecta *application.AuthContext en el contexto Gin.
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractBearerToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"detalle": err.Error(),
			})
			return
		}

		usuarioID, sesionID, tenantID, rol, issuer, err := m.validator.ValidarAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"detalle": "token inválido o expirado",
			})
			return
		}

		_ = sesionID  // disponible para logging/auditoría
		_ = issuer    // disponible para validación futura
		_ = rol       // disponible para autorización futura

		auth := &application.AuthContext{
			UsuarioID: usuarioID,
			TenantID:  tenantID,
			Permisos:  nil,
		}

		c.Set(ClaveAuthContext, auth)
		reqCtx := telemetry.WithUsuarioID(c.Request.Context(), usuarioID)
		c.Request = c.Request.WithContext(reqCtx)
		c.Next()
	}
}

// GetAuthContext extrae el AuthContext del contexto Gin.
// Retorna nil si no existe (ruta pública sin middleware).
func GetAuthContext(c *gin.Context) *application.AuthContext {
	auth, exists := c.Get(ClaveAuthContext)
	if !exists {
		return nil
	}
	return auth.(*application.AuthContext)
}

// extractBearerToken extrae el token Bearer del header Authorization.
func extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("header Authorization ausente")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("formato inválido, se esperaba: Bearer <token>")
	}

	if parts[1] == "" {
		return "", fmt.Errorf("token vacío")
	}

	return parts[1], nil
}

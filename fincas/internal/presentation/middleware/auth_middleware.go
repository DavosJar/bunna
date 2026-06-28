package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/davosjar/bunna/services/fincas/internal/application"
	iampostgres "github.com/davosjar/bunna/services/fincas/internal/infrastructure/security/iam/postgres"
	jwtvalidator "github.com/davosjar/bunna/services/fincas/internal/infrastructure/security/jwt"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry"
)

// Claves para el contexto Gin
const (
	ClaveAuthContext = "authContext"
)

// AuthMiddleware valida el token Bearer del header Authorization
// y construye el AuthContext para los casos de uso.
type AuthMiddleware struct {
	validator *jwtvalidator.TokenValidator
	iamRepo   *iampostgres.IAMRepositorio
}

// NewAuthMiddleware crea una nueva instancia de AuthMiddleware.
func NewAuthMiddleware(validator *jwtvalidator.TokenValidator, iamRepo *iampostgres.IAMRepositorio) *AuthMiddleware {
	return &AuthMiddleware{
		validator: validator,
		iamRepo:   iamRepo,
	}
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

		_ = sesionID // disponible para logging/auditoría
		_ = issuer   // disponible para validación futura

		log.Printf("[AuthMiddleware] JWT Claims -> UsuarioID: %s, TenantID: %s, Rol: %s", usuarioID, tenantID, rol)

		// Obtener permisos locales desde la base de datos de Fincas para este rol y tenant
		permisos, err := m.iamRepo.ObtenerPermisos(c.Request.Context(), rol, tenantID)
		if err != nil {
			// Si hay error (distinto a no encontrado) o no tiene permisos, dejamos la lista vacía o nil, pero registramos el error si es grave
			permisos = []string{}
		}

		auth := &application.AuthContext{
			UsuarioID: usuarioID,
			TenantID:  tenantID,
			Permisos:  permisos,
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

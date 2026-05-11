// Package middleware contiene los middlewares HTTP de la capa de presentación.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

const (
	// ClaveUsuarioID es la clave para extraer el usuarioID del contexto Gin.
	ClaveUsuarioID = "usuarioID"
	// ClaveSesionID es la clave para extraer el sesionID del contexto Gin.
	ClaveSesionID = "sesionID"
)

// JWTMiddleware valida el token Bearer del header Authorization.
// Si el token es válido, inyecta usuarioID y sesionID en el contexto Gin.
// Si el token es inválido o ausente, responde 401 y aborta.
func JWTMiddleware(tokenSvc sesiones_domain.TokenServicio) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"title":  "Unauthorized",
				"detail": "header Authorization ausente",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"title":  "Unauthorized",
				"detail": "formato de Authorization inválido, se esperaba: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]
		claims, err := tokenSvc.ValidarAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"title":  "Unauthorized",
				"detail": "token inválido o expirado",
			})
			return
		}

		// Inyectar claims en contexto Gin para uso en handlers
		c.Set(ClaveUsuarioID, claims.UsuarioID)
		c.Set(ClaveSesionID, claims.SesionID)
		c.Next()
	}
}

// Package middleware contiene los middlewares HTTP de la capa de presentación.
package middleware

import (
	"context"
	"net/http"
	"strings"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/gin-gonic/gin"
)

// Claves usadas en gin.Context.Set/Get (tipo string para compatibilidad gin).
const (
	ClaveUsuarioID = "usuarioID"
	ClaveSesionID  = "sesionID"
	ClaveTenantID  = "tenantID"
	ClaveRol       = "rol"
)

// ctxKeyUsuarioID y ctxKeySesionID son claves con tipo para context.Context (evita colisiones).
type ctxKey string

const (
	ctxKeyUsuarioID ctxKey = "usuarioID"
	ctxKeySesionID  ctxKey = "sesionID"
	ctxKeyTenantID  ctxKey = "tenantID"
	ctxKeyRol       ctxKey = "rol"
)

// GetUsuarioIDFromCtx extrae el usuarioID del context.Context (útil para handlers Huma).
func GetUsuarioIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyUsuarioID).(string); ok {
		return v
	}
	return ""
}

// GetSesionIDFromCtx extrae el sesionID del context.Context (útil para handlers Huma).
func GetSesionIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeySesionID).(string); ok {
		return v
	}
	return ""
}

// GetTenantIDFromCtx extrae el tenantID del context.Context (útil para handlers Huma).
func GetTenantIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyTenantID).(string); ok {
		return v
	}
	return ""
}

// GetRolFromCtx extrae el rol del context.Context (útil para handlers Huma).
func GetRolFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyRol).(string); ok {
		return v
	}
	return ""
}

// JWTMiddleware valida el token Bearer del header Authorization.
// Si el token es válido, inyecta usuarioID y sesionID en el contexto Gin y en el context.Context.
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

		// Inyectar claims en contexto Gin para handlers gin nativos
		c.Set(ClaveUsuarioID, claims.UsuarioID)
		c.Set(ClaveSesionID, claims.SesionID)
		c.Set(ClaveTenantID, claims.TenantID)
		c.Set(ClaveRol, claims.Rol)

		// Inyectar claims en context.Context para handlers Huma
		reqCtx := context.WithValue(c.Request.Context(), ctxKeyUsuarioID, claims.UsuarioID)
		reqCtx = context.WithValue(reqCtx, ctxKeySesionID, claims.SesionID)
		reqCtx = context.WithValue(reqCtx, ctxKeyTenantID, claims.TenantID)
		reqCtx = context.WithValue(reqCtx, ctxKeyRol, claims.Rol)
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}

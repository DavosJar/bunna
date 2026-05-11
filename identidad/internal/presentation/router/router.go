// Package router configura el router HTTP Gin con Huma v2 y registra todos los endpoints.
package router

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/handlers"
)

// Config contiene la configuración del router.
type Config struct {
	Version     string
	CORSOrigins []string
}

// New construye y retorna el router Gin con Huma configurado.
// Registra todos los endpoints de la API.
func New(facade facades.AuthFacade, cfg Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg.CORSOrigins))

	// Configurar Huma sobre Gin
	api := humagin.New(router, huma.DefaultConfig("Identidad API", cfg.Version))

	// Registrar handlers
	handlers.RegisterHealthHandler(api)
	handlers.NewRegisterHandler(facade).Register(api)
	handlers.NewLoginHandler(facade).Register(api)

	return router
}

// corsMiddleware configura CORS con los orígenes permitidos.
func corsMiddleware(origins []string) gin.HandlerFunc {
	allowed := "*"
	if len(origins) > 0 {
		allowed = origins[0]
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowed)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

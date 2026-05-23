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
func New(all *facades.AllFacades, cfg Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg.CORSOrigins))

	// Configurar Huma sobre Gin
	api := humagin.New(router, huma.DefaultConfig("Identidad API", cfg.Version))

	// Registrar handlers — Sistema
	handlers.RegisterHealthHandler(api)

	// Registrar handlers — Autenticación
	handlers.NewRegisterHandler(all.Auth).Register(api)
	handlers.NewLoginHandler(all.Auth).Register(api)
	handlers.NewRefreshHandler(all.Auth).Register(api)
	handlers.NewLogoutHandler(all.Auth).Register(api)
	handlers.NewLogoutAllHandler(all.Auth).Register(api)

	// Registrar handlers — Usuarios
	handlers.NewCrearUsuarioHandler(all.Usuario).Register(api)
	handlers.NewListarUsuariosHandler(all.Usuario).Register(api)
	handlers.NewModificarUsuarioHandler(all.Usuario).Register(api)
	handlers.NewDarDeBajaUsuarioHandler(all.Usuario).Register(api)
	handlers.NewExpulsarUsuarioHandler(all.Usuario).Register(api)

	// Registrar handlers — Mi Perfil
	handlers.NewVerMiPerfilHandler(all.Usuario).Register(api)
	handlers.NewModificarMiPerfilHandler(all.Usuario).Register(api)
	handlers.NewCambiarMiPasswordHandler(all.Seguridad).Register(api)

	// Registrar handlers — Seguridad
	handlers.NewResetearPasswordHandler(all.Seguridad).Register(api)
	handlers.NewDesbloquearCuentaHandler(all.Seguridad).Register(api)
	handlers.NewListarIPsBloqueadasHandler(all.Seguridad).Register(api)
	handlers.NewDesbloquearIPHandler(all.Seguridad).Register(api)
	handlers.NewConsultarCredencialesHandler(all.Seguridad).Register(api)

	// Registrar handlers — Sesiones
	handlers.NewListarSesionesHandler(all.Sesion).Register(api)
	handlers.NewForzarCierreSesionHandler(all.Sesion).Register(api)

	// Registrar handlers — Roles
	handlers.NewListarRolesHandler(all.Rbac).Register(api)
	handlers.NewCrearRolHandler(all.Rbac).Register(api)
	handlers.NewModificarRolHandler(all.Rbac).Register(api)
	handlers.NewEliminarRolHandler(all.Rbac).Register(api)
	handlers.NewAsignarRolHandler(all.Rbac).Register(api)
	handlers.NewRevocarRolHandler(all.Rbac).Register(api)
	handlers.NewAsignarPermisoARolHandler(all.Rbac).Register(api)
	handlers.NewRevocarPermisoDeRolHandler(all.Rbac).Register(api)

	// Registrar handlers — Tenants
	handlers.NewConfigurarTenantHandler(all.Tenant).Register(api)

	// Registrar handlers — Verificación
	handlers.NewSolicitarVerificacionHandler(all.Verificacion).Register(api)
	handlers.NewConfirmarVerificacionHandler(all.Verificacion).Register(api)
	handlers.NewReenviarVerificacionHandler(all.Verificacion).Register(api)

	// Registrar handlers — Recuperación
	handlers.NewSolicitarRecuperacionHandler(all.Recuperacion).Register(api)
	handlers.NewValidarTokenRecuperacionHandler(all.Recuperacion).Register(api)
	handlers.NewConfirmarRecuperacionHandler(all.Recuperacion).Register(api)

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

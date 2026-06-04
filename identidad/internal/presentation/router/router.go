// Package router configura el router HTTP Gin con Huma v2 y registra todos los endpoints.
package router

import (
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"

	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/handlers"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

// Config contiene la configuración del router.
type Config struct {
	Version           string
	CORSOrigins       []string
	APIGatewayEnabled bool
	TokenSvc          sesiones_domain.TokenServicio
	RateLimitIPMaxRequests int
	RateLimitIPVentana     time.Duration
}

// jwtIfRequired es un wrapper del middleware JWT que lo omite para rutas públicas.
func jwtIfRequired(tokenSvc sesiones_domain.TokenServicio) gin.HandlerFunc {
	jwtMid := middleware.JWTMiddleware(tokenSvc)
	publicPaths := []string{
		"/health",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/refresh",
		"/api/v1/recuperacion/solicitar",
		"/api/v1/recuperacion/validar",
		"/api/v1/recuperacion/confirmar",
		"/api/v1/verificacion/confirmar",
	}

	return func(c *gin.Context) {
		for _, p := range publicPaths {
			if strings.HasPrefix(c.Request.URL.Path, p) {
				c.Next()
				return
			}
		}
		jwtMid(c)
	}
}

// New construye y retorna el router Gin con Huma configurado.
// Registra todos los endpoints de la API.
func New(all *facades.AllFacades, cfg Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg))
	router.Use(middleware.NuevoRateLimitMiddleware(
		cfg.RateLimitIPMaxRequests,
		cfg.RateLimitIPVentana,
	))

	// Cuando está detrás de API Gateway, confiar en el proxy
	// y usar X-Forwarded-For para la IP real del cliente
	if cfg.APIGatewayEnabled {
		router.SetTrustedProxies([]string{"0.0.0.0/0"})
		router.Use(func(c *gin.Context) {
			c.Request.Header.Set("X-Real-IP", c.ClientIP())
			c.Next()
		})
	}

	// Aplicar JWT middleware condicionalmente
	if cfg.TokenSvc != nil {
		router.Use(jwtIfRequired(cfg.TokenSvc))
	}

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
	handlers.NewSwitchTenantHandler(all.Auth).Register(api)

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
	handlers.NewListarPermisosHandler(all.Rbac).Register(api)
	handlers.NewListarMisPermisosHandler(all.Rbac).Register(api)

	// Registrar handlers — Tenants
	handlers.NewConfigurarTenantHandler(all.Tenant).Register(api)
	handlers.NewListarMisTenantsHandler(all.Tenant).Register(api)

	// Registrar handlers — Verificación
	handlers.NewSolicitarVerificacionHandler(all.Verificacion).Register(api)
	handlers.NewConfirmarVerificacionHandler(all.Verificacion).Register(api)
	handlers.NewReenviarVerificacionHandler(all.Verificacion).Register(api)

	// Registrar handlers — Recuperación
	handlers.NewSolicitarRecuperacionHandler(all.Recuperacion).Register(api)
	handlers.NewValidarTokenRecuperacionHandler(all.Recuperacion).Register(api)
	handlers.NewConfirmarRecuperacionHandler(all.Recuperacion).Register(api)

	// Registrar handlers — Invitaciones
	handlers.NewCrearInvitacionHandler(all.Invitacion).Register(api)
	handlers.NewAceptarInvitacionHandler(all.Invitacion).Register(api)

	return router
}

// corsMiddleware configura CORS basado en la configuración del router.
// Si API Gateway está habilitado, el gateway maneja CORS y este middleware
// solo agrega headers si hay orígenes explícitamente configurados.
// Si API Gateway está deshabilitado, permite todos los orígenes (*).
func corsMiddleware(cfg Config) gin.HandlerFunc {
	if cfg.APIGatewayEnabled {
		// El API Gateway maneja CORS; solo agregar headers si hay orígenes configurados.
		if len(cfg.CORSOrigins) == 0 {
			return func(c *gin.Context) {
				c.Next()
			}
		}
		allowed := cfg.CORSOrigins[0]
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

	// Sin API Gateway: permitir todos los orígenes.
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

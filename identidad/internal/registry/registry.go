// Package registry centraliza la inyección de dependencias de la aplicación.
package registry

import (
	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/rate_limiter"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/security/bcrypt"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/login"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/logout"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/refresh"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	sesiones_postgres "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/persistence/postgres"
	sesiones_jwt "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/security/jwt"
	shared_idgenerator "github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/idgenerator"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/services/registro"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

// Registry centraliza todas las dependencias de la aplicación.
type Registry struct {
	// Repositorios
	usuarioRepository      usuario_domain.UsuarioRepositorio
	credencialesRepository seguridad_domain.CredencialesRepositorio
	intentoIPRepository    seguridad_domain.IntentoIPRepositorio

	// Servicios de dominio
	encriptacionServicio seguridad_domain.EncriptacionServicio
	tokenServicio        sesiones_domain.TokenServicio

	// Unit of Work
	usuarioUnitOfWork usuario_domain.UnitOfWork
	sesionUnitOfWork  sesiones_domain.UnitOfWork

	// Servicios de aplicación — sesiones
	ServicioLogin   *login.ServicioLogin
	ServicioRefresh *refresh.ServicioRefresh
	ServicioLogout  *logout.ServicioLogout

	// Servicios de aplicación — registro
	servicioRegistro *registro.ServicioRegistro

	// Servicios de aplicación — seguridad perimetral
	ServicioBloqueoIP *bloqueo_ip.ServicioBloqueoIP
	ServicioRateLimit *rate_limiter.ServicioRateLimit
}

// NewRegistry construye y conecta todas las dependencias de la aplicación.
func NewRegistry(db *gorm.DB, cfg *config.Config) *Registry {
	generadorID := shared_idgenerator.NewUUIDv7Generator()

	usuarioRepo := usuarios_postgres.NewUsuarioRepositorio(db)
	credencialesRepo := seguridad_postgres.NewCredencialesRepositorio(db)
	intentoIPRepo := seguridad_postgres.NewIntentoIPRepositorio(db)
	rateLimitRepo := seguridad_postgres.NewRateLimitRepositorio(db)
	sesionRepo := sesiones_postgres.NewSesionRepositorio(db)

	encriptacion := bcrypt.NewBcryptEncriptacion(cfg.BcryptCost)
	tokenSvc := sesiones_jwt.NewJWTTokenServicio(sesiones_jwt.ConfigJWT{
		Secret:            cfg.JWTSecret,
		ExpiracionAccess:  cfg.JWTAccessExpiracion,
		ExpiracionRefresh: cfg.JWTRefreshExpiracion,
	})

	usuarioUoW := usuarios_postgres.NewUnitOfWork(
		db,
		usuarioRepo,
		credencialesRepo,
		encriptacion,
		generadorID,
	)

	sesionUoW := sesiones_postgres.NewSesionUnitOfWork(
		db,
		sesionRepo,
		credencialesRepo,
		usuarioRepo,
		encriptacion,
		tokenSvc,
		generadorID,
	)

	bloqueoIPSvc := bloqueo_ip.NuevoServicioBloqueoIP(
		intentoIPRepo,
		generadorID,
		bloqueo_ip.ConfigBloqueoIP{
			MaxIntentos: cfg.BloqueoIPMaxIntentos,
			Ventana:     cfg.BloqueoIPVentana,
			Duracion:    cfg.BloqueoIPDuracion,
		},
	)
	rateLimitSvc := rate_limiter.NuevoServicioRateLimit(
		rateLimitRepo,
		generadorID,
		rate_limiter.ConfigRateLimit{
			MaxRequests: cfg.RateLimitMaxRequests,
			Ventana:     cfg.RateLimitVentana,
		},
	)

	loginSvc := login.NuevoServicioLogin(sesionUoW, bloqueoIPSvc, rateLimitSvc)
	refreshSvc := refresh.NuevoServicioRefresh(sesionUoW, refresh.ConfigRefresh{
		MaxRefrescos:    cfg.SesionMaxRefrescos,
		TimeoutAbsoluto: cfg.SesionTimeoutAbsoluto,
	})
	logoutSvc := logout.NuevoServicioLogout(sesionUoW)
	registroSvc := registro.NuevoServicioRegistro(usuarioUoW)

	return &Registry{
		usuarioRepository:      usuarioRepo,
		credencialesRepository: credencialesRepo,
		intentoIPRepository:    intentoIPRepo,
		encriptacionServicio:   encriptacion,
		tokenServicio:          tokenSvc,
		usuarioUnitOfWork:      usuarioUoW,
		sesionUnitOfWork:       sesionUoW,
		ServicioLogin:          loginSvc,
		ServicioRefresh:        refreshSvc,
		ServicioLogout:         logoutSvc,
		servicioRegistro:       registroSvc,
		ServicioBloqueoIP:      bloqueoIPSvc,
		ServicioRateLimit:      rateLimitSvc,
	}
}

// Getters

func (r *Registry) UsuarioRepository() usuario_domain.UsuarioRepositorio             { return r.usuarioRepository }
func (r *Registry) CredencialesRepository() seguridad_domain.CredencialesRepositorio { return r.credencialesRepository }
func (r *Registry) EncriptacionServicio() seguridad_domain.EncriptacionServicio       { return r.encriptacionServicio }
func (r *Registry) UsuarioUnitOfWork() usuario_domain.UnitOfWork                      { return r.usuarioUnitOfWork }
func (r *Registry) TokenServicio() sesiones_domain.TokenServicio                      { return r.tokenServicio }
func (r *Registry) GetServicioRegistro() *registro.ServicioRegistro                   { return r.servicioRegistro }

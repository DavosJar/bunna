package registry

import (
	"github.com/davosjar/bunna/services/identidad/internal/config"
	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	notificaciones_email "github.com/davosjar/bunna/services/identidad/internal/notificaciones/infrastructure/email"
	checkpermission "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/checkpermission"
	uc_listarmispermisos "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listarmispermisos"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/assignpermissiontorole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/assignrole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/createrole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/deleterole"
	listroles "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listroles"
	uc_listpermisos "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listpermisos"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/revokepermissionfromrole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/revokerole"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/updaterole"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	rbac_postgres "github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/persistence/postgres"
	uc_forgotpassword "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/forgotpassword"
	recuperacion_postgres "github.com/davosjar/bunna/services/identidad/internal/recuperacion/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/rate_limiter"
	uc_changemypassword "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/changemypassword"
	uc_listblockedips "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/listblockedips"
	uc_resetpassword "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/resetpassword"
	uc_unblockip "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/unblockip"
	uc_unlockaccount "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/unlockaccount"
	uc_viewcredentials "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/viewcredentials"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/security/bcrypt"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/login"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/logout"
	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/refresh"
	uc_listsessions "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/listsessions"
	uc_sesiones_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	uc_sesiones_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	uc_sesiones_refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	uc_sesiones_switchtenant "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/switchtenant"
	uc_terminatesession "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/terminatesession"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	sesiones_postgres "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/persistence/postgres"
	sesiones_jwt "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/security/jwt"
	shared_idgenerator "github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/idgenerator"
	uc_updatetenant "github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/updatetenant"
	tenant_domain "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	tenant_postgres "github.com/davosjar/bunna/services/identidad/internal/tenants/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/services/registro"
	uc_register "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/register"
	uc_createuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/createuser"
	uc_deleteuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/deleteuser"
	uc_expeluser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/expeluser"
	uc_listusers "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/listusers"
	uc_updatemyprofile "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/updatemyprofile"
	uc_updateuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/updateuser"
	uc_viewmyprofile "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/viewmyprofile"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	uc_verifyemail "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/verifyemail"
	verificacion_domain "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
	verificacion_postgres "github.com/davosjar/bunna/services/identidad/internal/verificacion/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

type Registry struct {
	// Repositorios
	usuarioRepository      usuario_domain.UsuarioRepositorio
	credencialesRepository seguridad_domain.CredencialesRepositorio
	intentoIPRepository    seguridad_domain.IntentoIPRepositorio
	sesionRepository       sesiones_domain.SesionRepositorio
	rolRepository          rbac.RolRepositorio
	permisoRepository      rbac.PermisoRepositorio
	rolPermisoRepository   rbac.RolPermisoRepositorio
	usuarioRolRepository   rbac.UsuarioRolRepositorio
	usuarioTenantRolRepo   rbac.UsuarioTenantRolRepositorio
	tenantRepository       tenant_domain.TenantRepositorio
	membresiaRepository    tenant_domain.MembresiaRepositorio
	verificacionRepo       verificacion_domain.VerificacionRepositorio

	// Servicios de dominio
	encriptacionServicio seguridad_domain.EncriptacionServicio
	tokenServicio        sesiones_domain.TokenServicio
	authService          *checkpermission.VerificarPermisoCasoDeUso
	emailServicio        notificaciones.EmailServicio

	// Unit of Work
	usuarioUnitOfWork usuario_domain.UnitOfWork
	sesionUnitOfWork  sesiones_domain.UnitOfWork

	// Servicios de aplicación — sesiones (antiguos)
	ServicioLogin   *login.ServicioLogin
	ServicioRefresh *refresh.ServicioRefresh
	ServicioLogout  *logout.ServicioLogout

	// Servicios de aplicación — registro
	servicioRegistro             *registro.ServicioRegistro
	RegistrarUsuarioCasoDeUso    *uc_register.RegistrarUsuarioCasoDeUso

	// Servicios de aplicación — seguridad perimetral
	ServicioBloqueoIP *bloqueo_ip.ServicioBloqueoIP
	ServicioRateLimit *rate_limiter.ServicioRateLimit

	// Casos de uso — auth
	IniciarSesionCasoDeUso    *uc_sesiones_login.IniciarSesionCasoDeUso
	CerrarSesionCasoDeUso     *uc_sesiones_logout.CerrarSesionCasoDeUso
	RenovarSesionCasoDeUso    *uc_sesiones_refresh.RenovarSesionCasoDeUso
	CambiarTenantCasoDeUso    *uc_sesiones_switchtenant.CambiarTenantCasoDeUso

	// Casos de uso — usuarios admin
	CrearUsuarioCasoDeUso     *uc_createuser.CrearUsuarioCasoDeUso
	ListarUsuariosCasoDeUso   *uc_listusers.ListarUsuariosCasoDeUso
	ModificarUsuarioCasoDeUso *uc_updateuser.ModificarUsuarioCasoDeUso
	DarDeBajaUsuarioCasoDeUso *uc_deleteuser.DarDeBajaUsuarioCasoDeUso
	ExpulsarUsuarioCasoDeUso  *uc_expeluser.ExpulsarUsuarioCasoDeUso

	// Casos de uso — autogestión
	VerMiPerfilCasoDeUso         *uc_viewmyprofile.VerMiPerfilCasoDeUso
	ModificarMiPerfilCasoDeUso   *uc_updatemyprofile.ModificarMiPerfilCasoDeUso
	CambiarMiContrasenaCasoDeUso *uc_changemypassword.CambiarMiContrasenaCasoDeUso

	// Casos de uso — seguridad
	ConsultarCredencialesCasoDeUso *uc_viewcredentials.ConsultarCredencialesCasoDeUso
	ResetearContrasenaCasoDeUso    *uc_resetpassword.ResetearContrasenaCasoDeUso
	DesbloquearCuentaCasoDeUso     *uc_unlockaccount.DesbloquearCuentaCasoDeUso
	ListarIPsBloqueadasCasoDeUso   *uc_listblockedips.ListarIPsBloqueadasCasoDeUso
	DesbloquearIPCasoDeUso         *uc_unblockip.DesbloquearIPCasoDeUso

	// Casos de uso — sesiones
	ListarSesionesCasoDeUso     *uc_listsessions.ListarSesionesCasoDeUso
	ForzarCierreSesionCasoDeUso *uc_terminatesession.ForzarCierreSesionCasoDeUso

	// Casos de uso — roles y permisos
	ListarRolesCasoDeUso         *listroles.ListarRolesCasoDeUso
	ListarPermisosCasoDeUso      *uc_listpermisos.ListarPermisosCasoDeUso
	ListarMisPermisosCasoDeUso   *uc_listarmispermisos.ListarMisPermisosCasoDeUso
	CrearRolCasoDeUso            *createrole.CrearRolCasoDeUso
	ModificarRolCasoDeUso        *updaterole.ModificarRolCasoDeUso
	EliminarRolCasoDeUso         *deleterole.EliminarRolCasoDeUso
	AsignarRolCasoDeUso          *assignrole.AsignarRolCasoDeUso
	RevocarRolCasoDeUso          *revokerole.RevocarRolCasoDeUso
	AsignarPermisoARolCasoDeUso  *assignpermissiontorole.AsignarPermisoARolCasoDeUso
	RevocarPermisoDeRolCasoDeUso *revokepermissionfromrole.RevocarPermisoDeRolCasoDeUso

	// Casos de uso — tenants
	ConfigurarTenantCasoDeUso *uc_updatetenant.ConfigurarTenantCasoDeUso

	// Casos de uso — verificación y recuperación
	VerificarCorreoCasoDeUso     *uc_verifyemail.VerificarCorreoCasoDeUso
	RecuperarContrasenaCasoDeUso *uc_forgotpassword.RecuperarContrasenaCasoDeUso
}

func NewRegistry(db *gorm.DB, cfg *config.Config) *Registry {
	generadorID := shared_idgenerator.NewUUIDv7Generator()

	usuarioRepo := usuarios_postgres.NewUsuarioRepositorio(db)
	credencialesRepo := seguridad_postgres.NewCredencialesRepositorio(db)
	intentoIPRepo := seguridad_postgres.NewIntentoIPRepositorio(db)
	rateLimitRepo := seguridad_postgres.NewRateLimitRepositorio(db)
	sesionRepo := sesiones_postgres.NewSesionRepositorio(db)
	rolRepo := rbac_postgres.NewRolRepositorio(db)
	permisoRepo := rbac_postgres.NewPermisoRepositorio(db)
	rolPermisoRepo := rbac_postgres.NewRolPermisoRepositorio(db)
	usuarioRolRepo := rbac_postgres.NewUsuarioRolRepositorio(db)
	usuarioTenantRolRepo := rbac_postgres.NewUsuarioTenantRolRepositorio(db)
	tenantRepo := tenant_postgres.NewTenantRepositorio(db)
	membresiaRepo := tenant_postgres.NewMembresiaRepositorio(db)
	verificacionRepo := verificacion_postgres.NewVerificacionRepositorio(db)
	tokenRecuperacionRepo := recuperacion_postgres.NewTokenRecuperacionRepositorio(db)
	usuarioRecuperacionRepo := recuperacion_postgres.NewUsuarioRecuperacionRepositorio(db)

	encriptacion := bcrypt.NewBcryptEncriptacion(cfg.BcryptCost)
	tokenSvc := sesiones_jwt.NewJWTTokenServicio(sesiones_jwt.ConfigJWT{
		Secret:            cfg.JWTSecret,
		Issuer:            cfg.JWTIssuer,
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

	authSvc := checkpermission.NewVerificarPermisoCasoDeUso(usuarioRolRepo, usuarioTenantRolRepo, permisoRepo)

	emailSvc := notificaciones_email.NewSMTPServicio(notificaciones_email.ConfigSMTP{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
		Async:    false,
	})

	// Servicios de aplicación antiguos
	loginSvc := login.NuevoServicioLogin(sesionUoW, bloqueoIPSvc, rateLimitSvc)
	refreshSvc := refresh.NuevoServicioRefresh(sesionUoW, refresh.ConfigRefresh{
		MaxRefrescos:    cfg.SesionMaxRefrescos,
		TimeoutAbsoluto: cfg.SesionTimeoutAbsoluto,
	})
	logoutSvc := logout.NuevoServicioLogout(sesionUoW)
	registroSvc := registro.NuevoServicioRegistro(usuarioUoW)
	registroUseCase := uc_register.NewRegistrarUsuarioCasoDeUso(
		usuarioRepo,
		credencialesRepo,
		encriptacion,
		generadorID,
		tenantRepo,
		membresiaRepo,
		rolRepo,
		usuarioTenantRolRepo,
	)

	listarMisPermisosCasoDeUso := uc_listarmispermisos.NewListarMisPermisosCasoDeUso(rolRepo, rolPermisoRepo)

	return &Registry{
		usuarioRepository:      usuarioRepo,
		credencialesRepository: credencialesRepo,
		intentoIPRepository:    intentoIPRepo,
		sesionRepository:       sesionRepo,
		rolRepository:          rolRepo,
		permisoRepository:      permisoRepo,
		rolPermisoRepository:   rolPermisoRepo,
		usuarioRolRepository:   usuarioRolRepo,
		usuarioTenantRolRepo:   usuarioTenantRolRepo,
		tenantRepository:       tenantRepo,
		membresiaRepository:    membresiaRepo,
		verificacionRepo:       verificacionRepo,

		encriptacionServicio: encriptacion,
		tokenServicio:        tokenSvc,
		authService:          authSvc,
		emailServicio:        emailSvc,

		usuarioUnitOfWork: usuarioUoW,
		sesionUnitOfWork:  sesionUoW,

		ServicioLogin:     loginSvc,
		ServicioRefresh:   refreshSvc,
		ServicioLogout:    logoutSvc,
		servicioRegistro:             registroSvc,
		RegistrarUsuarioCasoDeUso:    registroUseCase,
		ServicioBloqueoIP: bloqueoIPSvc,
		ServicioRateLimit: rateLimitSvc,

		// Casos de uso — auth
		IniciarSesionCasoDeUso: uc_sesiones_login.NewIniciarSesionCasoDeUso(
			sesionUoW, bloqueoIPSvc, rateLimitSvc,
			uc_sesiones_login.ConfigLogin{
				CuentaMaxIntentos:     cfg.CuentaBloqueoMaxIntentos,
				CuentaBloqueoDuracion: cfg.CuentaBloqueoDuracion,
			},
			membresiaRepo,
			usuarioTenantRolRepo,
		),
		CerrarSesionCasoDeUso:  uc_sesiones_logout.NewCerrarSesionCasoDeUso(sesionUoW),
		RenovarSesionCasoDeUso: uc_sesiones_refresh.NewRenovarSesionCasoDeUso(
			sesionUoW,
			uc_sesiones_refresh.ConfigRefresh{
				MaxRefrescos:    cfg.SesionMaxRefrescos,
				TimeoutAbsoluto: cfg.SesionTimeoutAbsoluto,
			},
			membresiaRepo,
			usuarioTenantRolRepo,
		),

		// Casos de uso — usuarios admin
		CrearUsuarioCasoDeUso:     uc_createuser.NewCrearUsuarioCasoDeUso(usuarioRepo, credencialesRepo, encriptacion, authSvc, generadorID),
		ListarUsuariosCasoDeUso:   uc_listusers.NewListarUsuariosCasoDeUso(usuarioRepo, authSvc),
		ModificarUsuarioCasoDeUso: uc_updateuser.NewModificarUsuarioCasoDeUso(usuarioRepo, authSvc),
		DarDeBajaUsuarioCasoDeUso: uc_deleteuser.NewDarDeBajaUsuarioCasoDeUso(usuarioRepo, authSvc),
		ExpulsarUsuarioCasoDeUso:  uc_expeluser.NewExpulsarUsuarioCasoDeUso(usuarioRepo, sesionRepo, authSvc),

		// Casos de uso — autogestión
		VerMiPerfilCasoDeUso:         uc_viewmyprofile.NewVerMiPerfilCasoDeUso(usuarioRepo),
		ModificarMiPerfilCasoDeUso:   uc_updatemyprofile.NewModificarMiPerfilCasoDeUso(usuarioRepo),
		CambiarMiContrasenaCasoDeUso: uc_changemypassword.NewCambiarMiContrasenaCasoDeUso(credencialesRepo, encriptacion),

		// Casos de uso — seguridad
		ConsultarCredencialesCasoDeUso: uc_viewcredentials.NewConsultarCredencialesCasoDeUso(credencialesRepo, authSvc),
		ResetearContrasenaCasoDeUso:    uc_resetpassword.NewResetearContrasenaCasoDeUso(credencialesRepo, sesionRepo, encriptacion, authSvc),
		DesbloquearCuentaCasoDeUso:     uc_unlockaccount.NewDesbloquearCuentaCasoDeUso(credencialesRepo, authSvc),
		ListarIPsBloqueadasCasoDeUso:   uc_listblockedips.NewListarIPsBloqueadasCasoDeUso(intentoIPRepo, authSvc),
		DesbloquearIPCasoDeUso:         uc_unblockip.NewDesbloquearIPCasoDeUso(intentoIPRepo, authSvc),

		// Casos de uso — sesiones
		ListarSesionesCasoDeUso:     uc_listsessions.NewListarSesionesCasoDeUso(sesionRepo, authSvc),
		ForzarCierreSesionCasoDeUso: uc_terminatesession.NewForzarCierreSesionCasoDeUso(sesionRepo, authSvc),

		// Casos de uso — switch tenant
		CambiarTenantCasoDeUso: uc_sesiones_switchtenant.NewCambiarTenantCasoDeUso(
			membresiaRepo,
			usuarioTenantRolRepo,
			sesionUoW,
		),

		// Casos de uso — tenants
		ConfigurarTenantCasoDeUso: uc_updatetenant.NewConfigurarTenantCasoDeUso(tenantRepo, authSvc),

		// Casos de uso — verificación y recuperación
		VerificarCorreoCasoDeUso: uc_verifyemail.NewVerificarCorreoCasoDeUso(
			verificacionRepo, emailSvc, generadorID,
			uc_verifyemail.ConfigVerificacion{
					FrontendURL:     cfg.FrontendURL,
				TokenExpiracion: cfg.VerificacionTokenExpiracion,
				MaxReenvios:     cfg.VerificacionMaxReenvios,
				VentanaReenvios: cfg.VerificacionVentanaReenvios,
			},
		),
		RecuperarContrasenaCasoDeUso: uc_forgotpassword.NewRecuperarContrasenaCasoDeUso(
			tokenRecuperacionRepo, usuarioRecuperacionRepo, sesionRepo, credencialesRepo, encriptacion, emailSvc, generadorID,
			uc_forgotpassword.ConfigRecuperacion{
				TokenExpiracion:     cfg.RecuperacionTokenExpiracion,
				RateLimitIPMax:      cfg.RecuperacionRateLimitIPMax,
				RateLimitUsuarioMax: cfg.RecuperacionRateLimitUsuarioMax,
				RateLimitVentana:    cfg.RecuperacionRateLimitVentana,
			},
		),

		// Casos de uso — roles y permisos
		ListarRolesCasoDeUso:         listroles.NewListarRolesCasoDeUso(rolRepo, permisoRepo, authSvc),
			ListarPermisosCasoDeUso:      uc_listpermisos.NewListarPermisosCasoDeUso(permisoRepo, authSvc),
		ListarMisPermisosCasoDeUso:   listarMisPermisosCasoDeUso,
		CrearRolCasoDeUso:            createrole.NewCrearRolCasoDeUso(rolRepo, permisoRepo, rolPermisoRepo, authSvc),
		ModificarRolCasoDeUso:        updaterole.NewModificarRolCasoDeUso(rolRepo, authSvc),
		EliminarRolCasoDeUso:         deleterole.NewEliminarRolCasoDeUso(rolRepo, authSvc),
		AsignarRolCasoDeUso:          assignrole.NewAsignarRolCasoDeUso(usuarioRolRepo, usuarioTenantRolRepo, rolRepo, authSvc),
		RevocarRolCasoDeUso:          revokerole.NewRevocarRolCasoDeUso(usuarioRolRepo, usuarioTenantRolRepo, authSvc),
		AsignarPermisoARolCasoDeUso:  assignpermissiontorole.NewAsignarPermisoARolCasoDeUso(rolRepo, permisoRepo, rolPermisoRepo, authSvc),
		RevocarPermisoDeRolCasoDeUso: revokepermissionfromrole.NewRevocarPermisoDeRolCasoDeUso(rolRepo, permisoRepo, rolPermisoRepo, authSvc),
	}
}

// Getters
func (r *Registry) UsuarioRepository() usuario_domain.UsuarioRepositorio { return r.usuarioRepository }
func (r *Registry) CredencialesRepository() seguridad_domain.CredencialesRepositorio {
	return r.credencialesRepository
}
func (r *Registry) EncriptacionServicio() seguridad_domain.EncriptacionServicio {
	return r.encriptacionServicio
}
func (r *Registry) UsuarioUnitOfWork() usuario_domain.UnitOfWork    { return r.usuarioUnitOfWork }
func (r *Registry) TokenServicio() sesiones_domain.TokenServicio    { return r.tokenServicio }
func (r *Registry) GetServicioRegistro() *registro.ServicioRegistro { return r.servicioRegistro }
func (r *Registry) GetRegistrarUsuarioCasoDeUso() *uc_register.RegistrarUsuarioCasoDeUso {
	return r.RegistrarUsuarioCasoDeUso
}
func (r *Registry) TenantRepository() tenant_domain.TenantRepositorio { return r.tenantRepository }
func (r *Registry) MembresiaRepository() tenant_domain.MembresiaRepositorio {
	return r.membresiaRepository
}
func (r *Registry) AuthService() *checkpermission.VerificarPermisoCasoDeUso { return r.authService }
func (r *Registry) EmailServicio() notificaciones.EmailServicio     { return r.emailServicio }
func (r *Registry) UsuarioTenantRolRepositorio() rbac.UsuarioTenantRolRepositorio {
	return r.usuarioTenantRolRepo
}

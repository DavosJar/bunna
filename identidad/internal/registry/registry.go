package registry

import (
	"context"
	"os"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	uc_aceptarInvitacion "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/aceptarinvitacion"
	uc_crearInvitacion "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/crearinvitacion"
	uc_eliminarInvitacion "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/eliminarinvitacion"
	uc_listarInvitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/listarinvitaciones"
	uc_obtenerInvitacion "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/obtenerinvitacion"
	uc_reenviarInvitacion "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/reenviarinvitacion"
	invitaciones_postgres "github.com/davosjar/bunna/services/identidad/internal/invitaciones/infrastructure/persistence/postgres"
	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	notificaciones_email "github.com/davosjar/bunna/services/identidad/internal/notificaciones/infrastructure/email"
	checkpermission "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/checkpermission"
	uc_listarmispermisos "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listarmispermisos"
	uc_listarrolesdeusuario "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listarrolesdeusuario"
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
	uc_confirmar_recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/confirmarrecuperacion"
	uc_solicitar_recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/solicitarrecuperacion"
	uc_validar_recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/validartokenrecuperacion"
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
	uc_listarmistenants "github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/uc_listarmistenants"
	uc_obtenertenantporid "github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/uc_obtenertenantporid"
	uc_obtenertenantporslug "github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/uc_obtenertenantporslug"
	tenant_domain "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	tenant_postgres "github.com/davosjar/bunna/services/identidad/internal/tenants/infrastructure/persistence/postgres"
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
	uc_confirmar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/confirmarverificacion"
	uc_reenviar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/reenviarverificacion"
	uc_solicitar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/solicitarverificacion"
	verificacion_domain "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
	verificacion_postgres "github.com/davosjar/bunna/services/identidad/internal/verificacion/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/consumers"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/publishers"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/decorator"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/gormplugin"
	"strings"

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

	RegistrarUsuarioCasoDeUso    decorator.UseCase[*uc_register.ComandoRegistrarUsuario, *uc_register.RespuestaRegistrarUsuario]

	// Servicios de aplicación — seguridad perimetral
	ServicioBloqueoIP *bloqueo_ip.ServicioBloqueoIP
	ServicioRateLimit *rate_limiter.ServicioRateLimit

	// Casos de uso — auth
	IniciarSesionCasoDeUso    decorator.UseCase[uc_sesiones_login.ComandoIniciarSesion, *uc_sesiones_login.RespuestaIniciarSesion]
	CerrarSesionCasoDeUso     decorator.LogoutUseCase
	RenovarSesionCasoDeUso    decorator.UseCase[uc_sesiones_refresh.ComandoRenovarSesion, *uc_sesiones_refresh.RespuestaRenovarSesion]
	CambiarTenantCasoDeUso    decorator.UseCase[uc_sesiones_switchtenant.ComandoCambiarTenant, *uc_sesiones_switchtenant.RespuestaCambiarTenant]

	// Casos de uso — usuarios admin
	CrearUsuarioCasoDeUso     decorator.UseCase[*uc_createuser.ComandoCrearUsuario, *uc_createuser.RespuestaCrearUsuario]
	ListarUsuariosCasoDeUso   decorator.UseCase[*uc_listusers.ComandoListarUsuarios, *uc_listusers.RespuestaListarUsuarios]
	ModificarUsuarioCasoDeUso decorator.UseCase[*uc_updateuser.ComandoModificarUsuario, *uc_updateuser.RespuestaModificarUsuario]
	DarDeBajaUsuarioCasoDeUso decorator.UseCase[*uc_deleteuser.ComandoDarDeBajaUsuario, *uc_deleteuser.RespuestaDarDeBajaUsuario]
	ExpulsarUsuarioCasoDeUso  decorator.UseCase[*uc_expeluser.ComandoExpulsarUsuario, *uc_expeluser.RespuestaExpulsarUsuario]

	// Casos de uso — autogestión
	VerMiPerfilCasoDeUso         decorator.UseCase[*uc_viewmyprofile.ComandoVerMiPerfil, *uc_viewmyprofile.RespuestaVerMiPerfil]
	ModificarMiPerfilCasoDeUso   decorator.UseCase[*uc_updatemyprofile.ComandoModificarMiPerfil, *uc_updatemyprofile.RespuestaModificarMiPerfil]
	CambiarMiContrasenaCasoDeUso decorator.UseCase[*uc_changemypassword.ComandoCambiarMiContrasena, *uc_changemypassword.RespuestaCambiarMiContrasena]

	// Casos de uso — seguridad
	ConsultarCredencialesCasoDeUso decorator.UseCase[*uc_viewcredentials.ComandoConsultarCredenciales, *uc_viewcredentials.RespuestaConsultarCredenciales]
	ResetearContrasenaCasoDeUso    decorator.UseCase[*uc_resetpassword.ComandoResetearContrasena, *uc_resetpassword.RespuestaResetearContrasena]
	DesbloquearCuentaCasoDeUso     decorator.UseCase[*uc_unlockaccount.ComandoDesbloquearCuenta, *uc_unlockaccount.RespuestaDesbloquearCuenta]
	ListarIPsBloqueadasCasoDeUso   decorator.UseCase[*uc_listblockedips.ComandoListarIPsBloqueadas, *uc_listblockedips.RespuestaListarIPsBloqueadas]
	DesbloquearIPCasoDeUso         decorator.UseCase[*uc_unblockip.ComandoDesbloquearIP, *uc_unblockip.RespuestaDesbloquearIP]

	// Casos de uso — sesiones
	ListarSesionesCasoDeUso     decorator.UseCase[*uc_listsessions.ComandoListarSesiones, *uc_listsessions.RespuestaListarSesiones]
	ForzarCierreSesionCasoDeUso decorator.UseCase[*uc_terminatesession.ComandoForzarCierreSesion, *uc_terminatesession.RespuestaForzarCierreSesion]

	// Casos de uso — roles y permisos
	ListarRolesCasoDeUso              decorator.UseCase[*listroles.ComandoListarRoles, *listroles.RespuestaListarRoles]
	ListarPermisosCasoDeUso           *uc_listpermisos.ListarPermisosCasoDeUso
	ListarMisPermisosCasoDeUso        *uc_listarmispermisos.ListarMisPermisosCasoDeUso
	ListarRolesDeUsuarioCasoDeUso     *uc_listarrolesdeusuario.ListarRolesDeUsuarioCasoDeUso
	CrearRolCasoDeUso            decorator.UseCase[*createrole.ComandoCrearRol, *createrole.RespuestaCrearRol]
	ModificarRolCasoDeUso        decorator.UseCase[*updaterole.ComandoModificarRol, *updaterole.RespuestaModificarRol]
	EliminarRolCasoDeUso         decorator.UseCase[*deleterole.ComandoEliminarRol, *deleterole.RespuestaEliminarRol]
	AsignarRolCasoDeUso          decorator.UseCase[*assignrole.ComandoAsignarRol, *assignrole.RespuestaAsignarRol]
	RevocarRolCasoDeUso          decorator.UseCase[*revokerole.ComandoRevocarRol, *revokerole.RespuestaRevocarRol]
	AsignarPermisoARolCasoDeUso  decorator.UseCase[*assignpermissiontorole.ComandoAsignarPermisoARol, *assignpermissiontorole.RespuestaAsignarPermisoARol]
	RevocarPermisoDeRolCasoDeUso decorator.UseCase[*revokepermissionfromrole.ComandoRevocarPermisoDeRol, *revokepermissionfromrole.RespuestaRevocarPermisoDeRol]

	// Casos de uso — verificación y recuperación
	SolicitarVerificacionCasoDeUso    decorator.UseCase[*uc_solicitar.ComandoSolicitarVerificacion, *uc_solicitar.RespuestaSolicitarVerificacion]
	ConfirmarVerificacionCasoDeUso    decorator.UseCase[*uc_confirmar.ComandoConfirmarVerificacion, *uc_confirmar.RespuestaConfirmarVerificacion]
	ReenviarVerificacionCasoDeUso     decorator.UseCase[*uc_reenviar.ComandoReenviarVerificacion, *uc_reenviar.RespuestaSolicitarVerificacion]
	SolicitarRecuperacionCasoDeUso    decorator.UseCase[*uc_solicitar_recuperacion.ComandoSolicitarRecuperacion, *uc_solicitar_recuperacion.RespuestaSolicitarRecuperacion]
	ValidarTokenRecuperacionCasoDeUso decorator.UseCase[*uc_validar_recuperacion.ComandoValidarTokenRecuperacion, *uc_validar_recuperacion.RespuestaValidarTokenRecuperacion]
	ConfirmarRecuperacionCasoDeUso    decorator.UseCase[*uc_confirmar_recuperacion.ComandoConfirmarRecuperacion, *uc_confirmar_recuperacion.RespuestaConfirmarRecuperacion]

	// Casos de uso — invitaciones
	CrearInvitacionCasoDeUso     decorator.UseCase[*uc_crearInvitacion.ComandoCrearInvitacion, *uc_crearInvitacion.RespuestaCrearInvitacion]
	AceptarInvitacionCasoDeUso   decorator.UseCase[*uc_aceptarInvitacion.ComandoAceptarInvitacion, *uc_aceptarInvitacion.RespuestaAceptarInvitacion]
	ObtenerInvitacionCasoDeUso   decorator.UseCase[*uc_obtenerInvitacion.ComandoObtenerInvitacion, *uc_obtenerInvitacion.RespuestaObtenerInvitacion]
	ListarInvitacionesCasoDeUso  decorator.UseCase[*uc_listarInvitaciones.ComandoListarInvitaciones, *uc_listarInvitaciones.RespuestaListarInvitaciones]
	ReenviarInvitacionCasoDeUso  decorator.UseCase[*uc_reenviarInvitacion.ComandoReenviarInvitacion, *uc_reenviarInvitacion.RespuestaReenviarInvitacion]
	EliminarInvitacionCasoDeUso  decorator.UseCase[*uc_eliminarInvitacion.ComandoEliminarInvitacion, *uc_eliminarInvitacion.RespuestaEliminarInvitacion]

	// Casos de uso — tenants
	ListarMisTenantsCasoDeUso     decorator.UseCase[string, *uc_listarmistenants.RespuestaListarMisTenants]
	ObtenerTenantPorIDCasoDeUso   decorator.UseCase[string, *uc_obtenertenantporid.RespuestaObtenerTenantPorID]
	ObtenerTenantPorSlugCasoDeUso decorator.UseCase[string, *uc_obtenertenantporslug.RespuestaObtenerTenantPorSlug]

	// Telemetry
	TelemetryWriter  buffer.BufferWriter
	TelemetryEnabled bool
	telemetryCancel  context.CancelFunc

	permisosConsumer *consumers.PermisosConsumer
	rolesPublisher   *publishers.RolesPublisher
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
	invitacionRepo := invitaciones_postgres.NewInvitacionRepositorio(db)

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

	// Roles Publisher
	rolesTopic := os.Getenv("KAFKA_TOPIC_ROLES")
	if rolesTopic == "" {
		env := os.Getenv("ENVIRONMENT")
		if env == "" {
			env = "dev"
		}
		rolesTopic = env + ".iam.roles"
	}
	rolesPublisher := publishers.NewRolesPublisher(strings.Split(cfg.KafkaBrokers, ","), rolesTopic)

	emailSvc := notificaciones_email.NewSMTPServicio(notificaciones_email.ConfigSMTP{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
		Async:    false,
	})

	registroUoW := usuarios_postgres.NewRegistroUnitOfWork(db)
	registroUseCase := uc_register.NewRegistrarUsuarioCasoDeUso(
		registroUoW,
		encriptacion,
		generadorID,
		rolesPublisher,
	)

	listarMisPermisosCasoDeUso := uc_listarmispermisos.NewListarMisPermisosCasoDeUso(rolRepo, rolPermisoRepo)

	// Inicialización de Telemetría
	var telemetryWriter buffer.BufferWriter
	var telemetryCancel context.CancelFunc = func() {}

	if cfg.TelemetryEnabled {
		// Configuración del buffer (valores por defecto o desde config si se expande)
		bufCfg := buffer.Config{
			Capacity:             10000,
			BatchSize:            1, // Enviamos de 1 en 1 para desarrollo
			FlushIntervalSeconds: 1,
			MaxRetries:           3,
			BackoffBase:          100 * time.Millisecond,
			BackoffMax:           2 * time.Second,
			KafkaBrokers:         strings.Split(cfg.KafkaBrokers, ","),
			KafkaTopic:           cfg.KafkaTopic,
		}
		ringBuf := buffer.NewRingBuffer(bufCfg)
		producer := buffer.NewKafkaProducer(bufCfg)

		ctx, cancel := context.WithCancel(context.Background())
		telemetryCancel = cancel
		buffer.StartConsumer(ctx, ringBuf, producer, bufCfg)

		telemetryWriter = ringBuf

		// Registrar plugin de telemetría para GORM (log_type: "BD")
		gormPlugin := gormplugin.NewTelemetryPlugin(telemetryWriter, gormplugin.DefaultConfig())
		db.Use(gormPlugin)
	} else {
		telemetryWriter = buffer.NewNoopWriter()
	}

	// Casos de uso base
	iniciarSesionUC := uc_sesiones_login.NewIniciarSesionCasoDeUso(
		sesionUoW, bloqueoIPSvc, rateLimitSvc,
		uc_sesiones_login.ConfigLogin{
			CuentaMaxIntentos:     cfg.CuentaBloqueoMaxIntentos,
			CuentaBloqueoDuracion: cfg.CuentaBloqueoDuracion,
		},
		membresiaRepo,
		usuarioTenantRolRepo,
	)

	renovarSesionUC := uc_sesiones_refresh.NewRenovarSesionCasoDeUso(
		sesionUoW,
		uc_sesiones_refresh.ConfigRefresh{
			MaxRefrescos:    cfg.SesionMaxRefrescos,
			TimeoutAbsoluto: cfg.SesionTimeoutAbsoluto,
		},
		membresiaRepo,
		usuarioTenantRolRepo,
	)

	cerrarSesionUC := uc_sesiones_logout.NewCerrarSesionCasoDeUso(sesionUoW)

	// ─────────────────────────────────────────────────────────────────────────
	// Construcción del Registry
	// ─────────────────────────────────────────────────────────────────────────

	// Casos de uso — usuarios admin
	crearUsuarioUC := uc_createuser.NewCrearUsuarioCasoDeUso(usuarioRepo, credencialesRepo, encriptacion, authSvc, generadorID)
	listarUsuariosUC := uc_listusers.NewListarUsuariosCasoDeUso(usuarioRepo, authSvc)
	modificarUsuarioUC := uc_updateuser.NewModificarUsuarioCasoDeUso(usuarioRepo, authSvc)
	darDeBajaUsuarioUC := uc_deleteuser.NewDarDeBajaUsuarioCasoDeUso(usuarioRepo, authSvc)
	expulsarUsuarioUC := uc_expeluser.NewExpulsarUsuarioCasoDeUso(membresiaRepo, usuarioTenantRolRepo, authSvc)

	// Casos de uso — autogestión
	verMiPerfilUC := uc_viewmyprofile.NewVerMiPerfilCasoDeUso(usuarioRepo)
	modificarMiPerfilUC := uc_updatemyprofile.NewModificarMiPerfilCasoDeUso(usuarioRepo)
	cambiarMiContrasenaUC := uc_changemypassword.NewCambiarMiContrasenaCasoDeUso(credencialesRepo, encriptacion)

	// Casos de uso — seguridad
	consultarCredencialesUC := uc_viewcredentials.NewConsultarCredencialesCasoDeUso(credencialesRepo, authSvc)
	resetearContrasenaUC := uc_resetpassword.NewResetearContrasenaCasoDeUso(credencialesRepo, sesionRepo, encriptacion, authSvc)
	desbloquearCuentaUC := uc_unlockaccount.NewDesbloquearCuentaCasoDeUso(credencialesRepo, authSvc)
	listarIPsBloqueadasUC := uc_listblockedips.NewListarIPsBloqueadasCasoDeUso(intentoIPRepo, authSvc)
	desbloquearIPUC := uc_unblockip.NewDesbloquearIPCasoDeUso(intentoIPRepo, authSvc)

	// Casos de uso — sesiones
	listarSesionesUC := uc_listsessions.NewListarSesionesCasoDeUso(sesionRepo, authSvc)
	forzarCierreSesionUC := uc_terminatesession.NewForzarCierreSesionCasoDeUso(sesionRepo, authSvc)

	// Casos de uso — switch tenant
	cambiarTenantUC := uc_sesiones_switchtenant.NewCambiarTenantCasoDeUso(
		membresiaRepo,
		usuarioTenantRolRepo,
		sesionUoW,
	)

	// Casos de uso — roles y permisos (solo los que tienen Ejecutar de 1 método)
	listarRolesUC := listroles.NewListarRolesCasoDeUso(rolRepo, permisoRepo, authSvc)
	listarRolesDeUsuarioUC := uc_listarrolesdeusuario.NewListarRolesDeUsuarioCasoDeUso(usuarioTenantRolRepo)
	crearRolUC := createrole.NewCrearRolCasoDeUso(rolRepo, permisoRepo, rolPermisoRepo, authSvc)
	modificarRolUC := updaterole.NewModificarRolCasoDeUso(rolRepo, authSvc)
	eliminarRolUC := deleterole.NewEliminarRolCasoDeUso(rolRepo, authSvc)
	asignarRolUC := assignrole.NewAsignarRolCasoDeUso(usuarioRolRepo, usuarioTenantRolRepo, rolRepo, authSvc)
	revocarRolUC := revokerole.NewRevocarRolCasoDeUso(usuarioRolRepo, usuarioTenantRolRepo, authSvc)
	asignarPermisoARolUC := assignpermissiontorole.NewAsignarPermisoARolCasoDeUso(rolRepo, permisoRepo, rolPermisoRepo, authSvc, rolesPublisher)
	revocarPermisoDeRolUC := revokepermissionfromrole.NewRevocarPermisoDeRolCasoDeUso(rolRepo, permisoRepo, rolPermisoRepo, authSvc, rolesPublisher)

	// Casos de uso — invitaciones
	crearInvitacionUC := uc_crearInvitacion.NewCrearInvitacionCasoDeUso(
		authSvc,
		invitacionRepo,
		tenantRepo,
		rolRepo,
		emailSvc,
		generadorID,
		cfg.FrontendURL,
		cfg.InvitacionTokenExpiracion,
	)
	aceptarInvitacionUC := uc_aceptarInvitacion.NewAceptarInvitacionCasoDeUso(
		invitacionRepo,
		membresiaRepo,
		usuarioTenantRolRepo,
		usuarioRepo,
	)
	obtenerInvitacionUC := uc_obtenerInvitacion.NewObtenerInvitacionCasoDeUso(
		invitacionRepo,
		tenantRepo,
		rolRepo,
	)

	listarInvitacionesUC := uc_listarInvitaciones.NewListarInvitacionesCasoDeUso(
		invitacionRepo,
		rolRepo,
	)

	reenviarInvitacionUC := uc_reenviarInvitacion.NewReenviarInvitacionCasoDeUso(
		invitacionRepo,
		tenantRepo,
		emailSvc,
		generadorID,
		cfg.FrontendURL,
		cfg.InvitacionTokenExpiracion,
	)

	eliminarInvitacionUC := uc_eliminarInvitacion.NewEliminarInvitacionCasoDeUso(
		invitacionRepo,
		authSvc,
	)

	// Casos de uso — tenants
	listarMisTenantsUC := uc_listarmistenants.NewListarMisTenantsCasoDeUso(tenantRepo)
	obtenerTenantPorIDUC := uc_obtenertenantporid.NewObtenerTenantPorIDCasoDeUso(tenantRepo)
	obtenerTenantPorSlugUC := uc_obtenertenantporslug.NewObtenerTenantPorSlugCasoDeUso(tenantRepo)

	// Casos de uso — recuperación de contraseña
	recuperacionConfig := uc_solicitar_recuperacion.ConfigRecuperacion{
		TokenExpiracion:     cfg.RecuperacionTokenExpiracion,
		RateLimitIPMax:      cfg.RecuperacionRateLimitIPMax,
		RateLimitUsuarioMax: cfg.RecuperacionRateLimitUsuarioMax,
		RateLimitVentana:    cfg.RecuperacionRateLimitVentana,
		FrontendURL:         cfg.FrontendURL,
	}
	solicitarRecuperacionUC := uc_solicitar_recuperacion.NewSolicitarRecuperacionCasoDeUso(
		tokenRecuperacionRepo, usuarioRecuperacionRepo, emailSvc, generadorID, recuperacionConfig,
	)
	validarTokenRecuperacionUC := uc_validar_recuperacion.NewValidarTokenRecuperacionCasoDeUso(tokenRecuperacionRepo)
	confirmarRecuperacionUC := uc_confirmar_recuperacion.NewConfirmarRecuperacionCasoDeUso(
		tokenRecuperacionRepo, usuarioRecuperacionRepo, sesionRepo, encriptacion, validarTokenRecuperacionUC,
	)

	// Casos de uso — verificación de correo
	solicitarVerificacionUC := uc_solicitar.NewSolicitarVerificacionCasoDeUso(
		verificacionRepo, emailSvc, generadorID,
		uc_solicitar.ConfigVerificacion{
			FrontendURL:     cfg.FrontendURL,
			TokenExpiracion: cfg.VerificacionTokenExpiracion,
			MaxReenvios:     cfg.VerificacionMaxReenvios,
			VentanaReenvios: cfg.VerificacionVentanaReenvios,
		},
	)
	confirmarVerificacionUC := uc_confirmar.NewConfirmarVerificacionCasoDeUso(
		verificacionRepo,
		uc_confirmar.ConfigVerificacion{
			FrontendURL:     cfg.FrontendURL,
			TokenExpiracion: cfg.VerificacionTokenExpiracion,
			MaxReenvios:     cfg.VerificacionMaxReenvios,
			VentanaReenvios: cfg.VerificacionVentanaReenvios,
		},
	)
	reenviarVerificacionUC := uc_reenviar.NewReenviarVerificacionCasoDeUso(
		verificacionRepo, emailSvc, generadorID, solicitarVerificacionUC,
		uc_reenviar.ConfigVerificacion{
			FrontendURL:     cfg.FrontendURL,
			TokenExpiracion: cfg.VerificacionTokenExpiracion,
			MaxReenvios:     cfg.VerificacionMaxReenvios,
			VentanaReenvios: cfg.VerificacionVentanaReenvios,
		},
	)

	// ─────────────────────────────────────────────────────────────────────────
	// Ensamblado del Registry
	// ─────────────────────────────────────────────────────────────────────────

	reg := &Registry{
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

		RegistrarUsuarioCasoDeUso:    registroUseCase,
		ServicioBloqueoIP: bloqueoIPSvc,
		ServicioRateLimit: rateLimitSvc,

		// Casos de uso — auth
		IniciarSesionCasoDeUso: iniciarSesionUC,
		CerrarSesionCasoDeUso:  cerrarSesionUC,
		RenovarSesionCasoDeUso: renovarSesionUC,

		// Casos de uso — usuarios admin
		CrearUsuarioCasoDeUso:     crearUsuarioUC,
		ListarUsuariosCasoDeUso:   listarUsuariosUC,
		ModificarUsuarioCasoDeUso: modificarUsuarioUC,
		DarDeBajaUsuarioCasoDeUso: darDeBajaUsuarioUC,
		ExpulsarUsuarioCasoDeUso:  expulsarUsuarioUC,

		// Casos de uso — autogestión
		VerMiPerfilCasoDeUso:         verMiPerfilUC,
		ModificarMiPerfilCasoDeUso:   modificarMiPerfilUC,
		CambiarMiContrasenaCasoDeUso: cambiarMiContrasenaUC,

		// Casos de uso — seguridad
		ConsultarCredencialesCasoDeUso: consultarCredencialesUC,
		ResetearContrasenaCasoDeUso:    resetearContrasenaUC,
		DesbloquearCuentaCasoDeUso:     desbloquearCuentaUC,
		ListarIPsBloqueadasCasoDeUso:   listarIPsBloqueadasUC,
		DesbloquearIPCasoDeUso:         desbloquearIPUC,

		// Casos de uso — sesiones
		ListarSesionesCasoDeUso:     listarSesionesUC,
		ForzarCierreSesionCasoDeUso: forzarCierreSesionUC,

		// Casos de uso — switch tenant
		CambiarTenantCasoDeUso: cambiarTenantUC,

		// Casos de uso — verificación
		SolicitarVerificacionCasoDeUso: solicitarVerificacionUC,
		ConfirmarVerificacionCasoDeUso: confirmarVerificacionUC,
		ReenviarVerificacionCasoDeUso:  reenviarVerificacionUC,

		// Casos de uso — recuperación
		SolicitarRecuperacionCasoDeUso:    solicitarRecuperacionUC,
		ValidarTokenRecuperacionCasoDeUso: validarTokenRecuperacionUC,
		ConfirmarRecuperacionCasoDeUso:    confirmarRecuperacionUC,

		// Casos de uso — roles y permisos
		ListarRolesCasoDeUso:              listarRolesUC,
		ListarPermisosCasoDeUso:           uc_listpermisos.NewListarPermisosCasoDeUso(permisoRepo, authSvc),
		ListarMisPermisosCasoDeUso:        listarMisPermisosCasoDeUso,
		ListarRolesDeUsuarioCasoDeUso:     listarRolesDeUsuarioUC,
		CrearRolCasoDeUso:            crearRolUC,
		ModificarRolCasoDeUso:        modificarRolUC,
		EliminarRolCasoDeUso:         eliminarRolUC,
		AsignarRolCasoDeUso:          asignarRolUC,
		RevocarRolCasoDeUso:          revocarRolUC,
		AsignarPermisoARolCasoDeUso:  asignarPermisoARolUC,
		RevocarPermisoDeRolCasoDeUso: revocarPermisoDeRolUC,

		// Casos de uso — invitaciones
		CrearInvitacionCasoDeUso:     crearInvitacionUC,
		AceptarInvitacionCasoDeUso:   aceptarInvitacionUC,
		ObtenerInvitacionCasoDeUso:   obtenerInvitacionUC,
		ListarInvitacionesCasoDeUso:  listarInvitacionesUC,
		ReenviarInvitacionCasoDeUso:  reenviarInvitacionUC,
		EliminarInvitacionCasoDeUso:  eliminarInvitacionUC,

		// Casos de uso — tenants
		ListarMisTenantsCasoDeUso:     listarMisTenantsUC,
		ObtenerTenantPorIDCasoDeUso:   obtenerTenantPorIDUC,
		ObtenerTenantPorSlugCasoDeUso: obtenerTenantPorSlugUC,

		// Telemetry fields
		TelemetryWriter:  telemetryWriter,
		TelemetryEnabled: cfg.TelemetryEnabled,
		telemetryCancel:  telemetryCancel,

		rolesPublisher: rolesPublisher,
	}

	// ─────────────────────────────────────────────────────────────────────────
	// Decoradores de Telemetría
	// ─────────────────────────────────────────────────────────────────────────

	if cfg.TelemetryEnabled {
		// Auth (login tiene UseCase, logout tiene custom 2-métodos)
		reg.IniciarSesionCasoDeUso = decorator.Wrap("Login", telemetryWriter, iniciarSesionUC)
		reg.RenovarSesionCasoDeUso = decorator.Wrap("Refresh", telemetryWriter, renovarSesionUC)
		reg.CerrarSesionCasoDeUso = decorator.NewDecoratorLogout(cerrarSesionUC, telemetryWriter)
		reg.RegistrarUsuarioCasoDeUso = decorator.Wrap("Registro", telemetryWriter, registroUseCase)

		// Usuarios admin
		reg.CrearUsuarioCasoDeUso = decorator.Wrap("CrearUsuario", telemetryWriter, crearUsuarioUC)
		reg.ListarUsuariosCasoDeUso = decorator.Wrap("ListarUsuarios", telemetryWriter, listarUsuariosUC)
		reg.ModificarUsuarioCasoDeUso = decorator.Wrap("ModificarUsuario", telemetryWriter, modificarUsuarioUC)
		reg.DarDeBajaUsuarioCasoDeUso = decorator.Wrap("DarDeBajaUsuario", telemetryWriter, darDeBajaUsuarioUC)
		reg.ExpulsarUsuarioCasoDeUso = decorator.Wrap("ExpulsarUsuario", telemetryWriter, expulsarUsuarioUC)

		// Autogestión
		reg.VerMiPerfilCasoDeUso = decorator.Wrap("VerMiPerfil", telemetryWriter, verMiPerfilUC)
		reg.ModificarMiPerfilCasoDeUso = decorator.Wrap("ModificarMiPerfil", telemetryWriter, modificarMiPerfilUC)
		reg.CambiarMiContrasenaCasoDeUso = decorator.Wrap("CambiarMiContrasena", telemetryWriter, cambiarMiContrasenaUC)

		// Seguridad
		reg.ConsultarCredencialesCasoDeUso = decorator.Wrap("ConsultarCredenciales", telemetryWriter, consultarCredencialesUC)
		reg.ResetearContrasenaCasoDeUso = decorator.Wrap("ResetearContrasena", telemetryWriter, resetearContrasenaUC)
		reg.DesbloquearCuentaCasoDeUso = decorator.Wrap("DesbloquearCuenta", telemetryWriter, desbloquearCuentaUC)
		reg.ListarIPsBloqueadasCasoDeUso = decorator.Wrap("ListarIPsBloqueadas", telemetryWriter, listarIPsBloqueadasUC)
		reg.DesbloquearIPCasoDeUso = decorator.Wrap("DesbloquearIP", telemetryWriter, desbloquearIPUC)

		// Sesiones
		reg.ListarSesionesCasoDeUso = decorator.Wrap("ListarSesiones", telemetryWriter, listarSesionesUC)
		reg.ForzarCierreSesionCasoDeUso = decorator.Wrap("ForzarCierreSesion", telemetryWriter, forzarCierreSesionUC)

		// Switch tenant
		reg.CambiarTenantCasoDeUso = decorator.Wrap("CambiarTenant", telemetryWriter, cambiarTenantUC)

		// Roles
		reg.ListarRolesCasoDeUso = decorator.Wrap("ListarRoles", telemetryWriter, listarRolesUC)
		reg.CrearRolCasoDeUso = decorator.Wrap("CrearRol", telemetryWriter, crearRolUC)
		reg.ModificarRolCasoDeUso = decorator.Wrap("ModificarRol", telemetryWriter, modificarRolUC)
		reg.EliminarRolCasoDeUso = decorator.Wrap("EliminarRol", telemetryWriter, eliminarRolUC)
		reg.AsignarRolCasoDeUso = decorator.Wrap("AsignarRol", telemetryWriter, asignarRolUC)
		reg.RevocarRolCasoDeUso = decorator.Wrap("RevocarRol", telemetryWriter, revocarRolUC)
		reg.AsignarPermisoARolCasoDeUso = decorator.Wrap("AsignarPermisoARol", telemetryWriter, asignarPermisoARolUC)
		reg.RevocarPermisoDeRolCasoDeUso = decorator.Wrap("RevocarPermisoDeRol", telemetryWriter, revocarPermisoDeRolUC)

		// Verificación
		reg.SolicitarVerificacionCasoDeUso = decorator.Wrap("SolicitarVerificacion", telemetryWriter, solicitarVerificacionUC)
		reg.ConfirmarVerificacionCasoDeUso = decorator.Wrap("ConfirmarVerificacion", telemetryWriter, confirmarVerificacionUC)
		reg.ReenviarVerificacionCasoDeUso = decorator.Wrap("ReenviarVerificacion", telemetryWriter, reenviarVerificacionUC)

		// Recuperación
		reg.SolicitarRecuperacionCasoDeUso = decorator.Wrap("SolicitarRecuperacion", telemetryWriter, solicitarRecuperacionUC)
		reg.ValidarTokenRecuperacionCasoDeUso = decorator.Wrap("ValidarTokenRecuperacion", telemetryWriter, validarTokenRecuperacionUC)
		reg.ConfirmarRecuperacionCasoDeUso = decorator.Wrap("ConfirmarRecuperacion", telemetryWriter, confirmarRecuperacionUC)

		// Invitaciones
		reg.CrearInvitacionCasoDeUso = decorator.Wrap("CrearInvitacion", telemetryWriter, crearInvitacionUC)
		reg.AceptarInvitacionCasoDeUso = decorator.Wrap("AceptarInvitacion", telemetryWriter, aceptarInvitacionUC)
		reg.ObtenerInvitacionCasoDeUso = decorator.Wrap("ObtenerInvitacion", telemetryWriter, obtenerInvitacionUC)
		reg.ListarInvitacionesCasoDeUso = decorator.Wrap("ListarInvitaciones", telemetryWriter, listarInvitacionesUC)
		reg.ReenviarInvitacionCasoDeUso = decorator.Wrap("ReenviarInvitacion", telemetryWriter, reenviarInvitacionUC)
		reg.EliminarInvitacionCasoDeUso = decorator.Wrap("EliminarInvitacion", telemetryWriter, eliminarInvitacionUC)

		// Tenants
		reg.ListarMisTenantsCasoDeUso = decorator.Wrap("ListarMisTenants", telemetryWriter, listarMisTenantsUC)
		reg.ObtenerTenantPorIDCasoDeUso = decorator.Wrap("ObtenerTenantPorID", telemetryWriter, obtenerTenantPorIDUC)
		reg.ObtenerTenantPorSlugCasoDeUso = decorator.Wrap("ObtenerTenantPorSlug", telemetryWriter, obtenerTenantPorSlugUC)
	}

	// ─────────────────────────────────────────────────────────────────────────
	// Consumidores Kafka
	// ─────────────────────────────────────────────────────────────────────────
	topicPermisos := cfg.KafkaTopicPermisos
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	permisosConsumer := consumers.NewPermisosConsumer(brokers, topicPermisos, permisoRepo, rolRepo, rolPermisoRepo, rolesPublisher, tenantRepo, generadorID)
	reg.permisosConsumer = permisosConsumer
	go permisosConsumer.Start(context.Background())

	return reg
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
func (r *Registry) GetRegistrarUsuarioCasoDeUso() decorator.UseCase[*uc_register.ComandoRegistrarUsuario, *uc_register.RespuestaRegistrarUsuario] {
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

// Close libera los recursos del Registry, incluyendo la telemetría.
func (r *Registry) Close() {
	if r.telemetryCancel != nil {
		r.telemetryCancel()
	}
	if r.permisosConsumer != nil {
		r.permisosConsumer.Close()
	}
	if r.rolesPublisher != nil {
		r.rolesPublisher.Close()
	}
}

package registry

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	fincasdomain "github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	fincaspostgres "github.com/davosjar/bunna/services/fincas/internal/fincas/infrastructure/persistence/postgres"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	diagnosticopostgres "github.com/davosjar/bunna/services/fincas/internal/diagnostico/infrastructure/persistence/postgres"
	nodosdomain "github.com/davosjar/bunna/services/fincas/internal/nodos/domain"
	nodospostgres "github.com/davosjar/bunna/services/fincas/internal/nodos/infrastructure/persistence/postgres"
	iampostgres "github.com/davosjar/bunna/services/fincas/internal/infrastructure/security/iam/postgres"
	iamconsumers "github.com/davosjar/bunna/services/fincas/internal/infrastructure/security/iam/consumers"
	"github.com/davosjar/bunna/services/fincas/internal/application"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"github.com/davosjar/bunna/services/fincas/internal/shared/infrastructure/idgenerator"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/eventpublisher"
	telemetryinfra "github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/decorator"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/gormplugin"
	telemetrymiddleware "github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/middleware"

	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarfinca"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/desactivarfinca"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/agregarlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/eliminarlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/tomarmuestra"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/listarmuestrasporlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/solicitardiagnosticomanual"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarinferencia"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/aceptardiagnostico"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/rechazardiagnostico"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/generarreporteporlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/listarnodos"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/obtenernodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/editarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/desactivarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/validarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarinferenciadesdenodo"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/facades"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/handler"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/router"
	jwtvalidator "github.com/davosjar/bunna/services/fincas/internal/infrastructure/security/jwt"
	fincasmiddleware "github.com/davosjar/bunna/services/fincas/internal/presentation/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Registry centraliza la creación y el ciclo de vida de todas las dependencias
// del microservicio de fincas: configuración, base de datos, repositorios, casos de uso, etc.
type Registry struct {
	db          *gorm.DB
	generadorID shared.GeneradorID
	publisher   application.EventPublisher
	serverPort  string

	// Repositorios (privados)
	fincaRepo       fincasdomain.FincaRepositorio
	loteRepo        fincasdomain.LoteRepositorio
	muestraRepo     diagnosticodomain.MuestraRepositorio
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio
	candidatoRepo   diagnosticodomain.CandidatoReentrenamientoRepositorio
	nodoRepo        nodosdomain.NodoRepositorio
	iamRepo         *iampostgres.IAMRepositorio

	// Unit of Work (privados)
	fincaUoW       *fincaspostgres.UnitOfWorkPostgres
	diagnosticoUoW *diagnosticopostgres.UnitOfWorkDiagnostico

	// Servicios de dominio (privados)
	fincaService *fincasdomain.FincaService

	// Casos de uso — fincas (públicos, ya decorados si telemetría activa)
	RegistrarFinca  facades.RegistrarFincaUseCase
	DesactivarFinca facades.DesactivarFincaUseCase

	// Casos de uso — lotes (públicos)
	AgregarLote  facades.AgregarLoteUseCase
	EliminarLote facades.EliminarLoteUseCase

	// Casos de uso — muestras (públicos)
	TomarMuestra          facades.TomarMuestraUseCase
	ListarMuestrasPorLote facades.ListarMuestrasPorLoteUseCase

	// Casos de uso — diagnósticos (públicos)
	SolicitarDiagnosticoManual facades.SolicitarDiagnosticoManualUseCase
	RegistrarInferencia        decorator.UseCase[registrarinferencia.Command, *registrarinferencia.Salida]
	AceptarDiagnostico         facades.AceptarDiagnosticoUseCase
	RechazarDiagnostico        facades.RechazarDiagnosticoUseCase

	// Casos de uso — reportes (públicos)
	GenerarReportePorLote facades.GenerarReportePorLoteUseCase

	// Facades (privados)
	fincasFacade       facades.FincasFacade
	lotesFacade        facades.LotesFacade
	muestrasFacade     facades.MuestrasFacade
	diagnosticosFacade facades.DiagnosticosFacade
	reportesFacade     facades.ReportesFacade

	// Seguridad (privados)
	tokenValidator *jwtvalidator.TokenValidator
	authMiddleware *fincasmiddleware.AuthMiddleware

	// Handlers (privados)
	fincaHandler       *handler.FincaHandler
	loteHandler        *handler.LoteHandler
	muestraHandler     *handler.MuestraHandler
	diagnosticoHandler *handler.DiagnosticoHandler
	reporteHandler     *handler.ReporteHandler
	nodoHandler        *handler.NodoHandler

	// Router (privado)
	router *gin.Engine

	// Telemetría
	TelemetryWriter  buffer.BufferWriter
	TelemetryEnabled bool
	telemetryCancel  context.CancelFunc

	rolesConsumer *iamconsumers.RolesConsumer
}

// NewRegistry crea todas las dependencias, ejecuta auto-migrate y devuelve un Registry listo para usar.
func NewRegistry() *Registry {
	cfg := loadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.dbHost, cfg.dbPort, cfg.dbUser, cfg.dbPassword, cfg.dbName, cfg.dbSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error conectando a BD: %v", err)
	}

	serviceInfo := telemetryinfra.ServiceInfo{
		Name:        cfg.serviceName,
		Environment: cfg.environment,
	}

	var telemetryWriter buffer.BufferWriter
	var telemetryCancel context.CancelFunc = func() {}

	if cfg.telemetryEnabled {
		bufCfg := buffer.Config{
			Capacity:             10000,
			BatchSize:            1,
			FlushIntervalSeconds: 1,
			MaxRetries:           3,
			BackoffBase:          100 * time.Millisecond,
			BackoffMax:           2 * time.Second,
			KafkaBrokers:         strings.Split(cfg.kafkaBrokers, ","),
			KafkaTopic:           cfg.kafkaTopic,
		}
		ringBuf := buffer.NewRingBuffer(bufCfg)
		producer := buffer.NewKafkaProducer(bufCfg)

		ctx, cancel := context.WithCancel(context.Background())
		telemetryCancel = cancel
		buffer.StartConsumer(ctx, ringBuf, producer, bufCfg)

		telemetryWriter = ringBuf

		gormPlugin := gormplugin.NewTelemetryPlugin(telemetryWriter, gormplugin.DefaultConfig())
		if err := db.Use(gormPlugin); err != nil {
			log.Fatalf("Error registrando plugin GORM de telemetría: %v", err)
		}
	} else {
		telemetryWriter = buffer.NewNoopWriter()
	}

	// AutoMigrate
	if err := runAutoMigrate(db); err != nil {
		log.Fatalf("Error en auto-migrate: %v", err)
	}
	log.Println("AutoMigrate completado")

	generador := idgenerator.NewGeneradorUUIDV7()

	// publisher: si hay brokers Kafka configurados, usa Kafka; si no, consola
	var publisher application.EventPublisher
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC_PERMISOS")
	if kafkaTopic == "" {
		kafkaTopic = cfg.environment + ".permisos"
	}
	if kafkaBrokers != "" {
		brokers := strings.Split(kafkaBrokers, ",")
		kp := eventpublisher.NewKafkaPublisher(brokers, kafkaTopic)
		publisher = kp
		log.Printf("EventPublisher: Kafka (%s → topic '%s')", kafkaBrokers, kafkaTopic)
	} else {
		publisher = eventpublisher.NewConsolePublisher()
		log.Println("EventPublisher: consola (sin Kafka)")
	}

	// Repositorios
	fincaRepo := fincaspostgres.NewFincaRepositorio(db)
	loteRepo := fincaspostgres.NewLoteRepositorio(db)
	muestraRepo := diagnosticopostgres.NewMuestraRepositorio(db)
	diagnosticoRepo := diagnosticopostgres.NewDiagnosticoRepositorio(db)
	candidatoRepo := diagnosticopostgres.NewCandidatoReentrenamientoRepositorio(db)
	nodoRepo := nodospostgres.NewNodoRepositorio(db)
	iamRepo := iampostgres.NewIAMRepositorio(db)

	// Unit of Work
	fincaUoW := fincaspostgres.NewUnitOfWorkPostgres(db, generador)
	diagnosticoUoW := diagnosticopostgres.NewUnitOfWorkDiagnostico(db, generador)

	// Servicios de dominio
	fincaService := fincasdomain.NewFincaService()

	// Casos de uso — instancia única por UC
	registrarFincaInner := registrarfinca.NewUseCase(fincaRepo, generador, publisher)
	desactivarFincaInner := desactivarfinca.NewUseCase(fincaRepo, fincaService, generador, publisher)
	agregarLoteInner := agregarlote.NewUseCase(fincaRepo, loteRepo, generador, publisher)
	eliminarLoteInner := eliminarlote.NewUseCase(loteRepo, generador, publisher)
	tomarMuestraInner := tomarmuestra.NewUseCase(loteRepo, muestraRepo, generador, publisher)
	listarMuestrasInner := listarmuestrasporlote.NewUseCase(loteRepo, muestraRepo)
	solicitarDiagnosticoInner := solicitardiagnosticomanual.NewUseCase(muestraRepo, generador, publisher)
	registrarInferenciaInner := registrarinferencia.NewUseCase(muestraRepo, diagnosticoRepo, generador, publisher)
	aceptarDiagnosticoInner := aceptardiagnostico.NewUseCase(diagnosticoRepo, generador, publisher)
	rechazarDiagnosticoInner := rechazardiagnostico.NewUseCase(diagnosticoRepo, candidatoRepo, diagnosticoUoW, generador, publisher)
	generarReporteInner := generarreporteporlote.NewUseCase(loteRepo, muestraRepo, diagnosticoRepo)

	// Casos de uso — nodos
	registrarNodoInner := registrarnodo.NewUseCase(nodoRepo, fincaRepo, generador)
	listarNodosInner := listarnodos.NewUseCase(nodoRepo)
	obtenerNodoInner := obtenernodo.NewUseCase(nodoRepo)
	editarNodoInner := editarnodo.NewUseCase(nodoRepo)
	desactivarNodoInner := desactivarnodo.NewUseCase(nodoRepo)
	validarNodoInner := validarnodo.NewUseCase(nodoRepo)
	registrarInferenciaDesdeNodoInner := registrarinferenciadesdenodo.NewUseCase(nodoRepo, loteRepo, diagnosticoUoW, generador, publisher)

	var registrarFincaUC facades.RegistrarFincaUseCase = registrarFincaInner
	var desactivarFincaUC facades.DesactivarFincaUseCase = desactivarFincaInner
	var agregarLoteUC facades.AgregarLoteUseCase = agregarLoteInner
	var eliminarLoteUC facades.EliminarLoteUseCase = eliminarLoteInner
	var tomarMuestraUC facades.TomarMuestraUseCase = tomarMuestraInner
	var listarMuestrasUC facades.ListarMuestrasPorLoteUseCase = listarMuestrasInner
	var solicitarDiagnosticoUC facades.SolicitarDiagnosticoManualUseCase = solicitarDiagnosticoInner
	var registrarInferenciaUC decorator.UseCase[registrarinferencia.Command, *registrarinferencia.Salida] = registrarInferenciaInner
	var aceptarDiagnosticoUC facades.AceptarDiagnosticoUseCase = aceptarDiagnosticoInner
	var rechazarDiagnosticoUC facades.RechazarDiagnosticoUseCase = rechazarDiagnosticoInner
	var generarReporteUC facades.GenerarReportePorLoteUseCase = generarReporteInner

	// Casos de uso — nodos (variables)
	var registrarNodoUC facades.RegistrarNodoUseCase = registrarNodoInner
	var listarNodosUC facades.ListarNodosUseCase = listarNodosInner
	var obtenerNodoUC facades.ObtenerNodoUseCase = obtenerNodoInner
	var editarNodoUC facades.EditarNodoUseCase = editarNodoInner
	var desactivarNodoUC facades.DesactivarNodoUseCase = desactivarNodoInner
	var validarNodoUC facades.ValidarNodoUseCase = validarNodoInner
	var inferenciaNodoUC facades.RegistrarInferenciaDesdeNodoUseCase = registrarInferenciaDesdeNodoInner

	// Decoradores de telemetría — capa NEGOCIO (APO)
	if cfg.telemetryEnabled {
		registrarFincaUC = decorator.WrapAuth("RegistrarFinca", telemetryWriter, serviceInfo, registrarFincaInner)
		desactivarFincaUC = decorator.WrapAuth("DesactivarFinca", telemetryWriter, serviceInfo, desactivarFincaInner)
		agregarLoteUC = decorator.WrapAuth("AgregarLote", telemetryWriter, serviceInfo, agregarLoteInner)
		eliminarLoteUC = decorator.WrapAuth("EliminarLote", telemetryWriter, serviceInfo, eliminarLoteInner)
		tomarMuestraUC = decorator.WrapAuth("TomarMuestra", telemetryWriter, serviceInfo, tomarMuestraInner)
		listarMuestrasUC = decorator.WrapAuth("ListarMuestrasPorLote", telemetryWriter, serviceInfo, listarMuestrasInner)
		solicitarDiagnosticoUC = decorator.WrapAuth("SolicitarDiagnosticoManual", telemetryWriter, serviceInfo, solicitarDiagnosticoInner)
		registrarInferenciaUC = decorator.Wrap("RegistrarInferencia", telemetryWriter, serviceInfo, registrarInferenciaInner)
		aceptarDiagnosticoUC = decorator.WrapAuth("AceptarDiagnostico", telemetryWriter, serviceInfo, aceptarDiagnosticoInner)
		rechazarDiagnosticoUC = decorator.WrapAuth("RechazarDiagnostico", telemetryWriter, serviceInfo, rechazarDiagnosticoInner)
		generarReporteUC = decorator.WrapAuth("GenerarReportePorLote", telemetryWriter, serviceInfo, generarReporteInner)
		registrarNodoUC = decorator.WrapAuth("RegistrarNodo", telemetryWriter, serviceInfo, registrarNodoInner)
		obtenerNodoUC = decorator.WrapAuth("ObtenerNodo", telemetryWriter, serviceInfo, obtenerNodoInner)
		editarNodoUC = decorator.WrapAuth("EditarNodo", telemetryWriter, serviceInfo, editarNodoInner)
		desactivarNodoUC = decorator.WrapAuth("DesactivarNodo", telemetryWriter, serviceInfo, desactivarNodoInner)
		validarNodoUC = decorator.Wrap("ValidarNodo", telemetryWriter, serviceInfo, validarNodoInner)
		inferenciaNodoUC = decorator.Wrap("RegistrarInferenciaDesdeNodo", telemetryWriter, serviceInfo, registrarInferenciaDesdeNodoInner)
	}

	// Facades (reciben los mismos UC decorados)
	fincasFacade := facades.NewFincasFacade(registrarFincaUC, desactivarFincaUC)
	lotesFacade := facades.NewLotesFacade(agregarLoteUC, eliminarLoteUC)
	muestrasFacade := facades.NewMuestrasFacade(tomarMuestraUC, listarMuestrasUC)
	diagnosticosFacade := facades.NewDiagnosticosFacade(solicitarDiagnosticoUC, aceptarDiagnosticoUC, rechazarDiagnosticoUC)
	reportesFacade := facades.NewReportesFacade(generarReporteUC)
	nodosFacade := facades.NewNodosFacade(registrarNodoUC, listarNodosUC, obtenerNodoUC, editarNodoUC, desactivarNodoUC, validarNodoUC, inferenciaNodoUC)

	// Token validator
	tokenValidator := jwtvalidator.NewTokenValidator(jwtvalidator.Config{
		Secret: cfg.jwtSecret,
		Issuer: cfg.jwtIssuer,
	})

	authMiddleware := fincasmiddleware.NewAuthMiddleware(tokenValidator, iamRepo)

	fincaHandler := handler.NewFincaHandler(fincasFacade)
	loteHandler := handler.NewLoteHandler(lotesFacade)
	muestraHandler := handler.NewMuestraHandler(muestrasFacade)
	diagnosticoHandler := handler.NewDiagnosticoHandler(diagnosticosFacade)
	reporteHandler := handler.NewReporteHandler(reportesFacade)
	nodoHandler := handler.NewNodoHandler(nodosFacade)

	ginEngine := router.New(router.Config{
		TelemetryEnabled: cfg.telemetryEnabled,
		TelemetryWriter:  telemetryWriter,
		TelemetryCfg: telemetrymiddleware.Config{
			MaxDurationWarning: 500 * time.Millisecond,
			MaxDurationError:   1000 * time.Millisecond,
			Service:            serviceInfo,
		},
		AuthMiddleware:     authMiddleware,
		FincaHandler:       fincaHandler,
		LoteHandler:        loteHandler,
		MuestraHandler:     muestraHandler,
		DiagnosticoHandler: diagnosticoHandler,
		ReporteHandler:     reporteHandler,
		NodoHandler:        nodoHandler,
	})

	// Publicar catálogo de permisos al iniciar
	go publicarCatalogoPermisos(context.Background(), publisher, generador)

	// Iniciar consumidor de roles
	var rolesConsumer *iamconsumers.RolesConsumer
	if kafkaBrokers != "" {
		topicRoles := os.Getenv("KAFKA_TOPIC_ROLES")
		if topicRoles == "" {
			topicRoles = "dev.iam.roles"
		}
		brokers := strings.Split(kafkaBrokers, ",")
		rolesConsumer = iamconsumers.NewRolesConsumer(brokers, topicRoles, iamRepo)
		go rolesConsumer.Start(context.Background())
	}

	return &Registry{
		db:             db,
		generadorID:    generador,
		publisher:      publisher,
		serverPort:     cfg.serverPort,
		fincaRepo:      fincaRepo,
		loteRepo:       loteRepo,
		muestraRepo:    muestraRepo,
		diagnosticoRepo: diagnosticoRepo,
		candidatoRepo:  candidatoRepo,
		nodoRepo:       nodoRepo,
		iamRepo:        iamRepo,
		fincaUoW:       fincaUoW,
		diagnosticoUoW: diagnosticoUoW,
		fincaService:   fincaService,

		RegistrarFinca:             registrarFincaUC,
		DesactivarFinca:            desactivarFincaUC,
		AgregarLote:                agregarLoteUC,
		EliminarLote:               eliminarLoteUC,
		TomarMuestra:               tomarMuestraUC,
		ListarMuestrasPorLote:      listarMuestrasUC,
		SolicitarDiagnosticoManual: solicitarDiagnosticoUC,
		RegistrarInferencia:        registrarInferenciaUC,
		AceptarDiagnostico:         aceptarDiagnosticoUC,
		RechazarDiagnostico:        rechazarDiagnosticoUC,
		GenerarReportePorLote:      generarReporteUC,

		fincasFacade:       fincasFacade,
		lotesFacade:        lotesFacade,
		muestrasFacade:     muestrasFacade,
		diagnosticosFacade: diagnosticosFacade,
		reportesFacade:     reportesFacade,

		tokenValidator: tokenValidator,
		authMiddleware: authMiddleware,

		fincaHandler:       fincaHandler,
		loteHandler:        loteHandler,
		muestraHandler:     muestraHandler,
		diagnosticoHandler: diagnosticoHandler,
		reporteHandler:     reporteHandler,
		nodoHandler:        nodoHandler,

		router:             ginEngine,
		TelemetryWriter:  telemetryWriter,
		TelemetryEnabled: cfg.telemetryEnabled,
		telemetryCancel:  telemetryCancel,
		rolesConsumer:    rolesConsumer,
	}
}

// --- Getters ---

func (r *Registry) FincaRepository() fincasdomain.FincaRepositorio {
	return r.fincaRepo
}

func (r *Registry) LoteRepository() fincasdomain.LoteRepositorio {
	return r.loteRepo
}

func (r *Registry) MuestraRepository() diagnosticodomain.MuestraRepositorio {
	return r.muestraRepo
}

func (r *Registry) DiagnosticoRepository() diagnosticodomain.DiagnosticoRepositorio {
	return r.diagnosticoRepo
}

func (r *Registry) CandidatoRepository() diagnosticodomain.CandidatoReentrenamientoRepositorio {
	return r.candidatoRepo
}

func (r *Registry) FincaService() *fincasdomain.FincaService {
	return r.fincaService
}

func (r *Registry) GeneradorID() shared.GeneradorID {
	return r.generadorID
}

func (r *Registry) EventPublisher() application.EventPublisher {
	return r.publisher
}

func (r *Registry) DB() *gorm.DB {
	return r.db
}

func (r *Registry) FincaUnitOfWork() *fincaspostgres.UnitOfWorkPostgres {
	return r.fincaUoW
}

func (r *Registry) DiagnosticoUnitOfWork() *diagnosticopostgres.UnitOfWorkDiagnostico {
	return r.diagnosticoUoW
}

func (r *Registry) ServerPort() string { return r.serverPort }

// Close libera recursos: telemetría y conexión a BD.
func (r *Registry) Close() {
	if r.telemetryCancel != nil {
		r.telemetryCancel()
	}
	if r.rolesConsumer != nil {
		r.rolesConsumer.Close()
	}
	sqlDB, err := r.db.DB()
	if err == nil {
		sqlDB.Close()
	}
	log.Println("Recursos liberados")
}

func (r *Registry) FincasFacade() facades.FincasFacade            { return r.fincasFacade }
func (r *Registry) LotesFacade() facades.LotesFacade               { return r.lotesFacade }
func (r *Registry) MuestrasFacade() facades.MuestrasFacade         { return r.muestrasFacade }
func (r *Registry) DiagnosticosFacade() facades.DiagnosticosFacade { return r.diagnosticosFacade }
func (r *Registry) ReportesFacade() facades.ReportesFacade         { return r.reportesFacade }
func (r *Registry) TokenValidator() *jwtvalidator.TokenValidator   { return r.tokenValidator }
func (r *Registry) AuthMiddleware() *fincasmiddleware.AuthMiddleware { return r.authMiddleware }
func (r *Registry) FincaHandler() *handler.FincaHandler             { return r.fincaHandler }
func (r *Registry) LoteHandler() *handler.LoteHandler               { return r.loteHandler }
func (r *Registry) MuestraHandler() *handler.MuestraHandler         { return r.muestraHandler }
func (r *Registry) DiagnosticoHandler() *handler.DiagnosticoHandler { return r.diagnosticoHandler }
func (r *Registry) ReporteHandler() *handler.ReporteHandler         { return r.reporteHandler }
func (r *Registry) NodoHandler() *handler.NodoHandler               { return r.nodoHandler }
func (r *Registry) NodoRepository() nodosdomain.NodoRepositorio     { return r.nodoRepo }
func (r *Registry) Router() *gin.Engine                             { return r.router }

type config struct {
	dbHost            string
	dbPort            string
	dbUser            string
	dbPassword        string
	dbName            string
	dbSSLMode         string
	jwtSecret         string
	jwtIssuer         string
	serverPort        string
	telemetryEnabled  bool
	kafkaBrokers      string
	kafkaTopic        string
	serviceName       string
	environment       string
}

func loadConfig() config {
	getEnv := func(key, fallback string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return fallback
	}
	getEnvBool := func(key string, fallback bool) bool {
		v := os.Getenv(key)
		if v == "" {
			return fallback
		}
		return v == "true" || v == "1"
	}
	return config{
		dbHost:           getEnv("DB_HOST", "localhost"),
		dbPort:           getEnv("DB_PORT", "5432"),
		dbUser:           getEnv("DB_USER", "fincas_user"),
		dbPassword:       getEnv("DB_PASSWORD", "fincas_pass_dev"),
		dbName:           getEnv("DB_NAME", "bunna_fincas"),
		dbSSLMode:        getEnv("DB_SSLMODE", "disable"),
		jwtSecret:        getEnv("JWT_SECRET", "clave-secreta-dev"),
		jwtIssuer:        getEnv("JWT_ISSUER", ""),
		serverPort:       getEnv("SERVER_PORT", "8082"),
		telemetryEnabled: getEnvBool("TELEMETRY_ENABLED", false),
		kafkaBrokers:     getEnv("KAFKA_BROKERS", "localhost:9092"),
		kafkaTopic:       getEnv("KAFKA_TOPIC", "telemetry"),
		serviceName:      getEnv("SERVICE_NAME", "fincas"),
		environment:      getEnv("ENVIRONMENT", "dev"),
	}
}

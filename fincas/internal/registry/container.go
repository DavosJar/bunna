package registry

import (
	"fmt"
	"log"
	"os"

	fincasdomain "github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	fincaspostgres "github.com/davosjar/bunna/services/fincas/internal/fincas/infrastructure/persistence/postgres"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	diagnosticopostgres "github.com/davosjar/bunna/services/fincas/internal/diagnostico/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/fincas/internal/application"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"github.com/davosjar/bunna/services/fincas/internal/shared/infrastructure/idgenerator"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/eventpublisher"

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
	// Infraestructura base (privados)
	db          *gorm.DB
	generadorID shared.GeneradorID
	publisher   application.EventPublisher

	// Repositorios (privados)
	fincaRepo       fincasdomain.FincaRepositorio
	loteRepo        fincasdomain.LoteRepositorio
	muestraRepo     diagnosticodomain.MuestraRepositorio
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio
	candidatoRepo   diagnosticodomain.CandidatoReentrenamientoRepositorio

	// Unit of Work (privados)
	fincaUoW       *fincaspostgres.UnitOfWorkPostgres
	diagnosticoUoW *diagnosticopostgres.UnitOfWorkDiagnostico

	// Servicios de dominio (privados)
	fincaService *fincasdomain.FincaService

	// Casos de uso — fincas (públicos)
	RegistrarFinca  *registrarfinca.UseCase
	DesactivarFinca *desactivarfinca.UseCase

	// Casos de uso — lotes (públicos)
	AgregarLote  *agregarlote.UseCase
	EliminarLote *eliminarlote.UseCase

	// Casos de uso — muestras (públicos)
	TomarMuestra          *tomarmuestra.UseCase
	ListarMuestrasPorLote *listarmuestrasporlote.UseCase

	// Casos de uso — diagnósticos (públicos)
	SolicitarDiagnosticoManual *solicitardiagnosticomanual.UseCase
	RegistrarInferencia        *registrarinferencia.UseCase
	AceptarDiagnostico         *aceptardiagnostico.UseCase
	RechazarDiagnostico        *rechazardiagnostico.UseCase

	// Casos de uso — reportes (públicos)
	GenerarReportePorLote *generarreporteporlote.UseCase

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

	// Router (privado)
	router *gin.Engine
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

	// AutoMigrate
	if err := runAutoMigrate(db); err != nil {
		log.Fatalf("Error en auto-migrate: %v", err)
	}
	log.Println("AutoMigrate completado")

	generador := idgenerator.NewGeneradorUUIDV7()
	publisher := eventpublisher.NewConsolePublisher()

	// Repositorios
	fincaRepo := fincaspostgres.NewFincaRepositorio(db)
	loteRepo := fincaspostgres.NewLoteRepositorio(db)
	muestraRepo := diagnosticopostgres.NewMuestraRepositorio(db)
	diagnosticoRepo := diagnosticopostgres.NewDiagnosticoRepositorio(db)
	candidatoRepo := diagnosticopostgres.NewCandidatoReentrenamientoRepositorio(db)

	// Unit of Work
	fincaUoW := fincaspostgres.NewUnitOfWorkPostgres(db, generador)
	diagnosticoUoW := diagnosticopostgres.NewUnitOfWorkDiagnostico(db, generador)

	// Servicios de dominio
	fincaService := fincasdomain.NewFincaService()

	// Facades
	fincasFacade := facades.NewFincasFacade(
		registrarfinca.NewUseCase(fincaRepo, generador, publisher),
		desactivarfinca.NewUseCase(fincaRepo, fincaService, generador, publisher),
	)
	lotesFacade := facades.NewLotesFacade(
		agregarlote.NewUseCase(fincaRepo, loteRepo, generador, publisher),
		eliminarlote.NewUseCase(loteRepo, generador, publisher),
	)
	muestrasFacade := facades.NewMuestrasFacade(
		tomarmuestra.NewUseCase(loteRepo, muestraRepo, generador, publisher),
		listarmuestrasporlote.NewUseCase(loteRepo, muestraRepo),
	)
	diagnosticosFacade := facades.NewDiagnosticosFacade(
		solicitardiagnosticomanual.NewUseCase(muestraRepo, generador, publisher),
		aceptardiagnostico.NewUseCase(diagnosticoRepo, generador, publisher),
		rechazardiagnostico.NewUseCase(diagnosticoRepo, candidatoRepo, diagnosticoUoW, generador, publisher),
	)
	reportesFacade := facades.NewReportesFacade(
		generarreporteporlote.NewUseCase(loteRepo, muestraRepo, diagnosticoRepo),
	)

	// Token validator
	tokenValidator := jwtvalidator.NewTokenValidator(jwtvalidator.Config{
		Secret: cfg.jwtSecret,
		Issuer: cfg.jwtIssuer,
	})

	// Auth middleware
	authMiddleware := fincasmiddleware.NewAuthMiddleware(tokenValidator)

	// Handlers
	fincaHandler := handler.NewFincaHandler(fincasFacade)
	loteHandler := handler.NewLoteHandler(lotesFacade)
	muestraHandler := handler.NewMuestraHandler(muestrasFacade)
	diagnosticoHandler := handler.NewDiagnosticoHandler(diagnosticosFacade)
	reporteHandler := handler.NewReporteHandler(reportesFacade)

	// Router
	ginEngine := router.New(router.Config{
		AuthMiddleware:      authMiddleware,
		FincaHandler:        fincaHandler,
		LoteHandler:         loteHandler,
		MuestraHandler:      muestraHandler,
		DiagnosticoHandler:  diagnosticoHandler,
		ReporteHandler:      reporteHandler,
	})

	return &Registry{
		db:             db,
		generadorID:    generador,
		publisher:      publisher,
		fincaRepo:      fincaRepo,
		loteRepo:       loteRepo,
		muestraRepo:    muestraRepo,
		diagnosticoRepo: diagnosticoRepo,
		candidatoRepo:  candidatoRepo,
		fincaUoW:       fincaUoW,
		diagnosticoUoW: diagnosticoUoW,
		fincaService:   fincaService,

		RegistrarFinca:             registrarfinca.NewUseCase(fincaRepo, generador, publisher),
		DesactivarFinca:            desactivarfinca.NewUseCase(fincaRepo, fincaService, generador, publisher),
		AgregarLote:                agregarlote.NewUseCase(fincaRepo, loteRepo, generador, publisher),
		EliminarLote:               eliminarlote.NewUseCase(loteRepo, generador, publisher),
		TomarMuestra:               tomarmuestra.NewUseCase(loteRepo, muestraRepo, generador, publisher),
		ListarMuestrasPorLote:      listarmuestrasporlote.NewUseCase(loteRepo, muestraRepo),
		SolicitarDiagnosticoManual: solicitardiagnosticomanual.NewUseCase(muestraRepo, generador, publisher),
		RegistrarInferencia:        registrarinferencia.NewUseCase(muestraRepo, diagnosticoRepo, generador, publisher),
		AceptarDiagnostico:         aceptardiagnostico.NewUseCase(diagnosticoRepo, generador, publisher),
		RechazarDiagnostico:        rechazardiagnostico.NewUseCase(diagnosticoRepo, candidatoRepo, diagnosticoUoW, generador, publisher),
		GenerarReportePorLote:      generarreporteporlote.NewUseCase(loteRepo, muestraRepo, diagnosticoRepo),

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

		router: ginEngine,
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

// Close cierra la conexión a base de datos.
func (r *Registry) Close() {
	sqlDB, err := r.db.DB()
	if err == nil {
		sqlDB.Close()
	}
	log.Println("Conexión a BD cerrada")
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
func (r *Registry) Router() *gin.Engine                             { return r.router }

type config struct {
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	dbSSLMode  string
	jwtSecret  string
	jwtIssuer  string
}

func loadConfig() config {
	getEnv := func(key, fallback string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return fallback
	}
	return config{
		dbHost:     getEnv("DB_HOST", "localhost"),
		dbPort:     getEnv("DB_PORT", "5432"),
		dbUser:     getEnv("DB_USER", "fincas_user"),
		dbPassword: getEnv("DB_PASSWORD", "fincas_pass_dev"),
		dbName:     getEnv("DB_NAME", "bunna_fincas"),
		dbSSLMode:  getEnv("DB_SSLMODE", "disable"),
		jwtSecret:  getEnv("JWT_SECRET", "clave-secreta-dev"),
		jwtIssuer:  getEnv("JWT_ISSUER", ""),
	}
}

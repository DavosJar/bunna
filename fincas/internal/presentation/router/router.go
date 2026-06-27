package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/buffer"
	telemetrymiddleware "github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/middleware"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/handler"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/middleware"
)

// Config contiene la configuración del router.
type Config struct {
	TelemetryEnabled bool
	TelemetryWriter  buffer.BufferWriter
	TelemetryCfg     telemetrymiddleware.Config

	AuthMiddleware      *middleware.AuthMiddleware
	FincaHandler        *handler.FincaHandler
	LoteHandler         *handler.LoteHandler
	MuestraHandler      *handler.MuestraHandler
	DiagnosticoHandler  *handler.DiagnosticoHandler
	ReporteHandler      *handler.ReporteHandler
	NodoHandler         *handler.NodoHandler
}

// New crea un nuevo engine Gin con todas las rutas registradas.
func New(cfg Config) *gin.Engine {
	r := gin.Default()

	// Telemetría RECURSO (HTTP): trace_id/span_id + eventos API
	if cfg.TelemetryEnabled && cfg.TelemetryWriter != nil {
		r.Use(telemetrymiddleware.NewTelemetryMiddleware(cfg.TelemetryWriter, cfg.TelemetryCfg))
	}

	// Health check (público)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Middleware de autenticación
	auth := cfg.AuthMiddleware.RequireAuth()

	// Rutas protegidas — fincas
	fincas := r.Group("/fincas")
	fincas.Use(auth)
	{
		fincas.POST("", cfg.FincaHandler.Registrar)
		fincas.POST("/:id/desactivar", cfg.FincaHandler.Desactivar)
	}

	// Rutas protegidas — lotes
	lotes := r.Group("/lotes")
	lotes.Use(auth)
	{
		lotes.POST("/:id/eliminar", cfg.LoteHandler.Eliminar)
	}

	// Lotes anidados en finca (requieren auth)
	fincas.POST("/:id/lotes", auth, cfg.LoteHandler.Agregar)
	fincas.GET("/:id/lotes/:loteID/muestras", auth, cfg.MuestraHandler.ListarPorLote)
	fincas.POST("/:id/lotes/:loteID/muestras", auth, cfg.MuestraHandler.Tomar)
	fincas.GET("/:id/lotes/:loteID/reporte", auth, cfg.ReporteHandler.GenerarPorLote)
	fincas.GET("/:id/muestras", auth, cfg.MuestraHandler.ListarPorLote)
	fincas.POST("/:id/muestras", auth, cfg.MuestraHandler.Tomar)

	// Rutas protegidas — diagnósticos
	diagnosticos := r.Group("/diagnosticos")
	diagnosticos.Use(auth)
	{
		diagnosticos.POST("/:id/aceptar", cfg.DiagnosticoHandler.Aceptar)
		diagnosticos.POST("/:id/rechazar", cfg.DiagnosticoHandler.Rechazar)
	}

	// Diagnóstico manual anidado en muestra
	r.POST("/muestras/:muestraID/diagnosticos/manual", auth, cfg.DiagnosticoHandler.SolicitarManual)

	// Rutas legacy de lotes (por ID directo)
	lotes.GET("/:id/muestras", auth, cfg.MuestraHandler.ListarPorLote)
	lotes.POST("/:id/muestras", auth, cfg.MuestraHandler.Tomar)
	lotes.GET("/:id/reporte", auth, cfg.ReporteHandler.GenerarPorLote)

	// Rutas internas — nodos (SIN JWT, para YOLO API)
	r.GET("/api/v1/nodos/validar", cfg.NodoHandler.Validar)
	r.POST("/api/v1/diagnosticos/inferencia", cfg.NodoHandler.RegistrarInferencia)

	// Rutas protegidas — nodos
	nodos := r.Group("/nodos")
	nodos.Use(auth)
	{
		nodos.POST("", cfg.NodoHandler.Registrar)
		nodos.GET("", cfg.NodoHandler.Listar)
		nodos.GET("/:id", cfg.NodoHandler.Obtener)
		nodos.PUT("/:id", cfg.NodoHandler.Editar)
		nodos.POST("/:id/desactivar", cfg.NodoHandler.Desactivar)
	}

	return r
}

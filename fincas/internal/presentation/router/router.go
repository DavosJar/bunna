package router

import (
	"net/http"

	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/buffer"
	telemetrymiddleware "github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry/middleware"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/handler"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/middleware"
	"github.com/gin-gonic/gin"
)

// Config contiene la configuración del router.
type Config struct {
	TelemetryEnabled bool
	TelemetryWriter  buffer.BufferWriter
	TelemetryCfg     telemetrymiddleware.Config

	AuthMiddleware     *middleware.AuthMiddleware
	FincaHandler       *handler.FincaHandler
	LoteHandler        *handler.LoteHandler
	MuestraHandler     *handler.MuestraHandler
	DiagnosticoHandler *handler.DiagnosticoHandler
	ReporteHandler     *handler.ReporteHandler
	NodoHandler        *handler.NodoHandler
}

// New crea un nuevo engine Gin con todas las rutas registradas.
func New(cfg Config) *gin.Engine {
	r := gin.Default()

	// Telemetría RECURSO (HTTP): trace_id/span_id + eventos API
	if cfg.TelemetryEnabled && cfg.TelemetryWriter != nil {
		r.Use(telemetrymiddleware.NewTelemetryMiddleware(cfg.TelemetryWriter, cfg.TelemetryCfg))
	}

	// Health check (público)
	healthHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
	r.GET("/health", healthHandler) // Mantenemos este por retrocompatibilidad
	r.GET("/api/v1/fincas/health", healthHandler) // Nueva ruta trazable solicitada

	// Middleware de autenticación
	auth := cfg.AuthMiddleware.RequireAuth()

	api := r.Group("/api/v1/fincas")

	// Rutas protegidas — fincas
	fincas := api.Group("/fincas")
	fincas.Use(auth)
	{
		fincas.POST("", cfg.FincaHandler.Registrar)
		fincas.GET("", cfg.FincaHandler.Listar)
		fincas.POST("/:id/desactivar", cfg.FincaHandler.Desactivar)
	}

	// Rutas protegidas — lotes
	lotes := api.Group("/lotes")
	lotes.Use(auth)
	{
		lotes.POST("/:id/eliminar", cfg.LoteHandler.Eliminar)
	}

	// Lotes anidados en finca (requieren auth)
	fincas.POST("/:id/lotes", auth, cfg.LoteHandler.Agregar)
	fincas.GET("/:id/lotes", auth, cfg.LoteHandler.Listar)
	fincas.GET("/:id/lotes/:loteID/muestras", auth, cfg.MuestraHandler.ListarPorLote)
	fincas.POST("/:id/lotes/:loteID/muestras", auth, cfg.MuestraHandler.Tomar)
	fincas.GET("/:id/lotes/:loteID/reporte", auth, cfg.ReporteHandler.GenerarPorLote)
	fincas.GET("/:id/muestras", auth, cfg.MuestraHandler.ListarPorLote)
	fincas.POST("/:id/muestras", auth, cfg.MuestraHandler.Tomar)

	// Rutas protegidas — diagnósticos
	diagnosticos := api.Group("/diagnosticos")
	diagnosticos.Use(auth)
	{
		diagnosticos.POST("/:id/aceptar", cfg.DiagnosticoHandler.Aceptar)
		diagnosticos.POST("/:id/rechazar", cfg.DiagnosticoHandler.Rechazar)
	}

	// Diagnóstico manual anidado en muestra
	api.POST("/muestras/:muestraID/diagnosticos/manual", auth, cfg.DiagnosticoHandler.SolicitarManual)
	api.POST("/muestras/:muestraID/diagnosticos/manual/resultado", auth, cfg.DiagnosticoHandler.GuardarResultadoManual)

	// Rutas legacy de lotes (por ID directo)
	lotes.GET("/:id/muestras", auth, cfg.MuestraHandler.ListarPorLote)
	lotes.POST("/:id/muestras", auth, cfg.MuestraHandler.Tomar)
	lotes.GET("/:id/reporte", auth, cfg.ReporteHandler.GenerarPorLote)

	// Rutas internas — nodos (SIN JWT, para YOLO API)
	r.GET("/api/v1/nodos/validar", cfg.NodoHandler.Validar)
	r.POST("/api/v1/diagnosticos/inferencia", cfg.NodoHandler.RegistrarInferencia)

	// Rutas protegidas — nodos
	nodos := api.Group("/nodos")
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

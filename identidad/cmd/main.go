package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/router"
	"github.com/davosjar/bunna/services/identidad/internal/rbac/application"
	rbac_postgres "github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	shared_idgenerator "github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/idgenerator"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Ejecutar seed de roles y permisos
	generadorID := shared_idgenerator.NewUUIDv7Generator()
	rolRepo := rbac_postgres.NewRolRepositorio(db)
	permisoRepo := rbac_postgres.NewPermisoRepositorio(db)
	rolPermisoRepo := rbac_postgres.NewRolPermisoRepositorio(db)
	seedSvc := application.NuevoSeedServicio(rolRepo, permisoRepo, rolPermisoRepo, generadorID)
	if err := seedSvc.Ejecutar(context.Background()); err != nil {
		log.Printf("Warning: error en seed de roles: %v", err)
	}

	reg := registry.NewRegistry(db, cfg)
	allFacades := facades.NewAllFacades(reg)
	r := router.New(allFacades, router.Config{
		Version:           "1.0.0",
		CORSOrigins:       []string{cfg.CORSOrigins},
		APIGatewayEnabled: cfg.APIGatewayEnabled,
		TokenSvc:          reg.TokenServicio(),
		RateLimitIPMaxRequests: cfg.RateLimitIPMaxRequests,
		RateLimitIPVentana:     cfg.RateLimitIPVentana,
		TelemetryEnabled:       cfg.TelemetryEnabled,
		TelemetryWriter:        reg.TelemetryWriter,
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	reg.Close()
	log.Println("Server exiting")
}

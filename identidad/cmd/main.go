package main

import (
	"context"
	"log"

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
	})

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

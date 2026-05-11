package main

import (
	"log"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/router"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
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

	reg := registry.NewRegistry(db, cfg)

	// Capa de presentación
	authFacade := facades.NewAuthFacade(reg.GetServicioRegistro(), reg.ServicioLogin)

	r := router.New(authFacade, router.Config{
		Version:     "1.0.0",
		CORSOrigins: []string{cfg.CORSOrigins},
	})

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

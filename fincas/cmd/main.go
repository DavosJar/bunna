package main

import (
	"fmt"
	"log"

	"github.com/davosjar/bunna/services/fincas/internal/config"
	"github.com/davosjar/bunna/services/fincas/internal/presentation"
	"github.com/davosjar/bunna/services/fincas/internal/registry"
	"github.com/davosjar/bunna/services/fincas/internal/shared"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	db, err := shared.NewDB(cfg)
	if err != nil {
		log.Fatalf("Error conectando a base de datos: %v", err)
	}

	if err := shared.AutoMigrate(db); err != nil {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	reg := registry.NewRegistry(db, cfg)

	router := presentation.New(reg.FincaFacade, presentation.RouterConfig{
		Version:     "1.0.0",
		CORSOrigins: []string{"*"},
		JWTSecret:   cfg.JWTSecret,
	})

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Servicio Fincas iniciando en %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

package main

import (
	"log"
	"strings"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/router"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error cargando configuración: %v", err)
	}

	db, err := config.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("error inicializando base de datos: %v", err)
	}

	reg := registry.NewRegistry(db, cfg)

	// Construir facade con interfaces — nunca structs concretos
	authFacade := facades.NewAuthFacade(
		reg.GetServicioRegistro(),
		reg.ServicioLogin,
	)

	// Parsear orígenes CORS
	corsOrigins := []string{}
	if cfg.CORSOrigins != "" && cfg.CORSOrigins != "*" {
		corsOrigins = strings.Split(cfg.CORSOrigins, ",")
	}

	r := router.New(authFacade, router.Config{
		Version:     "1.0.0",
		CORSOrigins: corsOrigins,
	})

	addr := ":" + cfg.Port
	log.Printf("servidor iniciando en %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}

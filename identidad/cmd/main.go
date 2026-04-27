package main

import (
	"log"
	"net/http"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/handler"
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

	reg := registry.NewRegistry(db)
	h := handler.NewHandler(reg)
	h.RegisterRoutes()

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

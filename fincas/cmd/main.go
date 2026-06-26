package main

import (
	"fmt"
	"log"

	"github.com/davosjar/bunna/services/fincas/internal/registry"
)

func main() {
	log.Println("Iniciando aplicación de fincas...")

	r := registry.NewRegistry()
	defer r.Close()

	port := r.ServerPort()
	log.Printf("Servidor HTTP escuchando en :%s", port)
	if err := r.Router().Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}

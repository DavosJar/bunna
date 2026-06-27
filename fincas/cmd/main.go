package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davosjar/bunna/services/fincas/internal/registry"
)

func main() {
	envFlag := flag.String("env", "", "entorno (dev, prod, test)")
	flag.Parse()

	if *envFlag != "" {
		os.Setenv("ENVIRONMENT", *envFlag)
	}

	log.Println("Iniciando aplicación de fincas...")

	r := registry.NewRegistry()
	defer r.Close()

	port := r.ServerPort()
	log.Printf("Servidor HTTP escuchando en :%s", port)
	if err := r.Router().Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}

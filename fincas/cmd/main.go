package main

import (
	"log"

	"github.com/davosjar/bunna/services/fincas/internal/registry"
)

func main() {
	log.Println("Iniciando aplicación de fincas...")

	r := registry.NewRegistry()
	defer r.Close()

	log.Println("Aplicación de fincas iniciada correctamente")

	select {}
}

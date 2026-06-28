package registry

import (
	"context"
	"log"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const (
	maxReintentos        = 10
	esperaInicial        = 2 * time.Second
	multiplicadorBackoff = 2
)

// publicarCatalogoPermisos publica el catálogo de permisos de fincas en Kafka
// durante el startup del servicio. Reintenta con backoff hasta que Kafka responda.
func publicarCatalogoPermisos(ctx context.Context, publisher application.EventPublisher, generador shared.GeneradorID) {
	eventID, err := generador.NextID(ctx)
	if err != nil {
		log.Printf("[WARN] No se pudo generar ID para evento de catálogo de permisos: %v", err)
		return
	}

	evento := application.CatalogoPermisosPublicado{
		EventID:   eventID,
		Tipo:      "permisos.catalogo",
		Origen:    "fincas",
		Modulo:    application.ModuloFincas,
		Permisos:  application.CatalogoFincas,
		Version:   "1.0",
		OcurredAt: time.Now(),
	}

	espera := esperaInicial
	for intento := 1; intento <= maxReintentos; intento++ {
		if err := publisher.Publish(ctx, "dev.permisos", evento); err == nil {
			log.Printf("[INFO] Catálogo de permisos publicado (%d permisos)", len(application.CatalogoFincas))
			return
		} else {
			log.Printf("[WARN] Intento %d/%d — Error publicando catálogo de permisos: %v", intento, maxReintentos, err)
		}

		if intento < maxReintentos {
			log.Printf("[INFO] Reintentando en %v...", espera)
			time.Sleep(espera)
			espera *= multiplicadorBackoff
		}
	}

	log.Printf("[ERROR] No se pudo publicar el catálogo de permisos tras %d intentos", maxReintentos)
}

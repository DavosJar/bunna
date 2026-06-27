package consumers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/segmentio/kafka-go"
)

// CatalogoPermisosPublicado coincide con la estructura enviada por otros módulos
type CatalogoPermisosPublicado struct {
	EventID   string           `json:"event_id"`
	Tipo      string           `json:"tipo"`
	Origen    string           `json:"origen"`
	Modulo    string           `json:"modulo"`
	Permisos  []PermisoEvento `json:"permisos"`
	Version   string           `json:"version"`
	OcurredAt time.Time        `json:"ocurred_at"`
}

type PermisoEvento struct {
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Modulo      string `json:"modulo"`
}

// PermisosConsumer escucha eventos de catálogo de permisos y los sincroniza en la DB
type PermisosConsumer struct {
	reader         *kafka.Reader
	permisoRepo    rbac.PermisoRepositorio
	generadorID    interface{ NextID(ctx context.Context) (string, error) }
}

func NewPermisosConsumer(brokers []string, topic string, permisoRepo rbac.PermisoRepositorio, generadorID interface{ NextID(ctx context.Context) (string, error) }) *PermisosConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "identidad-permisos-sync",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &PermisosConsumer{
		reader:      r,
		permisoRepo: permisoRepo,
		generadorID: generadorID,
	}
}

func (c *PermisosConsumer) Start(ctx context.Context) {
	log.Printf("[PermisosConsumer] Escuchando en tópico: %s", c.reader.Config().Topic)
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("[PermisosConsumer] Error leyendo mensaje: %v", err)
			continue
		}

		var event CatalogoPermisosPublicado
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("[PermisosConsumer] Error unmarshaling: %v", err)
			continue
		}

		if event.Tipo == "permisos.catalogo" {
			c.sincronizarPermisos(ctx, event)
		}
	}
}

func (c *PermisosConsumer) sincronizarPermisos(ctx context.Context, event CatalogoPermisosPublicado) {
	log.Printf("[PermisosConsumer] Recibido catálogo de '%s' con %d permisos", event.Origen, len(event.Permisos))

	for _, p := range event.Permisos {
		existente, err := c.permisoRepo.ObtenerPorCodigo(ctx, p.Codigo)
		if err != nil {
			log.Printf("[PermisosConsumer] Error consultando permiso %s: %v", p.Codigo, err)
			continue
		}

		if existente == nil {
			nuevoID, _ := c.generadorID.NextID(ctx)
			nuevo := &rbac.PermisoDB{
				ID:          nuevoID,
				Codigo:      p.Codigo,
				Nombre:      p.Nombre,
				Descripcion: p.Descripcion,
				Modulo:      p.Modulo,
			}
			if err := c.permisoRepo.Crear(ctx, nuevo); err != nil {
				log.Printf("[PermisosConsumer] Error creando permiso %s: %v", p.Codigo, err)
			} else {
				log.Printf("[PermisosConsumer] Permiso creado: %s", p.Codigo)
			}
		} else {
			// Siempre actualizar para que fincas mande la verdad
			if err := c.permisoRepo.ActualizarNombreDescripcion(ctx, existente.ID, p.Nombre, p.Descripcion); err != nil {
				log.Printf("[PermisosConsumer] Error actualizando permiso %s: %v", p.Codigo, err)
			}
		}
	}
}

func (c *PermisosConsumer) Close() error {
	return c.reader.Close()
}

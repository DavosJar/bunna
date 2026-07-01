package consumers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	tenant_domain "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
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
	rolRepo        rbac.RolRepositorio
	rolPermisoRepo rbac.RolPermisoRepositorio
	publisher      rbac.RolPublisher
	tenantRepo     tenant_domain.TenantRepositorio
	generadorID    interface{ NextID(ctx context.Context) (string, error) }
}

func NewPermisosConsumer(
	brokers []string,
	topic string,
	permisoRepo rbac.PermisoRepositorio,
	rolRepo rbac.RolRepositorio,
	rolPermisoRepo rbac.RolPermisoRepositorio,
	publisher rbac.RolPublisher,
	tenantRepo tenant_domain.TenantRepositorio,
	generadorID interface{ NextID(ctx context.Context) (string, error) },
) *PermisosConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "identidad-permisos-sync",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &PermisosConsumer{
		reader:         r,
		permisoRepo:    permisoRepo,
		rolRepo:        rolRepo,
		rolPermisoRepo: rolPermisoRepo,
		publisher:      publisher,
		tenantRepo:     tenantRepo,
		generadorID:    generadorID,
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

	// Obtener rol administrador (el único rol que se propaga a los tenants)
	admin, _ := c.rolRepo.ObtenerPorNombre(ctx, rbac.RolAdministrador)

	adminModificado := false

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
				// Auto-asignar al rol administrador en el tenant de sistema
				if admin != nil {
					_ = c.rolPermisoRepo.AsignarPermiso(ctx, admin.ID, nuevoID, rbac.TenantIDSistema, "")
					adminModificado = true
				}
			}
		} else {
			// Siempre actualizar para que fincas mande la verdad
			if err := c.permisoRepo.ActualizarNombreDescripcion(ctx, existente.ID, p.Nombre, p.Descripcion); err != nil {
				log.Printf("[PermisosConsumer] Error actualizando permiso %s: %v", p.Codigo, err)
			}
		}
	}

	// Registrar TODOS los permisos del catálogo en el mapa en memoria del admin.
	// Esto hace que TienePermisoEnRol los encuentre inmediatamente sin consultar la BD.
	var todosLosCodigos []string
	for _, p := range event.Permisos {
		todosLosCodigos = append(todosLosCodigos, p.Codigo)
	}
	rbac.RegistrarPermisosDeModulo(rbac.RolAdministrador, todosLosCodigos)

	// Si el rol administrador fue modificado, publicar la lista actualizada
	// de permisos para CADA tenant existente en el sistema
	if adminModificado && admin != nil {
		c.publicarPermisosPorTenant(ctx, admin.ID)
	}
}

// publicarPermisosPorTenant obtiene todos los permisos del rol administrador
// y publica un evento de Kafka por cada tenant existente en el sistema.
// Esto permite que Fincas valide usando el tenant_id exacto del JWT.
func (c *PermisosConsumer) publicarPermisosPorTenant(ctx context.Context, adminRolID string) {
	// Obtener todos los permisos del administrador (asignados en el tenant de sistema)
	permisosDB, err := c.rolPermisoRepo.ListarPorRolYTenant(ctx, adminRolID, rbac.TenantIDSistema)
	if err != nil {
		log.Printf("[PermisosConsumer] Error obteniendo permisos del administrador: %v", err)
		return
	}

	var codigos []string
	for _, p := range permisosDB {
		codigos = append(codigos, p.Codigo)
	}

	// Obtener todos los tenants del sistema
	tenants, err := c.tenantRepo.Listar(ctx)
	if err != nil {
		log.Printf("[PermisosConsumer] Error listando tenants: %v", err)
		return
	}

	// Publicar un evento por cada tenant
	for _, t := range tenants {
		if err := c.publisher.PublicarRolActualizado(ctx, rbac.RolAdministrador, t.ID(), codigos); err != nil {
			log.Printf("[PermisosConsumer] Error publicando permisos para tenant %s: %v", t.ID(), err)
		} else {
			log.Printf("[PermisosConsumer] Permisos de administrador publicados para tenant: %s (%s) con %d permisos", t.Nombre(), t.ID(), len(codigos))
		}
	}

	// También publicar para el tenant de sistema (por si acaso)
	_ = c.publisher.PublicarRolActualizado(ctx, rbac.RolAdministrador, rbac.TenantIDSistema, codigos)
	log.Printf("[PermisosConsumer] Permisos de administrador publicados para tenant de sistema con %d permisos", len(codigos))
}

func (c *PermisosConsumer) Close() error {
	return c.reader.Close()
}

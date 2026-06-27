package agregarlote

import (
	"context"
	"errors"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = application.PermisoCrearLote
const routingKey = "fincas.v1.lote.creado"

// LoteCreado es el evento publicado tras agregar un lote exitosamente.
type LoteCreado struct {
	EventID   string    `json:"event_id"`
	LoteID    string    `json:"lote_id"`
	FincaID   string    `json:"finca_id"`
	Nombre    string    `json:"nombre"`
	Area      float64   `json:"area"`
	TenantID  string    `json:"tenant_id,omitempty"`
	OcurredAt time.Time `json:"ocurred_at"`
}

// UseCase orquesta el registro de un nuevo lote en una finca.
type UseCase struct {
	fincaRepo domain.FincaRepositorio
	loteRepo  domain.LoteRepositorio
	generador shared.GeneradorID
	publisher application.EventPublisher
}

func NewUseCase(
	fincaRepo domain.FincaRepositorio,
	loteRepo domain.LoteRepositorio,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		fincaRepo: fincaRepo,
		loteRepo:  loteRepo,
		generador: generador,
		publisher: publisher,
	}
}

func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	// 1. Validar permisos
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	// 2. Validar campos del comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 3. Cargar finca
	finca, err := uc.fincaRepo.ObtenerPorID(ctx, cmd.FincaID)
	if err != nil {
		if errors.Is(err, domain.ErrFincaNoEncontrada) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	// 4. Validar tenencia: si el TenantID no coincide → ErrNotFound (sec 3.2)
	if auth.TenantID != "" && finca.TenantID() != nil && *finca.TenantID() != auth.TenantID {
		return nil, application.ErrNotFound
	}

	// 5. Verificar que la finca esté ACTIVA
	if finca.Estado() != domain.FincaActiva {
		return nil, application.ErrConflictoEstado("No se pueden agregar lotes a una finca en estado PENDIENTE_ELIMINACION")
	}

	// 6. Generar ID
	id, err := uc.generador.NextID(ctx)
	if err != nil {
		return nil, err
	}

	// 7. Asignar tenantID: heredar de la finca
	var tenantID string
	if finca.TenantID() != nil {
		tenantID = *finca.TenantID()
	}

	// 8. Construir entidad Lote y asignar ID vía reconstrucción desde persistencia
	now := time.Now()
	lote := domain.NewLoteFromPersistence(
		id, cmd.FincaID, tenantID, cmd.Nombre, cmd.Area, cmd.Descripcion,
		domain.LoteActivo, now, now,
	)

	// 9. Persistir
	if err := uc.loteRepo.Crear(ctx, lote); err != nil {
		return nil, err
	}

	// 10. Publicar evento
	eventID, _ := uc.generador.NextID(ctx)
	evento := LoteCreado{
		EventID:   eventID,
		LoteID:    id,
		FincaID:   cmd.FincaID,
		Nombre:    cmd.Nombre,
		Area:      cmd.Area,
		TenantID:  tenantID,
		OcurredAt: now,
	}
	_ = uc.publisher.Publish(ctx, routingKey, evento)

	return &Salida{
		ID:          id,
		FincaID:     cmd.FincaID,
		Nombre:      cmd.Nombre,
		Area:        cmd.Area,
		Descripcion: cmd.Descripcion,
		Estado:      string(domain.LoteActivo),
		CreatedAt:   lote.CreatedAt(),
	}, nil
}

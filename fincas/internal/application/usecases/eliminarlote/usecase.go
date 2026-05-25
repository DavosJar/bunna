package eliminarlote

import (
	"context"
	"errors"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = "ELIMINAR_LOTE"

const routingKey = "fincas.v1.lote.eliminado"

// LoteEliminado es el evento publicado tras eliminar un lote.
type LoteEliminado struct {
	EventID        string    `json:"event_id"`
	LoteID         string    `json:"lote_id"`
	FincaID        string    `json:"finca_id"`
	EstadoAnterior string    `json:"estado_anterior"`
	EstadoNuevo    string    `json:"estado_nuevo"`
	TenantID       string    `json:"tenant_id,omitempty"`
	OcurredAt      time.Time `json:"ocurred_at"`
}

// UseCase orquesta la eliminación de un lote (ACTIVO → ELIMINADO).
type UseCase struct {
	loteRepo  domain.LoteRepositorio
	generador shared.GeneradorID
	publisher application.EventPublisher
}

func NewUseCase(
	loteRepo domain.LoteRepositorio,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
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

	// 3. Cargar lote
	lote, err := uc.loteRepo.ObtenerPorID(ctx, cmd.LoteID)
	if err != nil {
		if errors.Is(err, domain.ErrLoteNoEncontrado) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	// 4. Validar tenencia: si el TenantID no coincide → ErrNotFound (sec 4.4)
	if auth.TenantID != "" && lote.TenantID() != auth.TenantID {
		return nil, application.ErrNotFound
	}

	// 5. Verificar estado actual — si ya está eliminado, no se puede volver a eliminar
	if lote.Estado() == domain.LoteEliminado {
		return nil, application.ErrConflictoEstado("el lote ya está eliminado")
	}

	// 6. Guardar estado anterior antes del cambio
	estadoAnterior := string(lote.Estado())

	// 7. Ejecutar cambio de estado (ACTIVO → ELIMINADO)
	if err := lote.CambiarEstado(domain.LoteEliminado); err != nil {
		if errors.Is(err, domain.ErrTransicionEstadoNoPermitida) {
			return nil, application.ErrConflictoEstado("el lote no está en estado ACTIVO")
		}
		return nil, err
	}

	// 8. Persistir cambio de estado
	if err := uc.loteRepo.Actualizar(ctx, lote); err != nil {
		return nil, err
	}

	// 9. Publicar evento
	eventID, _ := uc.generador.NextID(ctx)
	evento := LoteEliminado{
		EventID:        eventID,
		LoteID:         cmd.LoteID,
		FincaID:        lote.FincaID(),
		EstadoAnterior: estadoAnterior,
		EstadoNuevo:    string(domain.LoteEliminado),
		TenantID:       auth.TenantID,
		OcurredAt:      time.Now(),
	}
	_ = uc.publisher.Publish(ctx, routingKey, evento)

	return &Salida{
		ID:        lote.ID(),
		Estado:    string(lote.Estado()),
		UpdatedAt: lote.UpdatedAt(),
	}, nil
}

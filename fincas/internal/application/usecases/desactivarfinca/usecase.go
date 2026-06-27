package desactivarfinca

import (
	"context"
	"errors"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = application.PermisoDesactivarFinca

// FincaDesactivada es el evento publicado tras desactivar una finca.
type FincaDesactivada struct {
	EventID        string    `json:"event_id"`
	FincaID        string    `json:"finca_id"`
	EstadoAnterior string    `json:"estado_anterior"`
	EstadoNuevo    string    `json:"estado_nuevo"`
	UsuarioID      string    `json:"usuario_id"`
	TenantID       string    `json:"tenant_id,omitempty"`
	OcurredAt      time.Time `json:"ocurred_at"`
}

// UseCase orquesta la desactivación de una finca (ACTIVA → PENDIENTE_ELIMINACION).
type UseCase struct {
	fincaRepo    domain.FincaRepositorio
	fincaService *domain.FincaService
	generador    shared.GeneradorID
	publisher    application.EventPublisher
}

func NewUseCase(
	fincaRepo domain.FincaRepositorio,
	fincaService *domain.FincaService,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		fincaRepo:    fincaRepo,
		fincaService: fincaService,
		generador:    generador,
		publisher:    publisher,
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

	// 5. Contar lotes activos
	cantidadLotes, err := uc.fincaRepo.ContarLotes(ctx, cmd.FincaID)
	if err != nil {
		return nil, err
	}

	// 6. Invocar servicio de dominio
	estadoAnterior := string(finca.Estado())
	if err := uc.fincaService.EliminarFincaConLotes(finca, cantidadLotes, cmd.Confirmar); err != nil {
		if errors.Is(err, domain.ErrTransicionEstadoNoPermitida) {
			return nil, application.ErrConflictoEstado("la finca no está en estado ACTIVA")
		}
		// ErrFincaConLotes es un error con formato dinámico, mapeamos a ErrValidacion
		return nil, application.ErrValidacion(err.Error())
	}

	// 7. Persistir cambio de estado
	if err := uc.fincaRepo.Actualizar(ctx, finca); err != nil {
		return nil, err
	}

	// 8. Publicar evento
	eventID, _ := uc.generador.NextID(ctx)
	evento := FincaDesactivada{
		EventID:        eventID,
		FincaID:        cmd.FincaID,
		EstadoAnterior: estadoAnterior,
		EstadoNuevo:    string(domain.FincaPendienteEliminar),
		UsuarioID:      auth.UsuarioID,
		TenantID:       auth.TenantID,
		OcurredAt:      time.Now(),
	}
	_ = uc.publisher.Publish(ctx, "fincas.v1.finca.desactivada", evento)

	return &Salida{
		ID:        finca.ID(),
		Estado:    string(finca.Estado()),
		UpdatedAt: finca.UpdatedAt(),
	}, nil
}

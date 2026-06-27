package tomarmuestra

import (
	"context"
	"errors"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	fincasdomain "github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = application.PermisoCrearMuestra
const routingKey = "diagnosticos.v1.muestra.tomada"

// MuestraTomada es el evento publicado tras tomar una muestra exitosamente.
type MuestraTomada struct {
	EventID   string    `json:"event_id"`
	MuestraID string    `json:"muestra_id"`
	LoteID    string    `json:"lote_id"`
	Latitud   float64   `json:"latitud"`
	Longitud  float64   `json:"longitud"`
	TenantID  string    `json:"tenant_id,omitempty"`
	OcurredAt time.Time `json:"ocurred_at"`
}

// UseCase orquesta la toma de una muestra en un lote.
type UseCase struct {
	loteRepo   fincasdomain.LoteRepositorio
	muestraRepo diagnosticodomain.MuestraRepositorio
	generador  shared.GeneradorID
	publisher  application.EventPublisher
}

// NewUseCase crea una nueva instancia del caso de uso TomarMuestra.
func NewUseCase(
	loteRepo fincasdomain.LoteRepositorio,
	muestraRepo diagnosticodomain.MuestraRepositorio,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		loteRepo:    loteRepo,
		muestraRepo: muestraRepo,
		generador:   generador,
		publisher:   publisher,
	}
}

// Ejecutar valida permisos, carga el lote, construye y persiste la muestra,
// y publica el evento MuestraTomada.
func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	// 1. Validar permisos
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	// 2. Validar campos del comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 3. Cargar lote (opcional)
	var tenantID = auth.TenantID
	if cmd.LoteID != "" {
		lote, err := uc.loteRepo.ObtenerPorID(ctx, cmd.LoteID)
		if err != nil {
			if errors.Is(err, fincasdomain.ErrLoteNoEncontrado) {
				return nil, application.ErrNotFound
			}
			return nil, err
		}
		
		// 4. Validar tenencia: si el TenantID no coincide → ErrNotFound (sec 3.2)
		if tenantID != "" && lote.TenantID() != tenantID {
			return nil, application.ErrNotFound
		}

		// 5. Verificar que el lote esté ACTIVO
		if lote.Estado() == fincasdomain.LoteEliminado {
			return nil, application.ErrConflictoEstado("No se pueden tomar muestras en un lote eliminado")
		}

		if tenantID == "" {
			tenantID = lote.TenantID()
		}
	}

	// 6. Construir Ubicacion VO
	ubicacion, err := diagnosticodomain.NewUbicacion(cmd.Latitud, cmd.Longitud)
	if err != nil {
		return nil, application.ErrValidacion(err.Error())
	}

	// 7. Generar ID para la muestra
	id, err := uc.generador.NextID(ctx)
	if err != nil {
		return nil, err
	}

	// 8. Determinar tenantID si no se asignó aún
	if tenantID == "" {
		tenantID = auth.TenantID
	}

	// 9. Construir entidad Muestra
	muestra, err := diagnosticodomain.NewMuestra(id, cmd.FincaID, cmd.LoteID, *ubicacion, tenantID)
	if err != nil {
		return nil, application.ErrValidacion(err.Error())
	}

	// 10. Persistir
	if err := uc.muestraRepo.Crear(ctx, muestra); err != nil {
		return nil, err
	}

	// 11. Publicar evento
	eventID, _ := uc.generador.NextID(ctx)
	evento := MuestraTomada{
		EventID:   eventID,
		MuestraID: id,
		LoteID:    cmd.LoteID,
		Latitud:   cmd.Latitud,
		Longitud:  cmd.Longitud,
		TenantID:  tenantID,
		OcurredAt: time.Now(),
	}
	_ = uc.publisher.Publish(ctx, routingKey, evento)

	return &Salida{
		ID:        id,
		FincaID:   cmd.FincaID,
		LoteID:    cmd.LoteID,
		Latitud:   cmd.Latitud,
		Longitud:  cmd.Longitud,
		CreatedAt: time.Now(),
	}, nil
}

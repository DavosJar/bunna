package solicitardiagnosticomanual

import (
	"context"
	"errors"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = "SOLICITAR_DIAGNOSTICO"
const routingKey = "diagnosticos.v1.solicitud.diagnostico.manual"

// SolicitudDiagnosticoManual es el evento publicado tras solicitar un diagnóstico manual.
type SolicitudDiagnosticoManual struct {
	SolicitudID string    `json:"solicitud_id"`
	MuestraID   string    `json:"muestra_id"`
	ImageURL    string    `json:"image_url"`
	TenantID    string    `json:"tenant_id,omitempty"`
	UsuarioID   string    `json:"usuario_id,omitempty"`
	OcurredAt   time.Time `json:"ocurred_at"`
}

// UseCase orquesta la solicitud de un diagnóstico manual para una muestra.
type UseCase struct {
	muestraRepo diagnosticodomain.MuestraRepositorio
	generador   shared.GeneradorID
	publisher   application.EventPublisher
}

// NewUseCase crea una nueva instancia del caso de uso SolicitarDiagnosticoManual.
func NewUseCase(
	muestraRepo diagnosticodomain.MuestraRepositorio,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		muestraRepo: muestraRepo,
		generador:   generador,
		publisher:   publisher,
	}
}

// Ejecutar valida permisos, carga la muestra, verifica tenencia, genera un ID de
// correlación y publica el evento SolicitudDiagnosticoManual. No persiste nada en BD.
func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	// 1. Validar permisos
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	// 2. Validar campos del comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 3. Cargar muestra
	muestra, err := uc.muestraRepo.ObtenerPorID(ctx, cmd.MuestraID)
	if err != nil {
		if errors.Is(err, diagnosticodomain.ErrMuestraNoEncontrada) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	// 4. Validar tenencia: si el TenantID no coincide → ErrNotFound (sec 4.7)
	if auth.TenantID != "" && muestra.TenantID() != auth.TenantID {
		return nil, application.ErrNotFound
	}

	// 5. Generar SolicitudID para correlación
	solicitudID, err := uc.generador.NextID(ctx)
	if err != nil {
		return nil, err
	}

	// 6. Publicar evento — NO persiste nada en BD
	evento := SolicitudDiagnosticoManual{
		SolicitudID: solicitudID,
		MuestraID:   muestra.ID(),
		ImageURL:    cmd.ImageURL,
		TenantID:    muestra.TenantID(),
		UsuarioID:   auth.UsuarioID,
		OcurredAt:   time.Now(),
	}

	if err := uc.publisher.Publish(ctx, routingKey, evento); err != nil {
		return nil, err
	}

	// 7. Retornar salida
	return &Salida{
		SolicitudID:  solicitudID,
		MuestraID:    muestra.ID(),
		SolicitadoEn: time.Now(),
	}, nil
}

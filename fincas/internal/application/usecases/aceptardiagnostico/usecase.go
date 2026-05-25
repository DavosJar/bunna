package aceptardiagnostico

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = "ACEPTAR_DIAGNOSTICO"
const routingKey = "diagnosticos.v1.diagnostico.aceptado"

// DiagnosticoAceptado es el evento publicado tras aceptar un diagnóstico exitosamente.
type DiagnosticoAceptado struct {
	EventID        string    `json:"event_id"`
	DiagnosticoID  string    `json:"diagnostico_id"`
	MuestraID      string    `json:"muestra_id"`
	EstadoAnterior string    `json:"estado_anterior"`
	EstadoNuevo    string    `json:"estado_nuevo"`
	UsuarioID      string    `json:"usuario_id"`
	TenantID       string    `json:"tenant_id,omitempty"`
	OcurredAt      time.Time `json:"ocurred_at"`
}

// UseCase orquesta la aceptación de un diagnóstico (PENDIENTE → ACEPTADO).
type UseCase struct {
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio
	generador       shared.GeneradorID
	publisher       application.EventPublisher
}

// NewUseCase crea una nueva instancia del caso de uso AceptarDiagnostico.
func NewUseCase(
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		diagnosticoRepo: diagnosticoRepo,
		generador:       generador,
		publisher:       publisher,
	}
}

// Ejecutar valida permisos, carga el diagnóstico, verifica tenencia y estado,
// transiciona a ACEPTADO, persiste y publica el evento.
func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	// 1. Validar permisos
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	// 2. Validar campos del comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 3. Cargar diagnóstico
	diagnostico, err := uc.diagnosticoRepo.ObtenerPorID(ctx, cmd.DiagnosticoID)
	if err != nil {
		if errors.Is(err, diagnosticodomain.ErrDiagnosticoNoEncontrado) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	// 4. Validar tenencia: si el TenantID no coincide → ErrNotFound (sec 3.2)
	if auth.TenantID != "" && diagnostico.TenantID() != auth.TenantID {
		return nil, application.ErrNotFound
	}

	// 5. Validar estado actual: debe ser PENDIENTE
	if diagnostico.Estado() != diagnosticodomain.EstadoDiagnosticoPendiente {
		return nil, application.ErrConflictoEstado(
			fmt.Sprintf("El diagnóstico no está pendiente. Estado actual: %s", diagnostico.Estado()),
		)
	}

	// 6. Guardar estado anterior antes de mutar
	estadoAnterior := string(diagnostico.Estado())

	// 7. Transicionar a ACEPTADO
	if err := diagnostico.MarcarComoAceptado(); err != nil {
		return nil, application.ErrConflictoEstado(
			fmt.Sprintf("El diagnóstico no puede ser aceptado. Estado actual: %s", estadoAnterior),
		)
	}

	// 8. Persistir cambio de estado
	if err := uc.diagnosticoRepo.Actualizar(ctx, diagnostico); err != nil {
		return nil, err
	}

	// 9. Publicar evento (ignorar error de publish)
	eventID, _ := uc.generador.NextID(ctx)
	evento := DiagnosticoAceptado{
		EventID:        eventID,
		DiagnosticoID:  cmd.DiagnosticoID,
		MuestraID:      diagnostico.MuestrasId(),
		EstadoAnterior: estadoAnterior,
		EstadoNuevo:    string(diagnosticodomain.EstadoDiagnosticoAceptado),
		UsuarioID:      auth.UsuarioID,
		TenantID:       auth.TenantID,
		OcurredAt:      time.Now(),
	}
	_ = uc.publisher.Publish(ctx, routingKey, evento)

	return &Salida{
		ID:        diagnostico.ID(),
		Estado:    string(diagnostico.Estado()),
		UpdatedAt: diagnostico.UpdatedAt(),
	}, nil
}

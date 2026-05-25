package rechazardiagnostico

import (
	"context"
	"errors"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = "RECHAZAR_DIAGNOSTICO"
const routingKey = "diagnosticos.v1.diagnostico.rechazado"

// DiagnosticoRechazado es el evento publicado tras rechazar un diagnóstico exitosamente.
type DiagnosticoRechazado struct {
	EventID            string    `json:"event_id"`
	DiagnosticoID      string    `json:"diagnostico_id"`
	MuestraID          string    `json:"muestra_id"`
	EstadoAnterior     string    `json:"estado_anterior"`
	EstadoNuevo        string    `json:"estado_nuevo"`
	Motivo             string    `json:"motivo,omitempty"`
	EsCandidatoRetrain bool      `json:"es_candidato_retrain"`
	UsuarioID          string    `json:"usuario_id"`
	TenantID           string    `json:"tenant_id"`
	OcurredAt          time.Time `json:"ocurred_at"`
}

// UseCase orquesta el rechazo de un diagnóstico con creación atómica
// del candidato a reentrenamiento.
type UseCase struct {
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio
	candidatoRepo   diagnosticodomain.CandidatoReentrenamientoRepositorio
	uow             application.UnitOfWorkDiagnostico
	generador       shared.GeneradorID
	publisher       application.EventPublisher
}

// NewUseCase crea una nueva instancia del caso de uso RechazarDiagnostico.
func NewUseCase(
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio,
	candidatoRepo diagnosticodomain.CandidatoReentrenamientoRepositorio,
	uow application.UnitOfWorkDiagnostico,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		diagnosticoRepo: diagnosticoRepo,
		candidatoRepo:   candidatoRepo,
		uow:             uow,
		generador:       generador,
		publisher:       publisher,
	}
}

// Ejecutar valida el permiso, el comando, carga el diagnóstico, cambia el estado,
// persiste atómicamente (diagnóstico + candidato) y publica el evento.
func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	// 1. Validar permiso
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

	// 4. Validar tenencia
	if diagnostico.TenantID() != auth.TenantID {
		return nil, application.ErrNotFound
	}

	// 5. Validar estado: solo se puede rechazar un diagnóstico PENDIENTE
	if diagnostico.Estado() != diagnosticodomain.EstadoDiagnosticoPendiente {
		return nil, application.ErrConflictoEstado(
			"El diagnóstico no está pendiente. Estado actual: " + string(diagnostico.Estado()),
		)
	}

	// 6. Guardar estado anterior para el evento
	estadoAnterior := string(diagnostico.Estado())

	// 7. Ejecutar cambio de estado en la entidad
	if err := diagnostico.MarcarComoRechazado(); err != nil {
		return nil, application.ErrConflictoEstado("No se puede rechazar el diagnóstico")
	}

	// 8. Extraer datos de la inferencia para el candidato
	resultado := diagnostico.ResultadoInferencia()

	// 9. Persistir atómicamente dentro de la Unit of Work
	err = uc.uow.Transaccional(ctx, func(tx application.UnitOfWorkDiagnostico) error {
		// a. Actualizar estado del diagnóstico
		if err := tx.DiagnosticoRepo().Actualizar(ctx, diagnostico); err != nil {
			return err
		}

		// b. Generar ID para el candidato a reentrenamiento
		candidatoID, err := uc.generador.NextID(ctx)
		if err != nil {
			return err
		}

		// c. Construir candidato con los datos de la inferencia
		candidato, err := diagnosticodomain.NewCandidatoReentrenamiento(
			candidatoID,
			diagnostico.ID(),
			resultado.ImageUrl(),
			resultado.TieneClorosis(),
			resultado.Confianza(),
			cmd.Motivo,
			auth.UsuarioID,
		)
		if err != nil {
			return err
		}

		// d. Persistir candidato
		return tx.CandidatoRepo().Crear(ctx, candidato)
	})
	if err != nil {
		return nil, err
	}

	// 10. Publicar evento (ignorar error)
	eventID, _ := uc.generador.NextID(ctx)
	evento := DiagnosticoRechazado{
		EventID:            eventID,
		DiagnosticoID:      diagnostico.ID(),
		MuestraID:          diagnostico.MuestrasId(),
		EstadoAnterior:     estadoAnterior,
		EstadoNuevo:        string(diagnostico.Estado()),
		Motivo:             cmd.Motivo,
		EsCandidatoRetrain: true,
		UsuarioID:          auth.UsuarioID,
		TenantID:           auth.TenantID,
		OcurredAt:          time.Now(),
	}
	_ = uc.publisher.Publish(ctx, routingKey, evento)

	return &Salida{
		ID:        diagnostico.ID(),
		Estado:    string(diagnostico.Estado()),
		Motivo:    cmd.Motivo,
		UpdatedAt: diagnostico.UpdatedAt(),
	}, nil
}

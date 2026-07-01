package registrarinferencia

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const routingKey = "diagnosticos.v1.diagnostico.creado"

// DiagnosticoCreado es el evento publicado tras registrar un diagnóstico exitosamente.
type DiagnosticoCreado struct {
	EventID        string    `json:"event_id"`
	DiagnosticoID  string    `json:"diagnostico_id"`
	MuestraID      string    `json:"muestra_id"`
	Estado         string    `json:"estado"`
	TieneClorosis  bool      `json:"tiene_clorosis"`
	Confianza      float64   `json:"confianza"`
	TenantID       string    `json:"tenant_id,omitempty"`
	OcurredAt      time.Time `json:"ocurred_at"`
}

// UseCase orquesta el registro de una inferencia (caso de uso interno).
// No verifica permisos de usuario porque es disparado por el consumer de RabbitMQ.
type UseCase struct {
	muestraRepo    diagnosticodomain.MuestraRepositorio
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio
	generador      shared.GeneradorID
	publisher      application.EventPublisher
}

// NewUseCase crea una nueva instancia del caso de uso RegistrarInferencia.
func NewUseCase(
	muestraRepo diagnosticodomain.MuestraRepositorio,
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		muestraRepo:    muestraRepo,
		diagnosticoRepo: diagnosticoRepo,
		generador:      generador,
		publisher:      publisher,
	}
}

// Ejecutar valida el comando, carga la muestra, construye el diagnóstico,
// lo persiste y publica el evento. No recibe AuthContext por ser operación interna.
func (uc *UseCase) Ejecutar(ctx context.Context, cmd Command) (*Salida, error) {
	// 1. Validar campos del comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 2. Cargar muestra
	muestra, err := uc.muestraRepo.ObtenerPorID(ctx, cmd.MuestraID)
	if err != nil {
		if errors.Is(err, diagnosticodomain.ErrMuestraNoEncontrada) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	// 3. Construir ResultadoInferencia VO
	resultadoInferencia, err := diagnosticodomain.NewResultadoInferencia(
		cmd.ImageURL, "", cmd.TieneClorosis, cmd.Confianza, cmd.ProcesadoAt,
	)
	if err != nil {
		return nil, application.ErrValidacion(err.Error())
	}

	// 4. Generar ID para el diagnóstico
	id, err := uc.generador.NextID(ctx)
	if err != nil {
		return nil, err
	}

	// 5. Generar nombre autogenerado
	nombre := fmt.Sprintf("INF-%s-%04d", time.Now().Format("20060102"), rand.Intn(10000))

	// 6. Extraer tenantID de la muestra cargada
	tenantID := muestra.TenantID()

	// 7. Construir entidad Diagnostico
	diagnostico, err := diagnosticodomain.NewDiagnostico(id, nombre, cmd.MuestraID, tenantID, resultadoInferencia)
	if err != nil {
		return nil, application.ErrValidacion(err.Error())
	}

	// 8. Persistir diagnóstico
	if err := uc.diagnosticoRepo.Crear(ctx, diagnostico); err != nil {
		return nil, err
	}

	// 9. Publicar evento
	eventID, _ := uc.generador.NextID(ctx)
	evento := DiagnosticoCreado{
		EventID:       eventID,
		DiagnosticoID: id,
		MuestraID:     cmd.MuestraID,
		Estado:        string(diagnostico.Estado()),
		TieneClorosis: cmd.TieneClorosis,
		Confianza:     cmd.Confianza,
		TenantID:      tenantID,
		OcurredAt:     time.Now(),
	}
	_ = uc.publisher.Publish(ctx, routingKey, evento)

	return &Salida{
		ID:            id,
		MuestraID:     cmd.MuestraID,
		Nombre:        nombre,
		Estado:        string(diagnostico.Estado()),
		TieneClorosis: cmd.TieneClorosis,
		Confianza:     cmd.Confianza,
		ImageURL:      cmd.ImageURL,
		ImageBase64:   "", // Manual flow keeps it locally
		ProcesadoAt:   cmd.ProcesadoAt,
		CreatedAt:     diagnostico.CreatedAt(),
	}, nil
}

package registrarinferenciadesdenodo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	fincasdomain "github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	nodosdomain "github.com/davosjar/bunna/services/fincas/internal/nodos/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const routingKey = "nodos.v1.inferencia.desde_nodo"

type InferenciaDesdeNodo struct {
	EventID        string    `json:"event_id"`
	MuestraID      string    `json:"muestra_id"`
	DiagnosticoID  string    `json:"diagnostico_id"`
	NodoID         string    `json:"nodo_id"`
	LoteID         string    `json:"lote_id"`
	TieneClorosis  bool      `json:"tiene_clorosis"`
	Confianza      float64   `json:"confianza"`
	TenantID       string    `json:"tenant_id,omitempty"`
	OcurredAt      time.Time `json:"ocurred_at"`
}

type UseCase struct {
	nodoRepo       nodosdomain.NodoRepositorio
	loteRepo       fincasdomain.LoteRepositorio
	uow            application.UnitOfWorkDiagnostico
	generador      shared.GeneradorID
	publisher      application.EventPublisher
}

func NewUseCase(
	nodoRepo nodosdomain.NodoRepositorio,
	loteRepo fincasdomain.LoteRepositorio,
	uow application.UnitOfWorkDiagnostico,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		nodoRepo:  nodoRepo,
		loteRepo:  loteRepo,
		uow:       uow,
		generador: generador,
		publisher: publisher,
	}
}

func (uc *UseCase) Ejecutar(ctx context.Context, cmd Command) (*Salida, error) {
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	nodo, err := uc.nodoRepo.ObtenerPorID(ctx, cmd.NodoID)
	if err != nil {
		if errors.Is(err, nodosdomain.ErrNodoNoEncontrado) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	if !nodo.IsActivo() {
		return nil, application.ErrConflictoEstado("el nodo no está activo")
	}

	fincaID := cmd.FincaID
	if fincaID == "" {
		fincaID = nodo.FincaID()
	}
	if fincaID == "" {
		return nil, application.ErrValidacion("no se pudo determinar el fincaID: el nodo no tiene finca asignada y no se proporcionó en el comando")
	}

	loteID := cmd.LoteID
	if loteID == "" && nodo.LoteID() != nil {
		loteID = *nodo.LoteID()
	}

	if loteID != "" {
		_, err = uc.loteRepo.ObtenerPorID(ctx, loteID)
		if err != nil {
			if errors.Is(err, fincasdomain.ErrLoteNoEncontrado) {
				return nil, application.ErrNotFound
			}
			return nil, err
		}
	}

	tenantID := nodo.TenantID()

	var muestraID, diagnosticoID string
	var estado string

	err = uc.uow.Transaccional(ctx, func(tx application.UnitOfWorkDiagnostico) error {
		idMuestra, err := uc.generador.NextID(ctx)
		if err != nil {
			return err
		}

		muestra, err := diagnosticodomain.NewMuestraSinUbicacion(idMuestra, fincaID, loteID, tenantID)
		if err != nil {
			return application.ErrValidacion(err.Error())
		}

		if err := tx.MuestraRepo().Crear(ctx, muestra); err != nil {
			return err
		}

		idDiagnostico, err := uc.generador.NextID(ctx)
		if err != nil {
			return err
		}

		resultado, err := diagnosticodomain.NewResultadoInferencia(
			cmd.ImageURL, cmd.ImageBase64, cmd.TieneClorosis, cmd.Confianza, cmd.ProcesadoAt,
		)
		if err != nil {
			return application.ErrValidacion(err.Error())
		}

		nombre := fmt.Sprintf("NOD-%s-%04d", time.Now().Format("20060102"), rand.Intn(10000))

		diagnostico, err := diagnosticodomain.NewDiagnostico(idDiagnostico, nombre, idMuestra, tenantID, resultado)
		if err != nil {
			return application.ErrValidacion(err.Error())
		}

		if err := tx.DiagnosticoRepo().Crear(ctx, diagnostico); err != nil {
			return err
		}

		muestraID = idMuestra
		diagnosticoID = idDiagnostico
		estado = string(diagnostico.Estado())

		return nil
	})
	if err != nil {
		return nil, err
	}

	eventID, _ := uc.generador.NextID(ctx)
	evento := InferenciaDesdeNodo{
		EventID:       eventID,
		MuestraID:     muestraID,
		DiagnosticoID: diagnosticoID,
		NodoID:        cmd.NodoID,
		LoteID:        loteID,
		TieneClorosis: cmd.TieneClorosis,
		Confianza:     cmd.Confianza,
		TenantID:      tenantID,
		OcurredAt:     time.Now(),
	}
	_ = uc.publisher.Publish(ctx, routingKey, evento)

	return &Salida{
		MuestraID:     muestraID,
		DiagnosticoID: diagnosticoID,
		Estado:        estado,
		TieneClorosis: cmd.TieneClorosis,
		Confianza:     cmd.Confianza,
		ImageURL:      cmd.ImageURL,
		ImageBase64:   cmd.ImageBase64,
		CreatedAt:     time.Now(),
	}, nil
}

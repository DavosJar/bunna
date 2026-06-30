package facades

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/aceptardiagnostico"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/rechazardiagnostico"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarinferencia"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/solicitardiagnosticomanual"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/mapper"
)

type (
	SolicitarDiagnosticoManualUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd solicitardiagnosticomanual.Command) (*solicitardiagnosticomanual.Salida, error)
	}
	AceptarDiagnosticoUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd aceptardiagnostico.Command) (*aceptardiagnostico.Salida, error)
	}
	RegistrarInferenciaUseCase interface {
		Ejecutar(ctx context.Context, cmd registrarinferencia.Command) (*registrarinferencia.Salida, error)
	}
	RechazarDiagnosticoUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd rechazardiagnostico.Command) (*rechazardiagnostico.Salida, error)
	}
)

type DiagnosticosFacade interface {
	SolicitarManual(ctx context.Context, auth *application.AuthContext, muestraID string, req dto.SolicitarDiagnosticoManualRequest) (*dto.SolicitudDiagnosticoResponse, error)
	Aceptar(ctx context.Context, auth *application.AuthContext, diagnosticoID string) (*dto.EstadoCambioResponse, error)
	Rechazar(ctx context.Context, auth *application.AuthContext, diagnosticoID string, req dto.RechazarDiagnosticoRequest) (*dto.EstadoCambioResponse, error)
	GuardarResultadoManual(ctx context.Context, auth *application.AuthContext, muestraID string, req dto.GuardarResultadoManualRequest) (*dto.DiagnosticoResponse, error)
}

type diagnosticosFacade struct {
	solicitarUC SolicitarDiagnosticoManualUseCase
	aceptarUC   AceptarDiagnosticoUseCase
	rechazarUC  RechazarDiagnosticoUseCase
	inferenciaUC RegistrarInferenciaUseCase
	mapper      mapper.DiagnosticoMapper
}

func NewDiagnosticosFacade(
	solicitarUC SolicitarDiagnosticoManualUseCase,
	aceptarUC AceptarDiagnosticoUseCase,
	rechazarUC RechazarDiagnosticoUseCase,
	inferenciaUC RegistrarInferenciaUseCase,
) DiagnosticosFacade {
	return &diagnosticosFacade{
		solicitarUC: solicitarUC,
		aceptarUC:   aceptarUC,
		rechazarUC:  rechazarUC,
		inferenciaUC: inferenciaUC,
		mapper:      mapper.DiagnosticoMapper{},
	}
}

func (f *diagnosticosFacade) SolicitarManual(ctx context.Context, auth *application.AuthContext, muestraID string, req dto.SolicitarDiagnosticoManualRequest) (*dto.SolicitudDiagnosticoResponse, error) {
	cmd := solicitardiagnosticomanual.Command{
		MuestraID: muestraID,
		ImageURL:  req.ImageURL,
	}
	salida, err := f.solicitarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.SolicitarManualSalidaToResponse(salida)
	return &resp, nil
}

func (f *diagnosticosFacade) Aceptar(ctx context.Context, auth *application.AuthContext, diagnosticoID string) (*dto.EstadoCambioResponse, error) {
	cmd := aceptardiagnostico.Command{DiagnosticoID: diagnosticoID}
	salida, err := f.aceptarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.AceptarSalidaToEstadoCambio(salida)
	return &resp, nil
}

func (f *diagnosticosFacade) Rechazar(ctx context.Context, auth *application.AuthContext, diagnosticoID string, req dto.RechazarDiagnosticoRequest) (*dto.EstadoCambioResponse, error) {
	cmd := rechazardiagnostico.Command{
		DiagnosticoID: diagnosticoID,
		Motivo:        req.Motivo,
	}
	salida, err := f.rechazarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.RechazarSalidaToEstadoCambio(salida)
	return &resp, nil
}

func (f *diagnosticosFacade) GuardarResultadoManual(ctx context.Context, auth *application.AuthContext, muestraID string, req dto.GuardarResultadoManualRequest) (*dto.DiagnosticoResponse, error) {
	cmd := registrarinferencia.Command{
		MuestraID:     muestraID,
		ImageURL:      req.ImageURL,
		TieneClorosis: req.TieneClorosis,
		Confianza:     req.Confianza,
		ProcesadoAt:   req.ProcesadoAt,
	}
	salida, err := f.inferenciaUC.Ejecutar(ctx, cmd)
	if err != nil {
		return nil, err
	}
	resp := dto.DiagnosticoResponse{
		ID:            salida.ID,
		MuestraID:     salida.MuestraID,
		Nombre:        salida.Nombre,
		Estado:        salida.Estado,
		TieneClorosis: salida.TieneClorosis,
		Confianza:     salida.Confianza,
		ImageURL:      salida.ImageURL,
		ProcesadoAt:   salida.ProcesadoAt,
		CreatedAt:     salida.CreatedAt,
	}
	return &resp, nil
}

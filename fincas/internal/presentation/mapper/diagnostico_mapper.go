package mapper

import (
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/aceptardiagnostico"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/rechazardiagnostico"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarinferencia"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/solicitardiagnosticomanual"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
)

// DiagnosticoMapper convierte entre capa de aplicación y DTOs de diagnóstico.
type DiagnosticoMapper struct{}

func (DiagnosticoMapper) SolicitarManualSalidaToResponse(s *solicitardiagnosticomanual.Salida) dto.SolicitudDiagnosticoResponse {
	return dto.SolicitudDiagnosticoResponse{
		SolicitudID:  s.SolicitudID,
		MuestraID:    s.MuestraID,
		SolicitadoEn: s.SolicitadoEn,
	}
}

func (DiagnosticoMapper) AceptarSalidaToEstadoCambio(s *aceptardiagnostico.Salida) dto.EstadoCambioResponse {
	return dto.EstadoCambioResponse{
		ID:        s.ID,
		Estado:    s.Estado,
		UpdatedAt: s.UpdatedAt,
	}
}

func (DiagnosticoMapper) RechazarSalidaToEstadoCambio(s *rechazardiagnostico.Salida) dto.EstadoCambioResponse {
	return dto.EstadoCambioResponse{
		ID:        s.ID,
		Estado:    s.Estado,
		Motivo:    s.Motivo,
		UpdatedAt: s.UpdatedAt,
	}
}

func (DiagnosticoMapper) InferenciaSalidaToResponse(s *registrarinferencia.Salida) dto.DiagnosticoResponse {
	return dto.DiagnosticoResponse{
		ID:            s.ID,
		MuestraID:     s.MuestraID,
		Nombre:        s.Nombre,
		Estado:        s.Estado,
		TieneClorosis: s.TieneClorosis,
		Confianza:     s.Confianza,
		ImageURL:      s.ImageURL,
		ProcesadoAt:   s.ProcesadoAt,
		CreatedAt:     s.CreatedAt,
	}
}

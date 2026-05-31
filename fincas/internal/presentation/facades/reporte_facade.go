package facades

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/generarreporteporlote"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/mapper"
)

type GenerarReportePorLoteUseCase interface {
	Ejecutar(ctx context.Context, auth *application.AuthContext, cmd generarreporteporlote.Command) (*generarreporteporlote.Salida, error)
}

type ReportesFacade interface {
	GenerarPorLote(ctx context.Context, auth *application.AuthContext, loteID string) (*dto.ReporteLoteResponse, error)
}

type reportesFacade struct {
	reporteUC GenerarReportePorLoteUseCase
	mapper    mapper.ReporteMapper
}

func NewReportesFacade(reporteUC GenerarReportePorLoteUseCase) ReportesFacade {
	return &reportesFacade{
		reporteUC: reporteUC,
		mapper:    mapper.ReporteMapper{},
	}
}

func (f *reportesFacade) GenerarPorLote(ctx context.Context, auth *application.AuthContext, loteID string) (*dto.ReporteLoteResponse, error) {
	cmd := generarreporteporlote.Command{LoteID: loteID}
	salida, err := f.reporteUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}

	resp := f.mapper.SalidaToResponse(salida)
	return &resp, nil
}

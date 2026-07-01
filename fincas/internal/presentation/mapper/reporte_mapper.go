package mapper

import (
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/generarreporteporlote"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
)

// ReporteMapper convierte entre capa de aplicación y DTOs de reporte.
type ReporteMapper struct{}

func (ReporteMapper) SalidaToResponse(s *generarreporteporlote.Salida) dto.ReporteLoteResponse {
	muestras := make([]dto.MuestraReporteItem, len(s.Muestras))
	for i, m := range s.Muestras {
		muestras[i] = dto.MuestraReporteItem{
			ID:                m.ID,
			Latitud:           m.Latitud,
			Longitud:          m.Longitud,
			DiagnosticoID:     m.DiagnosticoID,
			EstadoDiagnostico: m.EstadoDiagnostico,
			ImageURL:          m.ImageURL,
			ImageBase64:       m.ImageBase64,
			TieneClorosis:     m.TieneClorosis,
			Confianza:         m.Confianza,
		}
	}

	zonas := make([]dto.ZonaAfectadaDTO, len(s.Zonas))
	for i, z := range s.Zonas {
		zonas[i] = dto.ZonaAfectadaDTO{
			Latitud:  z.Latitud,
			Longitud: z.Longitud,
			RadioMts: z.RadioMts,
		}
	}

	return dto.ReporteLoteResponse{
		ID:        s.ID,
		Nombre:    s.Nombre,
		AreaTotal: s.AreaTotal,
		Estado:    s.Estado,
		Muestras:  muestras,
		Zonas:     zonas,
		Metricas: dto.MetricasDTO{
			TotalMuestras:        s.Metricas.TotalMuestras,
			ConClorosis:          s.Metricas.ConClorosis,
			SinClorosis:          s.Metricas.SinClorosis,
			Pendientes:           s.Metricas.Pendientes,
			AreaAfectadaEstimada: s.Metricas.AreaAfectadaEstimada,
			PorcentajeAfectado:   s.Metricas.PorcentajeAfectado,
		},
		GeneradoEn: s.GeneradoEn,
	}
}

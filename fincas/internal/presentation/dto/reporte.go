package dto

import "time"

// MuestraReporteItem representa una muestra individual en el reporte.
type MuestraReporteItem struct {
	ID                string   `json:"id"`
	Latitud           float64  `json:"latitud"`
	Longitud          float64  `json:"longitud"`
	DiagnosticoID     string   `json:"diagnosticoID,omitempty"`
	EstadoDiagnostico string   `json:"estadoDiagnostico,omitempty"`
	ImageURL          string   `json:"imageURL,omitempty"`
	ImageBase64       string   `json:"imageBase64,omitempty"`
	TieneClorosis     *bool    `json:"tieneClorosis,omitempty"`
	Confianza         *float64 `json:"confianza,omitempty"`
}

// ZonaAfectadaDTO representa un punto con clorosis en el lote.
type ZonaAfectadaDTO struct {
	Latitud  float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
	RadioMts float64 `json:"radioMts"`
}

// MetricasDTO contiene los indicadores calculados del reporte.
type MetricasDTO struct {
	TotalMuestras        int     `json:"totalMuestras"`
	ConClorosis          int     `json:"conClorosis"`
	SinClorosis          int     `json:"sinClorosis"`
	Pendientes           int     `json:"pendientes"`
	AreaAfectadaEstimada float64 `json:"areaAfectadaEstimada"`
	PorcentajeAfectado   float64 `json:"porcentajeAfectado"`
}

// ReporteLoteResponse es la respuesta completa del reporte de un lote.
type ReporteLoteResponse struct {
	ID         string               `json:"id"`
	Nombre     string               `json:"nombre"`
	AreaTotal  float64              `json:"areaTotal"`
	Estado     string               `json:"estado"`
	Muestras   []MuestraReporteItem `json:"muestras"`
	Zonas      []ZonaAfectadaDTO    `json:"zonas"`
	Metricas   MetricasDTO          `json:"metricas"`
	GeneradoEn time.Time            `json:"generadoEn"`
}

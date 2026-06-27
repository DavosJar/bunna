package generarreporteporlote

import "time"

// MuestraReporte representa una muestra individual dentro del reporte de un lote.
// Los campos DiagnosticoID, EstadoDiagnostico, TieneClorosis y Confianza son
// opcionales: se llenan solo si existe un diagnóstico para la muestra.
// TieneClorosis y Confianza usan punteros para distinguir entre "no hay diagnóstico"
// (nil) y "el diagnóstico indica falso/cero" (false/0).
type MuestraReporte struct {
	ID                string
	Latitud           float64
	Longitud          float64
	DiagnosticoID     string   // opcional, solo si existe diagnóstico
	EstadoDiagnostico string   // opcional
	ImageURL          string   // opcional, URL de la imagen en YOLO
	TieneClorosis     *bool    // puntero, nil si no hay diagnóstico
	Confianza         *float64 // puntero, nil si no hay diagnóstico
}

// ZonaAfectada representa un punto con clorosis dentro del lote.
// Se usa para renderizar un mapa de calor.
type ZonaAfectada struct {
	Latitud  float64
	Longitud float64
	RadioMts float64 // siempre 2.0
}

// Metricas contiene los indicadores calculados a partir de las muestras
// y diagnósticos aceptados del lote.
type Metricas struct {
	TotalMuestras        int
	ConClorosis          int
	SinClorosis          int
	Pendientes           int
	AreaAfectadaEstimada float64 // en metros cuadrados
	PorcentajeAfectado   float64 // entre 0 y 100
}

// Salida es la respuesta completa del caso de uso GenerarReportePorLote.
type Salida struct {
	ID        string
	Nombre    string
	AreaTotal float64
	Estado    string
	Muestras  []MuestraReporte
	Zonas     []ZonaAfectada
	Metricas  Metricas
	GeneradoEn time.Time
}

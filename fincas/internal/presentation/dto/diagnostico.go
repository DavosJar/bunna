package dto

import "time"

// SolicitarDiagnosticoManualRequest es el body para solicitar diagnóstico manual.
type SolicitarDiagnosticoManualRequest struct {
	ImageURL string `json:"imageURL"`
}

// RechazarDiagnosticoRequest es el body para rechazar un diagnóstico.
type RechazarDiagnosticoRequest struct {
	Motivo string `json:"motivo"`
}

// DiagnosticoResponse es la respuesta completa de un diagnóstico.
type DiagnosticoResponse struct {
	ID            string    `json:"id"`
	MuestraID     string    `json:"muestraID"`
	Nombre        string    `json:"nombre"`
	Estado        string    `json:"estado"`
	TieneClorosis bool      `json:"tieneClorosis"`
	Confianza     float64   `json:"confianza"`
	ImageURL      string    `json:"imageURL"`
	ProcesadoAt   time.Time `json:"procesadoAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

// SolicitudDiagnosticoResponse es la respuesta al solicitar un diagnóstico manual.
type SolicitudDiagnosticoResponse struct {
	SolicitudID  string    `json:"solicitudID"`
	MuestraID    string    `json:"muestraID"`
	SolicitadoEn time.Time `json:"solicitadoEn"`
}

// EstadoCambioResponse es la respuesta genérica para cambios de estado
// (desactivar finca, eliminar lote, aceptar/rechazar diagnóstico).
type EstadoCambioResponse struct {
	ID        string    `json:"id"`
	Estado    string    `json:"estado"`
	Motivo    string    `json:"motivo,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

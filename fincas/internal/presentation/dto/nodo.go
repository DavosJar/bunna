package dto

import "time"

type RegistrarNodoRequest struct {
	FincaID string  `json:"finca_id"`
	NodeKey string  `json:"node_key"`
	LoteID  *string `json:"lote_id,omitempty"`
	Nombre  string  `json:"nombre"`
}

type EditarNodoRequest struct {
	LoteID *string `json:"lote_id"`
	Nombre string  `json:"nombre"`
}

type DesactivarNodoRequest struct {
	Estado string `json:"estado"`
}

type NodoResponse struct {
	ID            string    `json:"id"`
	FincaID       string    `json:"finca_id"`
	LoteID        *string   `json:"lote_id,omitempty"`
	TenantID      string    `json:"tenant_id"`
	Nombre        string    `json:"nombre"`
	NodeKey       string    `json:"node_key"`
	Estado        string    `json:"estado"`
	CreadoEn      time.Time `json:"creado_en"`
	ActualizadoEn time.Time `json:"actualizado_en,omitempty"`
}

type ValidarNodoResponse struct {
	NodoID   string  `json:"nodoID"`
	FincaID  string  `json:"fincaID"`
	LoteID   *string `json:"loteID,omitempty"`
	TenantID string  `json:"tenantID"`
}

type RegistrarInferenciaDesdeNodoRequest struct {
	NodoID        string    `json:"nodoID"`
	FincaID       string    `json:"fincaID"`
	LoteID        string    `json:"loteID"`
	TenantID      string    `json:"tenantID"`
	ImageURL      string    `json:"imageURL"`
	TieneClorosis bool      `json:"tieneClorosis"`
	Confianza     float64   `json:"confianza"`
	ProcesadoAt   time.Time `json:"procesadoAt"`
}

type InferenciaResponse struct {
	MuestraID     string    `json:"muestra_id"`
	DiagnosticoID string    `json:"diagnostico_id"`
	Estado        string    `json:"estado"`
	TieneClorosis bool      `json:"tieneClorosis"`
	Confianza     float64   `json:"confianza"`
	ImageURL      string    `json:"imageURL"`
	CreatedAt     time.Time `json:"created_at"`
}

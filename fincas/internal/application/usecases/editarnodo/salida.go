package editarnodo

import "time"

type Salida struct {
	ID            string    `json:"id"`
	FincaID       string    `json:"finca_id"`
	LoteID        *string   `json:"lote_id,omitempty"`
	TenantID      string    `json:"tenant_id"`
	Nombre        string    `json:"nombre"`
	NodeKey       string    `json:"node_key"`
	Estado        string    `json:"estado"`
	ActualizadoEn time.Time `json:"actualizado_en"`
}

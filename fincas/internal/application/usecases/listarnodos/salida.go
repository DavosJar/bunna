package listarnodos

import "time"

type NodoItem struct {
	ID        string    `json:"id"`
	FincaID   string    `json:"finca_id"`
	LoteID    *string   `json:"lote_id,omitempty"`
	TenantID  string    `json:"tenant_id"`
	Nombre    string    `json:"nombre"`
	NodeKey   string    `json:"node_key"`
	Estado    string    `json:"estado"`
	CreadoEn  time.Time `json:"creado_en"`
}

type Salida struct {
	Nodos        []NodoItem
	Total        int
	Pagina       int
	TotalPaginas int
}

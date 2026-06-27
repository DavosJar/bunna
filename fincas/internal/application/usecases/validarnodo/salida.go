package validarnodo

type Salida struct {
	NodoID   string  `json:"nodoID"`
	FincaID  string  `json:"fincaID"`
	LoteID   *string `json:"loteID,omitempty"`
	TenantID string  `json:"tenantID"`
}

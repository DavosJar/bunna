package dto

import "time"

// AgregarLoteRequest es el body para agregar un lote.
type AgregarLoteRequest struct {
	Nombre      string  `json:"nombre"`
	Area        float64 `json:"area"`
	Descripcion string  `json:"descripcion"`
}

// LoteResponse es la respuesta de un lote (agregar, listar).
type LoteResponse struct {
	ID          string    `json:"id"`
	FincaID     string    `json:"fincaID"`
	Nombre      string    `json:"nombre"`
	Area        float64   `json:"area"`
	Descripcion string    `json:"descripcion"`
	Estado      string    `json:"estado"`
	CreatedAt   time.Time `json:"createdAt"`
}

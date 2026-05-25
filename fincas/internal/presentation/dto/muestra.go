package dto

import "time"

// TomarMuestraRequest es el body para tomar una muestra.
type TomarMuestraRequest struct {
	Latitud  float64 `json:"latitud"`
	Longitud float64 `json:"longitud"`
}

// MuestraResponse es la respuesta al tomar una muestra.
type MuestraResponse struct {
	ID        string    `json:"id"`
	LoteID    string    `json:"loteID"`
	Latitud   float64   `json:"latitud"`
	Longitud  float64   `json:"longitud"`
	CreatedAt time.Time `json:"createdAt"`
}

// MuestraItemResponse es un item en el listado de muestras por lote.
type MuestraItemResponse struct {
	ID        string    `json:"id"`
	LoteID    string    `json:"loteID"`
	Latitud   float64   `json:"latitud"`
	Longitud  float64   `json:"longitud"`
	CreatedAt time.Time `json:"createdAt"`
}

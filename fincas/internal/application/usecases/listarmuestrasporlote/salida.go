package listarmuestrasporlote

import "time"

// MuestraItem representa una muestra individual en la respuesta del listado por lote.
type MuestraItem struct {
	ID        string
	LoteID    string
	Latitud   float64
	Longitud  float64
	CreatedAt time.Time
}

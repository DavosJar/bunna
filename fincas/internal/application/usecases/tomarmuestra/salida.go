package tomarmuestra

import "time"

// Salida contiene los datos de respuesta tras tomar una muestra exitosamente.
type Salida struct {
	ID        string
	FincaID   string
	LoteID    string
	Latitud   float64
	Longitud  float64
	CreatedAt time.Time
}

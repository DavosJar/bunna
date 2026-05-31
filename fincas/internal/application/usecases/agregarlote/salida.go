package agregarlote

import "time"

// Salida contiene los datos de respuesta tras agregar un lote exitosamente.
type Salida struct {
	ID          string
	FincaID     string
	Nombre      string
	Area        float64
	Descripcion string
	Estado      string
	CreatedAt   time.Time
}

package eliminarlote

import "time"

// Salida contiene los datos de respuesta tras eliminar un lote.
type Salida struct {
	ID        string
	Estado    string
	UpdatedAt time.Time
}

package desactivarfinca

import "time"

// Salida contiene los datos de respuesta tras desactivar una finca.
type Salida struct {
	ID        string
	Estado    string
	UpdatedAt time.Time
}

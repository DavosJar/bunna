package registrarfinca

import "time"

// Salida contiene los datos de respuesta tras registrar una finca exitosamente.
type Salida struct {
	ID          string
	Nombre      string
	Ubicacion   string
	Descripcion string
	Estado      string
	CreatedAt   time.Time
}

package aceptardiagnostico

import "time"

// Salida contiene los datos de respuesta tras aceptar un diagnóstico exitosamente.
type Salida struct {
	ID        string
	Estado    string
	UpdatedAt time.Time
}

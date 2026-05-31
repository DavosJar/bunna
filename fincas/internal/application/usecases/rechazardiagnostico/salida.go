package rechazardiagnostico

import "time"

// Salida contiene los datos de respuesta tras rechazar un diagnóstico exitosamente.
type Salida struct {
	ID        string
	Estado    string
	Motivo    string
	UpdatedAt time.Time
}

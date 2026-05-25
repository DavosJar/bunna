package solicitardiagnosticomanual

import "time"

// Salida contiene los datos de respuesta tras solicitar un diagnóstico manual exitosamente.
type Salida struct {
	SolicitudID  string
	MuestraID    string
	SolicitadoEn time.Time
}

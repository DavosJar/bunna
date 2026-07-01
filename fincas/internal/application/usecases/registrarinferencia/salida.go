package registrarinferencia

import "time"

// Salida contiene los datos de respuesta tras registrar una inferencia exitosamente.
type Salida struct {
	ID             string
	MuestraID      string
	Nombre         string
	Estado         string
	TieneClorosis  bool
	Confianza      float64
	ImageURL       string
	ImageBase64    string
	ProcesadoAt    time.Time
	CreatedAt      time.Time
}

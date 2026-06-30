package listarfincas

import "time"

type FincaSalida struct {
	ID          string
	Nombre      string
	Ubicacion   string
	Descripcion string
	Estado      string
	CreatedAt   time.Time
}

type Salida struct {
	Fincas []FincaSalida
}

package listarlotes

import "time"

type LoteSalida struct {
	ID          string
	FincaID     string
	Nombre      string
	Area        float64
	Descripcion string
	Estado      string
	CreatedAt   time.Time
}

type Salida struct {
	Lotes []LoteSalida
}

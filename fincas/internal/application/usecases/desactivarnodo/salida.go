package desactivarnodo

import "time"

type Salida struct {
	ID            string    `json:"id"`
	Estado        string    `json:"estado"`
	ActualizadoEn time.Time `json:"actualizado_en"`
}

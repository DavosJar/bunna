package registrarinferenciadesdenodo

import "time"

type Salida struct {
	MuestraID     string    `json:"muestra_id"`
	DiagnosticoID string    `json:"diagnostico_id"`
	Estado        string    `json:"estado"`
	TieneClorosis bool      `json:"tieneClorosis"`
	Confianza     float64   `json:"confianza"`
	ImageURL      string    `json:"imageURL"`
	CreatedAt     time.Time `json:"created_at"`
}

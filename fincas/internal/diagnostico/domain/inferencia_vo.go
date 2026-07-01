package domain

import (
	"time"
)

type ResultadoInferencia struct {
	imageUrl      string
	imageBase64   string
	tieneClorosis bool
	confianza     float64
	procesadoAt   time.Time
}

func NewResultadoInferencia(imageUrl string, imageBase64 string, tieneClorosis bool, confianza float64, procesadoAt time.Time) (*ResultadoInferencia, error) {
	if imageUrl == "" {
		return nil, ErrImageUrlRequerida
	}
	if confianza < 0 || confianza > 1 {
		return nil, ErrConfianzaInvalida
	}
	return &ResultadoInferencia{
		imageUrl:      imageUrl,
		imageBase64:   imageBase64,
		tieneClorosis: tieneClorosis,
		confianza:     confianza,
		procesadoAt:   procesadoAt,
	}, nil
}

func (r *ResultadoInferencia) ImageUrl() string       { return r.imageUrl }
func (r *ResultadoInferencia) ImageBase64() string    { return r.imageBase64 }
func (r *ResultadoInferencia) TieneClorosis() bool    { return r.tieneClorosis }
func (r *ResultadoInferencia) Confianza() float64     { return r.confianza }
func (r *ResultadoInferencia) ProcesadoAt() time.Time { return r.procesadoAt }

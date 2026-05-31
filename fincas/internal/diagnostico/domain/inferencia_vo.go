package domain

import (
	"time"
)

type ResultadoInferencia struct {
	imageUrl      string
	tieneClorosis bool
	confianza     float64
	procesadoAt   time.Time
}

func NewResultadoInferencia(imageUrl string, tieneClorosis bool, confianza float64, procesadoAt time.Time) (*ResultadoInferencia, error) {
	if imageUrl == "" {
		return nil, ErrImageUrlRequerida
	}
	if confianza < 0 || confianza > 1 {
		return nil, ErrConfianzaInvalida
	}
	return &ResultadoInferencia{
		imageUrl:      imageUrl,
		tieneClorosis: tieneClorosis,
		confianza:     confianza,
		procesadoAt:   procesadoAt,
	}, nil
}

func (r *ResultadoInferencia) ImageUrl() string       { return r.imageUrl }
func (r *ResultadoInferencia) TieneClorosis() bool    { return r.tieneClorosis }
func (r *ResultadoInferencia) Confianza() float64     { return r.confianza }
func (r *ResultadoInferencia) ProcesadoAt() time.Time { return r.procesadoAt }

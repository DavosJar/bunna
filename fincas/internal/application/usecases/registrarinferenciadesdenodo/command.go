package registrarinferenciadesdenodo

import (
	"fmt"
	"time"
)

type Command struct {
	NodoID        string    `json:"nodoID"`
	FincaID       string    `json:"fincaID"`
	LoteID        string    `json:"loteID"`
	TenantID      string    `json:"tenantID"`
	ImageURL      string    `json:"imageURL"`
	TieneClorosis bool      `json:"tieneClorosis"`
	Confianza     float64   `json:"confianza"`
	ProcesadoAt   time.Time `json:"procesadoAt"`
}

func (c *Command) Validar() error {
	if c.NodoID == "" {
		return fmt.Errorf("validación: el nodoID es requerido")
	}
	if c.ImageURL == "" {
		return fmt.Errorf("validación: la imageURL es requerida")
	}
	if c.Confianza < 0 || c.Confianza > 1 {
		return fmt.Errorf("validación: la confianza debe estar entre 0 y 1")
	}
	if c.ProcesadoAt.IsZero() {
		return fmt.Errorf("validación: el procesadoAt es requerido")
	}
	if c.ProcesadoAt.After(time.Now().Add(5 * time.Minute)) {
		return fmt.Errorf("validación: el procesadoAt no puede ser una fecha futura")
	}
	return nil
}

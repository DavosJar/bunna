package registrarinferencia

import (
	"fmt"
	"time"
)

// Command contiene los datos de entrada para registrar una inferencia.
// Es un comando interno disparado por el consumer de RabbitMQ al recibir
// el resultado de YOLO.
type Command struct {
	MuestraID     string
	ImageURL      string
	TieneClorosis bool
	Confianza     float64
	ProcesadoAt   time.Time
}

// Validar verifica que todos los campos cumplen las restricciones de formato.
func (c *Command) Validar() error {
	if c.MuestraID == "" {
		return fmt.Errorf("validación: el muestraID es requerido")
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
	if c.ProcesadoAt.After(time.Now()) {
		return fmt.Errorf("validación: el procesadoAt no puede ser una fecha futura")
	}
	return nil
}

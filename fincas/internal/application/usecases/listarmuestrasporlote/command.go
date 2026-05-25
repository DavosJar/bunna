package listarmuestrasporlote

import "fmt"

// Command contiene los datos de entrada para listar muestras de un lote.
type Command struct {
	LoteID string
}

// Validar verifica que el campo LoteID cumple las restricciones.
func (c *Command) Validar() error {
	if c.LoteID == "" {
		return fmt.Errorf("validación: el loteID es requerido")
	}
	return nil
}

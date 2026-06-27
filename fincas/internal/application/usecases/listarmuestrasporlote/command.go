package listarmuestrasporlote

import "fmt"

// Command contiene los datos de entrada para listar muestras de un lote.
type Command struct {
	FincaID string
	LoteID  string
}

// Validar verifica que el campo LoteID cumple las restricciones.
func (c *Command) Validar() error {
	if c.FincaID == "" {
		return fmt.Errorf("validación: el fincaID es requerido")
	}
	return nil
}

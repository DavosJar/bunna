package generarreporteporlote

import "fmt"

// Command contiene los datos de entrada para generar el reporte de un lote.
type Command struct {
	LoteID string
}

// Validar verifica que los campos obligatorios estén presentes.
func (c *Command) Validar() error {
	if c.LoteID == "" {
		return fmt.Errorf("validación: el loteID es requerido")
	}
	return nil
}

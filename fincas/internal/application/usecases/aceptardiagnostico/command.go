package aceptardiagnostico

import "fmt"

// Command contiene los datos de entrada para aceptar un diagnóstico.
type Command struct {
	DiagnosticoID string
}

// Validar verifica que los campos obligatorios estén presentes.
func (c *Command) Validar() error {
	if c.DiagnosticoID == "" {
		return fmt.Errorf("validación: el diagnosticoID es requerido")
	}
	return nil
}

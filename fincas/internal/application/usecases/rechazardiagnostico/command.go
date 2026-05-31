package rechazardiagnostico

import "fmt"

// Command contiene los datos de entrada para rechazar un diagnóstico.
type Command struct {
	DiagnosticoID string
	Motivo        string
}

// Validar verifica que todos los campos cumplen las restricciones de formato.
func (c *Command) Validar() error {
	if c.DiagnosticoID == "" {
		return fmt.Errorf("validación: el diagnosticoID es requerido")
	}
	if len(c.Motivo) > 500 {
		return fmt.Errorf("validación: el motivo no debe exceder 500 caracteres")
	}
	return nil
}

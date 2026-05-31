package registrarfinca

import "fmt"

// Command contiene los datos de entrada para registrar una finca nueva.
type Command struct {
	Nombre      string
	Ubicacion   string
	Descripcion string
}

// Validar verifica que todos los campos cumplen las restricciones de formato.
func (c *Command) Validar() error {
	if c.Nombre == "" || len(c.Nombre) < 3 {
		return fmt.Errorf("validación: el nombre debe tener al menos 3 caracteres")
	}
	if len(c.Nombre) > 200 {
		return fmt.Errorf("validación: el nombre no debe exceder 200 caracteres")
	}
	if c.Ubicacion == "" {
		return fmt.Errorf("validación: la ubicación es requerida")
	}
	if len(c.Ubicacion) > 500 {
		return fmt.Errorf("validación: la ubicación no debe exceder 500 caracteres")
	}
	if len(c.Descripcion) > 1000 {
		return fmt.Errorf("validación: la descripción no debe exceder 1000 caracteres")
	}
	return nil
}

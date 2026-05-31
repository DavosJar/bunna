package desactivarfinca

import "fmt"

// Command contiene los datos de entrada para desactivar una finca.
type Command struct {
	FincaID   string
	Confirmar bool
}

// Validar verifica que los campos obligatorios estén presentes.
func (c *Command) Validar() error {
	if c.FincaID == "" {
		return fmt.Errorf("validación: el fincaID es requerido")
	}
	return nil
}

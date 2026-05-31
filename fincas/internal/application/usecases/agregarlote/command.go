package agregarlote

import "fmt"

// Command contiene los datos de entrada para agregar un lote a una finca.
type Command struct {
	FincaID     string
	Nombre      string
	Area        float64
	Descripcion string
}

// Validar verifica que todos los campos cumplen las restricciones de formato.
func (c *Command) Validar() error {
	if c.FincaID == "" {
		return fmt.Errorf("validación: el fincaID es requerido")
	}
	if c.Nombre == "" || len(c.Nombre) < 3 {
		return fmt.Errorf("validación: el nombre debe tener al menos 3 caracteres")
	}
	if len(c.Nombre) > 150 {
		return fmt.Errorf("validación: el nombre no debe exceder 150 caracteres")
	}
	if c.Area <= 0 {
		return fmt.Errorf("validación: el área debe ser mayor a 0")
	}
	if c.Area > 99999.99 {
		return fmt.Errorf("validación: el área no debe exceder 99999.99")
	}
	if len(c.Descripcion) > 1000 {
		return fmt.Errorf("validación: la descripción no debe exceder 1000 caracteres")
	}
	return nil
}

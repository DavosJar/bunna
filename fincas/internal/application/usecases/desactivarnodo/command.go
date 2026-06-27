package desactivarnodo

import "fmt"

type Command struct {
	NodoID string
	Estado string
}

func (c *Command) Validar() error {
	if c.NodoID == "" {
		return fmt.Errorf("validación: el nodoID es requerido")
	}
	estadosValidos := map[string]bool{"INACTIVO": true, "MANTENIMIENTO": true}
	if !estadosValidos[c.Estado] {
		return fmt.Errorf("validación: estado debe ser INACTIVO o MANTENIMIENTO")
	}
	return nil
}

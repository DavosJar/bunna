package editarnodo

import "fmt"

type Command struct {
	NodoID string
	LoteID *string
	Nombre string
}

func (c *Command) Validar() error {
	if c.NodoID == "" {
		return fmt.Errorf("validación: el nodoID es requerido")
	}
	return nil
}

package obtenernodo

import "fmt"

type Command struct {
	NodoID string
}

func (c *Command) Validar() error {
	if c.NodoID == "" {
		return fmt.Errorf("validación: el nodoID es requerido")
	}
	return nil
}

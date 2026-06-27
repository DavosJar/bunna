package validarnodo

import "fmt"

type Command struct {
	NodeKey string
}

func (c *Command) Validar() error {
	if c.NodeKey == "" {
		return fmt.Errorf("validación: el nodeKey es requerido")
	}
	return nil
}

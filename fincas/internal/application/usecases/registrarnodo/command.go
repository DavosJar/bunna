package registrarnodo

import "fmt"

type Command struct {
	TenantID string
	FincaID  string
	NodeKey  string
	LoteID   *string
	Nombre   string
}

func (c *Command) Validar() error {
	if c.TenantID == "" {
		return fmt.Errorf("validación: el tenantID es requerido")
	}
	if c.FincaID == "" {
		return fmt.Errorf("validación: el fincaID es requerido")
	}
	if c.NodeKey == "" {
		return fmt.Errorf("validación: el nodeKey es requerido")
	}
	if c.Nombre == "" {
		return fmt.Errorf("validación: el nombre es requerido")
	}
	return nil
}

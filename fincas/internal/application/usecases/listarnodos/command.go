package listarnodos

import "fmt"

type Command struct {
	TenantID string
	Estado   string
}

func (c *Command) Validar() error {
	if c.TenantID == "" {
		return fmt.Errorf("validación: el tenantID es requerido")
	}
	return nil
}

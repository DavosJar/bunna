package solicitardiagnosticomanual

import (
	"fmt"
	"strings"
)

// Command contiene los datos de entrada para solicitar un diagnóstico manual.
type Command struct {
	MuestraID string
	ImageURL  string
}

// Validar verifica que todos los campos cumplen las restricciones de formato.
func (c *Command) Validar() error {
	if c.MuestraID == "" {
		return fmt.Errorf("validación: el muestraID es requerido")
	}
	if c.ImageURL == "" {
		return fmt.Errorf("validación: la imageURL es requerida")
	}
	if !strings.HasPrefix(c.ImageURL, "https://") {
		return fmt.Errorf("validación: la imageURL debe ser HTTPS")
	}
	return nil
}

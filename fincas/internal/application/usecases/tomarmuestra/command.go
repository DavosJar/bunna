package tomarmuestra

import "fmt"

// Command contiene los datos de entrada para tomar una muestra en un lote.
type Command struct {
	FincaID  string
	LoteID   string
	Latitud  float64
	Longitud float64
}

// Validar verifica que todos los campos cumplen las restricciones de formato.
func (c *Command) Validar() error {
	if c.FincaID == "" {
		return fmt.Errorf("validación: el fincaID es requerido")
	}
	if c.Latitud < -90 || c.Latitud > 90 {
		return fmt.Errorf("validación: la latitud debe estar entre -90 y 90")
	}
	if c.Longitud < -180 || c.Longitud > 180 {
		return fmt.Errorf("validación: la longitud debe estar entre -180 y 180")
	}
	return nil
}

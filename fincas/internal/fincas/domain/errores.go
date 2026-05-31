package domain

import "fmt"

// Errores de dominio — solo errores de NEGOCIO, no de formato
var (
	ErrNoPropietario               = fmt.Errorf("no tienes permisos sobre este recurso")
	ErrFincaNoEncontrada           = fmt.Errorf("finca no encontrada")
	ErrLoteNoEncontrado            = fmt.Errorf("lote no encontrado")
	ErrTransicionEstadoNoPermitida = fmt.Errorf("transición de estado no permitida")
)

func ErrFincaConLotes(count int) error {
	return fmt.Errorf("la finca tiene %d lote(s) asociado(s). confirma la eliminación", count)
}

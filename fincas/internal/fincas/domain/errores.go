package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNombreFincaRequerido = errors.New("El nombre de la finca es requerido y debe tener entre 3 y 200 caracteres")
	ErrNombreFincaLargo     = errors.New("El nombre de la finca no puede exceder 200 caracteres")
	ErrUbicacionRequerida   = errors.New("La ubicación de la finca es requerida")
	ErrUbicacionLarga       = errors.New("La ubicación no puede exceder 500 caracteres")
	ErrDescripcionLarga     = errors.New("La descripción no puede exceder 1000 caracteres")

	ErrNombreLoteRequerido = errors.New("El nombre del lote es requerido y debe tener entre 3 y 200 caracteres")
	ErrNombreLoteLargo     = errors.New("El nombre del lote no puede exceder 200 caracteres")
	ErrAreaRequerida       = errors.New("El área del lote es requerida y debe ser mayor a 0")

	ErrFincaNoEncontrada = errors.New("Finca no encontrada")
	ErrLoteNoEncontrado  = errors.New("Lote no encontrado")
	ErrNoPropietario     = errors.New("No tienes permisos sobre este recurso")
)

func ErrFincaConLotes(count int) error {
	return fmt.Errorf("La finca tiene %d lote(s) asociado(s). Confirma la eliminación con ?confirm=true", count)
}

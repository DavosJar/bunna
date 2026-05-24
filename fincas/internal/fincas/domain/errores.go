package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNombreFincaRequerido = errors.New("el nombre de la finca es requerido y debe tener entre 3 y 200 caracteres")
	ErrNombreFincaLargo     = errors.New("el nombre de la finca no puede exceder 200 caracteres")
	ErrUbicacionRequerida   = errors.New("la ubicación de la finca es requerida")
	ErrUbicacionLarga       = errors.New("la ubicación no puede exceder 500 caracteres")
	ErrDescripcionLarga     = errors.New("la descripción no puede exceder 1000 caracteres")
	ErrNoPropietario        = errors.New("el usuarioID de la finca es requerido")

	ErrFincaIDRequerido   = errors.New("el ID de la finca es requerido")
	ErrNombreLoteRequerido = errors.New("el nombre del lote es requerido y debe tener entre 3 y 200 caracteres")
	ErrNombreLoteLargo     = errors.New("el nombre del lote no puede exceder 200 caracteres")
	ErrAreaRequerida       = errors.New("el área del lote es requerida y debe ser mayor a 0")

	ErrFincaNoEncontrada = errors.New("finca no encontrada")
	ErrLoteNoEncontrado  = errors.New("lote no encontrado")
)

func ErrFincaConLotes(count int) error {
	return fmt.Errorf("la finca tiene %d lote(s) asociado(s). confirma la eliminación", count)
}

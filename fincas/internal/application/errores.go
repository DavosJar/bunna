package application

import "fmt"

// Errores de aplicación — mapean a códigos HTTP específicos.
// No son errores de dominio; son la traducción que hace la capa de aplicación
// de los errores de dominio a un lenguaje que la presentación entiende.
var (
	ErrForbidden = fmt.Errorf("no tiene permisos para esta operación")
	ErrNotFound  = fmt.Errorf("recurso no encontrado")
)

func ErrConflictoEstado(msg string) error {
	return fmt.Errorf("conflicto de estado: %s", msg)
}

func ErrValidacion(msg string) error {
	return fmt.Errorf("validación: %s", msg)
}

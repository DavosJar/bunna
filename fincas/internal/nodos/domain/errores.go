package domain

import "fmt"

var (
	ErrNodoNoEncontrado           = fmt.Errorf("nodo no encontrado")
	ErrNodoDuplicado              = fmt.Errorf("el nodeKey ya está registrado")
	ErrTransicionEstadoNoPermitida = fmt.Errorf("transición de estado no permitida")
)

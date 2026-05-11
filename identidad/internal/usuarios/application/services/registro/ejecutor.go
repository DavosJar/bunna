package registro

import "context"

// EjecutorRegistro define el contrato del caso de uso de registro.
// Permite mockear el servicio en pruebas de capas superiores.
type EjecutorRegistro interface {
	Ejecutar(ctx context.Context, comando *ComandoRegistro) (*DtoRespuestaRegistro, error)
}

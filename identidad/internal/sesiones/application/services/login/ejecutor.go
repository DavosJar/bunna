package login

import "context"

// EjecutorLogin define el contrato del caso de uso de login.
// Permite mockear el servicio en pruebas de capas superiores.
type EjecutorLogin interface {
	Ejecutar(ctx context.Context, cmd ComandoLogin) (*RespuestaLogin, error)
}

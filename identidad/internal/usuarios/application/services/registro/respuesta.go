package registro

import "time"

// DtoRespuestaRegistro es el DTO de salida para el caso de uso de registro
type DtoRespuestaRegistro struct {
	UsuarioID string
	Correo    string
	Estado    string
	Timestamp time.Time
}

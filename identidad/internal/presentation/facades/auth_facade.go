// Package facades agrupa los casos de uso de la capa de presentación.
// Las facades son la única puerta de entrada desde los handlers hacia la aplicación.
package facades

import "context"

// ComandoRegistro contiene los datos necesarios para registrar un nuevo usuario.
type ComandoRegistro struct {
	Nombre   string
	Apellido string
	Correo   string
	Password string
	Telefono string
}

// RespuestaRegistro contiene los datos del usuario recién creado.
type RespuestaRegistro struct {
	UsuarioID string
	Correo    string
	Estado    string
}

// ComandoLogin contiene las credenciales y el origen del intento de login.
type ComandoLogin struct {
	Email     string
	Password  string
	IPOrigen  string
}

// RespuestaLogin contiene los tokens generados tras un login exitoso.
type RespuestaLogin struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
	UsuarioID    string
	SesionID     string
}

// AuthFacade es la interfaz que expone los casos de uso de autenticación.
// Los handlers dependen de esta interfaz, nunca de la implementación concreta.
type AuthFacade interface {
	Registrar(ctx context.Context, cmd ComandoRegistro) (*RespuestaRegistro, error)
	Login(ctx context.Context, cmd ComandoLogin) (*RespuestaLogin, error)
}
// Package facades agrupa los casos de uso de la capa de presentación.
// Las facades son la única puerta de entrada desde los handlers hacia la aplicación.
package facades

import (
	"context"

	uc_sesiones_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	uc_sesiones_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	uc_sesiones_refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
)

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
	Email    string
	Password string
	IPOrigen string
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
	Refresh(ctx context.Context, cmd ComandoRefresh) (*RespuestaRefresh, error)
	Logout(ctx context.Context, cmd ComandoLogout) (*RespuestaLogout, error)
	LogoutAll(ctx context.Context, cmd ComandoLogoutAll) (*RespuestaLogout, error)
}

// ComandoRefresh contiene el token de refresco.
type ComandoRefresh struct {
	RefreshToken string
}

// RespuestaRefresh contiene los tokens renovados.
type RespuestaRefresh struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
	UsuarioID    string
}

// RespuestaLogout contiene la cantidad de sesiones revocadas.
type RespuestaLogout struct {
	SesionesRevocadas int
}

// ComandoLogout contiene los datos para cerrar una sesión.
type ComandoLogout struct {
	SesionID  string
	UsuarioID string
}

// ComandoLogoutAll contiene el usuarioID para cerrar todas las sesiones.
type ComandoLogoutAll struct {
	UsuarioID string
}

// LoginUseCase define el contrato del caso de uso de inicio de sesión.
// Satisfecho por uc_sesiones_login.IniciarSesionCasoDeUso.
type LoginUseCase interface {
	Ejecutar(ctx context.Context, cmd uc_sesiones_login.ComandoIniciarSesion) (*uc_sesiones_login.RespuestaIniciarSesion, error)
}

// RefreshUseCase define el contrato del caso de uso de renovación de sesión.
// Satisfecho por uc_sesiones_refresh.RenovarSesionCasoDeUso.
type RefreshUseCase interface {
	Ejecutar(ctx context.Context, cmd uc_sesiones_refresh.ComandoRenovarSesion) (*uc_sesiones_refresh.RespuestaRenovarSesion, error)
}

// LogoutUseCase define el contrato del caso de uso de cierre de sesión.
// Satisfecho por uc_sesiones_logout.CerrarSesionCasoDeUso.
type LogoutUseCase interface {
	Ejecutar(ctx context.Context, cmd uc_sesiones_logout.ComandoCerrarSesion) (*uc_sesiones_logout.RespuestaCerrarSesion, error)
	CerrarTodas(ctx context.Context, cmd uc_sesiones_logout.ComandoCerrarTodasLasSesiones) (*uc_sesiones_logout.RespuestaCerrarSesion, error)
}

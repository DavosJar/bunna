package facades

import (
	"context"
	"time"

	svc_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/login"
	svc_registro "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/services/registro"
)

// authFacadeImpl implementa AuthFacade orquestando los servicios de aplicación.
type authFacadeImpl struct {
	servicioRegistro svc_registro.EjecutorRegistro
	servicioLogin    svc_login.EjecutorLogin
}

// NewAuthFacade construye la implementación concreta de AuthFacade.
func NewAuthFacade(
	servicioRegistro svc_registro.EjecutorRegistro,
	servicioLogin svc_login.EjecutorLogin,
) AuthFacade {
	return &authFacadeImpl{
		servicioRegistro: servicioRegistro,
		servicioLogin:    servicioLogin,
	}
}

// Registrar delega al ServicioRegistro y traduce DTOs.
func (f *authFacadeImpl) Registrar(ctx context.Context, cmd ComandoRegistro) (*RespuestaRegistro, error) {
	respuesta, err := f.servicioRegistro.Ejecutar(ctx, &svc_registro.ComandoRegistro{
		Correo:   cmd.Correo,
		Password: cmd.Password,
		Nombre:   cmd.Nombre,
		Apellido: cmd.Apellido,
		Telefono: cmd.Telefono,
	})
	if err != nil {
		return nil, err
	}

	return &RespuestaRegistro{
		UsuarioID: respuesta.UsuarioID,
		Correo:    respuesta.Correo,
		Estado:    respuesta.Estado,
	}, nil
}

// Login delega al ServicioLogin y traduce DTOs.
func (f *authFacadeImpl) Login(ctx context.Context, cmd ComandoLogin) (*RespuestaLogin, error) {
	respuesta, err := f.servicioLogin.Ejecutar(ctx, svc_login.ComandoLogin{
		Email:    cmd.Email,
		Password: cmd.Password,
		IPOrigen: cmd.IPOrigen,
	})
	if err != nil {
		return nil, err
	}

	expiresIn := int64(time.Until(respuesta.ExpiracionAccess).Seconds())

	return &RespuestaLogin{
		AccessToken:  respuesta.AccessToken,
		RefreshToken: respuesta.RefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		UsuarioID:    respuesta.UsuarioID,
		SesionID:     respuesta.SesionID,
	}, nil
}

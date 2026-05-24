package facades

import (
	"context"
	"time"

	uc_sesiones_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	uc_sesiones_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	uc_sesiones_refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	svc_registro "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/services/registro"
)

// authFacadeImpl implementa AuthFacade orquestando los servicios de aplicación.
type authFacadeImpl struct {
	servicioRegistro svc_registro.EjecutorRegistro
	loginUseCase     LoginUseCase
	refreshUseCase   RefreshUseCase
	logoutUseCase    LogoutUseCase
}

// NewAuthFacade construye la implementación concreta de AuthFacade.
func NewAuthFacade(
	servicioRegistro svc_registro.EjecutorRegistro,
	loginUseCase LoginUseCase,
	refreshUseCase RefreshUseCase,
	logoutUseCase LogoutUseCase,
) AuthFacade {
	return &authFacadeImpl{
		servicioRegistro: servicioRegistro,
		loginUseCase:     loginUseCase,
		refreshUseCase:   refreshUseCase,
		logoutUseCase:    logoutUseCase,
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
	respuesta, err := f.loginUseCase.Ejecutar(ctx, uc_sesiones_login.ComandoIniciarSesion{
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

// Refresh renueva la sesión usando el refresh token.
func (f *authFacadeImpl) Refresh(ctx context.Context, cmd ComandoRefresh) (*RespuestaRefresh, error) {
	respuesta, err := f.refreshUseCase.Ejecutar(ctx, uc_sesiones_refresh.ComandoRenovarSesion{
		RefreshToken: cmd.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	expiresIn := int64(time.Until(respuesta.ExpiracionAccess).Seconds())

	return &RespuestaRefresh{
		AccessToken:  respuesta.AccessToken,
		RefreshToken: respuesta.RefreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		UsuarioID:    respuesta.UsuarioID,
	}, nil
}

// Logout cierra una sesión específica del usuario autenticado.
func (f *authFacadeImpl) Logout(ctx context.Context, cmd ComandoLogout) (*RespuestaLogout, error) {
	respuesta, err := f.logoutUseCase.Ejecutar(ctx, uc_sesiones_logout.ComandoCerrarSesion{
		SesionID:  cmd.SesionID,
		UsuarioID: cmd.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaLogout{
		SesionesRevocadas: respuesta.SesionesRevocadas,
	}, nil
}

// LogoutAll cierra todas las sesiones del usuario autenticado.
func (f *authFacadeImpl) LogoutAll(ctx context.Context, cmd ComandoLogoutAll) (*RespuestaLogout, error) {
	respuesta, err := f.logoutUseCase.CerrarTodas(ctx, uc_sesiones_logout.ComandoCerrarTodasLasSesiones{
		UsuarioID: cmd.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaLogout{
		SesionesRevocadas: respuesta.SesionesRevocadas,
	}, nil
}

package facades

import (
	"context"
	"fmt"
	"time"

	uc_sesiones_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	uc_sesiones_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	uc_sesiones_refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	uc_register "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/register"
	uc_verifyemail "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/verifyemail"
)

type authFacadeImpl struct {
	registroUseCase     RegistroUseCase
	verificacionUseCase *uc_verifyemail.VerificarCorreoCasoDeUso
	loginUseCase        LoginUseCase
	refreshUseCase      RefreshUseCase
	logoutUseCase       LogoutUseCase
}

func NewAuthFacade(
	registroUseCase RegistroUseCase,
	verificacionUseCase *uc_verifyemail.VerificarCorreoCasoDeUso,
	loginUseCase LoginUseCase,
	refreshUseCase RefreshUseCase,
	logoutUseCase LogoutUseCase,
) AuthFacade {
	return &authFacadeImpl{
		registroUseCase:     registroUseCase,
		verificacionUseCase: verificacionUseCase,
		loginUseCase:        loginUseCase,
		refreshUseCase:      refreshUseCase,
		logoutUseCase:       logoutUseCase,
	}
}

func (f *authFacadeImpl) Registrar(ctx context.Context, cmd ComandoRegistro) (*RespuestaRegistro, error) {
	respuesta, err := f.registroUseCase.Ejecutar(ctx, &uc_register.ComandoRegistrarUsuario{
		Correo:   cmd.Correo,
		Password: cmd.Password,
		Nombre:   cmd.Nombre,
		Apellido: cmd.Apellido,
		Telefono: cmd.Telefono,
	})
	if err != nil {
		return nil, err
	}

	// Post-registro: enviar email de verificación (best-effort)
	go func() {
		if _, err := f.verificacionUseCase.Solicitar(context.Background(), uc_verifyemail.ComandoSolicitarVerificacion{
			UsuarioID: respuesta.UsuarioID,
		}); err != nil {
			fmt.Printf("[AuthFacade] Error al solicitar verificación: %v\n", err)
		}
	}()

	return &RespuestaRegistro{
		UsuarioID: respuesta.UsuarioID,
		Correo:    respuesta.Correo,
		Estado:    respuesta.Estado,
	}, nil
}

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
		TenantID:     respuesta.TenantID,
		Rol:          respuesta.Rol,
	}, nil
}

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
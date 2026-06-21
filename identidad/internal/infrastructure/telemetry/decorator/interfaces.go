package decorator

import (
	"context"
	uc_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	uc_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	uc_refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	uc_register "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/register"
)

type LoginUseCase interface {
	Ejecutar(ctx context.Context, cmd uc_login.ComandoIniciarSesion) (*uc_login.RespuestaIniciarSesion, error)
}

type RefreshUseCase interface {
	Ejecutar(ctx context.Context, cmd uc_refresh.ComandoRenovarSesion) (*uc_refresh.RespuestaRenovarSesion, error)
}

type LogoutUseCase interface {
	Ejecutar(ctx context.Context, cmd uc_logout.ComandoCerrarSesion) (*uc_logout.RespuestaCerrarSesion, error)
	CerrarTodas(ctx context.Context, cmd uc_logout.ComandoCerrarTodasLasSesiones) (*uc_logout.RespuestaCerrarSesion, error)
}

type RegistroUseCase interface {
	Ejecutar(ctx context.Context, cmd *uc_register.ComandoRegistrarUsuario) (*uc_register.RespuestaRegistrarUsuario, error)
}

package decorator

import (
	"context"

	uc_logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
)

// UseCase is a generic interface for any use case with a single Ejecutar method.
// Cmd can be a struct, pointer, or scalar. Resp is typically a pointer to a response struct.
// This replaces per-use-case interfaces (LoginUseCase, RefreshUseCase, etc.)
// by leveraging Go structural typing.
type UseCase[Cmd, Resp any] interface {
	Ejecutar(context.Context, Cmd) (Resp, error)
}

// LogoutUseCase is kept separate because Logout has two methods (Ejecutar + CerrarTodas).
type LogoutUseCase interface {
	Ejecutar(ctx context.Context, cmd uc_logout.ComandoCerrarSesion) (*uc_logout.RespuestaCerrarSesion, error)
	CerrarTodas(ctx context.Context, cmd uc_logout.ComandoCerrarTodasLasSesiones) (*uc_logout.RespuestaCerrarSesion, error)
}

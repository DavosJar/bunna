package facades

import (
	"context"

	decorator "github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/decorator"
	uc_confirmar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/confirmarverificacion"
	uc_reenviar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/reenviarverificacion"
	uc_solicitar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/solicitarverificacion"
)

type ComandoSolicitarVerificacion struct {
	UsuarioID string
}

type RespuestaSolicitarVerificacion struct {
	Mensaje string
}

type ComandoConfirmarVerificacion struct {
	Token string
}

type RespuestaConfirmarVerificacion struct {
	Mensaje string
}

type ComandoReenviarVerificacion struct {
	UsuarioID string
}

type VerificacionFacade interface {
	SolicitarVerificacion(ctx context.Context, cmd ComandoSolicitarVerificacion) (*RespuestaSolicitarVerificacion, error)
	ConfirmarVerificacion(ctx context.Context, cmd ComandoConfirmarVerificacion) (*RespuestaConfirmarVerificacion, error)
	ReenviarVerificacion(ctx context.Context, cmd ComandoReenviarVerificacion) (*RespuestaSolicitarVerificacion, error)
}

type verificacionFacadeImpl struct {
	solicitarUseCase decorator.UseCase[*uc_solicitar.ComandoSolicitarVerificacion, *uc_solicitar.RespuestaSolicitarVerificacion]
	confirmarUseCase decorator.UseCase[*uc_confirmar.ComandoConfirmarVerificacion, *uc_confirmar.RespuestaConfirmarVerificacion]
	reenviarUseCase  decorator.UseCase[*uc_reenviar.ComandoReenviarVerificacion, *uc_reenviar.RespuestaSolicitarVerificacion]
}

func NewVerificacionFacade(
	solicitarUseCase decorator.UseCase[*uc_solicitar.ComandoSolicitarVerificacion, *uc_solicitar.RespuestaSolicitarVerificacion],
	confirmarUseCase decorator.UseCase[*uc_confirmar.ComandoConfirmarVerificacion, *uc_confirmar.RespuestaConfirmarVerificacion],
	reenviarUseCase decorator.UseCase[*uc_reenviar.ComandoReenviarVerificacion, *uc_reenviar.RespuestaSolicitarVerificacion],
) VerificacionFacade {
	return &verificacionFacadeImpl{
		solicitarUseCase: solicitarUseCase,
		confirmarUseCase: confirmarUseCase,
		reenviarUseCase:  reenviarUseCase,
	}
}

func (f *verificacionFacadeImpl) SolicitarVerificacion(ctx context.Context, cmd ComandoSolicitarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	resp, err := f.solicitarUseCase.Ejecutar(ctx, &uc_solicitar.ComandoSolicitarVerificacion{
		UsuarioID: cmd.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaSolicitarVerificacion{Mensaje: resp.Mensaje}, nil
}

func (f *verificacionFacadeImpl) ConfirmarVerificacion(ctx context.Context, cmd ComandoConfirmarVerificacion) (*RespuestaConfirmarVerificacion, error) {
	resp, err := f.confirmarUseCase.Ejecutar(ctx, &uc_confirmar.ComandoConfirmarVerificacion{
		Token: cmd.Token,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaConfirmarVerificacion{Mensaje: resp.Mensaje}, nil
}

func (f *verificacionFacadeImpl) ReenviarVerificacion(ctx context.Context, cmd ComandoReenviarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	resp, err := f.reenviarUseCase.Ejecutar(ctx, &uc_reenviar.ComandoReenviarVerificacion{
		UsuarioID: cmd.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaSolicitarVerificacion{Mensaje: resp.Mensaje}, nil
}

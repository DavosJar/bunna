package facades

import (
	"context"

	decorator "github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/decorator"
	uc_confirmar "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/confirmarrecuperacion"
	uc_solicitar "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/solicitarrecuperacion"
	uc_validar "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/validartokenrecuperacion"
)

type ComandoSolicitarRecuperacion struct {
	Email    string
	IPOrigen string
}

type RespuestaSolicitarRecuperacion struct {
	Mensaje string
}

type ComandoValidarTokenRecuperacion struct {
	Token string
}

type RespuestaValidarTokenRecuperacion struct {
	UsuarioID string
	Valido    bool
}

type ComandoConfirmarRecuperacion struct {
	Token         string
	NuevaPassword string
}

type RespuestaConfirmarRecuperacion struct {
	Mensaje string
}

type RecuperacionFacade interface {
	SolicitarRecuperacion(ctx context.Context, cmd ComandoSolicitarRecuperacion) (*RespuestaSolicitarRecuperacion, error)
	ValidarTokenRecuperacion(ctx context.Context, cmd ComandoValidarTokenRecuperacion) (*RespuestaValidarTokenRecuperacion, error)
	ConfirmarRecuperacion(ctx context.Context, cmd ComandoConfirmarRecuperacion) (*RespuestaConfirmarRecuperacion, error)
}

type recuperacionFacadeImpl struct {
	solicitarUseCase decorator.UseCase[*uc_solicitar.ComandoSolicitarRecuperacion, *uc_solicitar.RespuestaSolicitarRecuperacion]
	validarUseCase   decorator.UseCase[*uc_validar.ComandoValidarTokenRecuperacion, *uc_validar.RespuestaValidarTokenRecuperacion]
	confirmarUseCase decorator.UseCase[*uc_confirmar.ComandoConfirmarRecuperacion, *uc_confirmar.RespuestaConfirmarRecuperacion]
}

func NewRecuperacionFacade(
	solicitarUseCase decorator.UseCase[*uc_solicitar.ComandoSolicitarRecuperacion, *uc_solicitar.RespuestaSolicitarRecuperacion],
	validarUseCase decorator.UseCase[*uc_validar.ComandoValidarTokenRecuperacion, *uc_validar.RespuestaValidarTokenRecuperacion],
	confirmarUseCase decorator.UseCase[*uc_confirmar.ComandoConfirmarRecuperacion, *uc_confirmar.RespuestaConfirmarRecuperacion],
) RecuperacionFacade {
	return &recuperacionFacadeImpl{
		solicitarUseCase: solicitarUseCase,
		validarUseCase:   validarUseCase,
		confirmarUseCase: confirmarUseCase,
	}
}

func (f *recuperacionFacadeImpl) SolicitarRecuperacion(ctx context.Context, cmd ComandoSolicitarRecuperacion) (*RespuestaSolicitarRecuperacion, error) {
	resp, err := f.solicitarUseCase.Ejecutar(ctx, &uc_solicitar.ComandoSolicitarRecuperacion{
		Email:    cmd.Email,
		IPOrigen: cmd.IPOrigen,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaSolicitarRecuperacion{Mensaje: resp.Mensaje}, nil
}

func (f *recuperacionFacadeImpl) ValidarTokenRecuperacion(ctx context.Context, cmd ComandoValidarTokenRecuperacion) (*RespuestaValidarTokenRecuperacion, error) {
	resp, err := f.validarUseCase.Ejecutar(ctx, &uc_validar.ComandoValidarTokenRecuperacion{
		Token: cmd.Token,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaValidarTokenRecuperacion{
		UsuarioID: resp.UsuarioID,
		Valido:    resp.Valido,
	}, nil
}

func (f *recuperacionFacadeImpl) ConfirmarRecuperacion(ctx context.Context, cmd ComandoConfirmarRecuperacion) (*RespuestaConfirmarRecuperacion, error) {
	resp, err := f.confirmarUseCase.Ejecutar(ctx, &uc_confirmar.ComandoConfirmarRecuperacion{
		Token:         cmd.Token,
		NuevaPassword: cmd.NuevaPassword,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaConfirmarRecuperacion{Mensaje: resp.Mensaje}, nil
}

package facades

import (
	"context"

	uc_verifyemail "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/verifyemail"
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
	verificarCorreo *uc_verifyemail.VerificarCorreoCasoDeUso
}

func NewVerificacionFacade(verificarCorreo *uc_verifyemail.VerificarCorreoCasoDeUso) VerificacionFacade {
	return &verificacionFacadeImpl{verificarCorreo: verificarCorreo}
}

func (f *verificacionFacadeImpl) SolicitarVerificacion(ctx context.Context, cmd ComandoSolicitarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	resp, err := f.verificarCorreo.Solicitar(ctx, uc_verifyemail.ComandoSolicitarVerificacion{
		UsuarioID: cmd.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaSolicitarVerificacion{Mensaje: resp.Mensaje}, nil
}

func (f *verificacionFacadeImpl) ConfirmarVerificacion(ctx context.Context, cmd ComandoConfirmarVerificacion) (*RespuestaConfirmarVerificacion, error) {
	resp, err := f.verificarCorreo.Confirmar(ctx, uc_verifyemail.ComandoConfirmarVerificacion{
		Token: cmd.Token,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaConfirmarVerificacion{Mensaje: resp.Mensaje}, nil
}

func (f *verificacionFacadeImpl) ReenviarVerificacion(ctx context.Context, cmd ComandoReenviarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	resp, err := f.verificarCorreo.Reenviar(ctx, uc_verifyemail.ComandoReenviarVerificacion{
		UsuarioID: cmd.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaSolicitarVerificacion{Mensaje: resp.Mensaje}, nil
}

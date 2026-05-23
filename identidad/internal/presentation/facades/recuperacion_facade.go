package facades

import (
	"context"

	uc_forgotpassword "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/forgotpassword"
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
	recuperarContrasena *uc_forgotpassword.RecuperarContrasenaCasoDeUso
}

func NewRecuperacionFacade(recuperarContrasena *uc_forgotpassword.RecuperarContrasenaCasoDeUso) RecuperacionFacade {
	return &recuperacionFacadeImpl{recuperarContrasena: recuperarContrasena}
}

func (f *recuperacionFacadeImpl) SolicitarRecuperacion(ctx context.Context, cmd ComandoSolicitarRecuperacion) (*RespuestaSolicitarRecuperacion, error) {
	resp, err := f.recuperarContrasena.Solicitar(ctx, uc_forgotpassword.ComandoSolicitarRecuperacion{
		Email:    cmd.Email,
		IPOrigen: cmd.IPOrigen,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaSolicitarRecuperacion{Mensaje: resp.Mensaje}, nil
}

func (f *recuperacionFacadeImpl) ValidarTokenRecuperacion(ctx context.Context, cmd ComandoValidarTokenRecuperacion) (*RespuestaValidarTokenRecuperacion, error) {
	resp, err := f.recuperarContrasena.ValidarToken(ctx, uc_forgotpassword.ComandoValidarTokenRecuperacion{
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
	resp, err := f.recuperarContrasena.Confirmar(ctx, uc_forgotpassword.ComandoConfirmarRestablecimiento{
		Token:         cmd.Token,
		NuevaPassword: cmd.NuevaPassword,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaConfirmarRecuperacion{Mensaje: resp.Mensaje}, nil
}

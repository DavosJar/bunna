package facades

import (
	"context"

	uc_aceptar "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/aceptarinvitacion"
	uc_crear "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/crearinvitacion"
)

type ComandoCrearInvitacion struct {
	TenantID  string
	RolID     string
	Correo    string
	CreadoPor string
}

type RespuestaCrearInvitacion struct {
	Mensaje string
}

type ComandoAceptarInvitacion struct {
	Token     string
	UsuarioID string
}

type RespuestaAceptarInvitacion struct {
	TenantID string
	RolID    string
}

type InvitacionFacade interface {
	CrearInvitacion(ctx context.Context, cmd ComandoCrearInvitacion) (*RespuestaCrearInvitacion, error)
	AceptarInvitacion(ctx context.Context, cmd ComandoAceptarInvitacion) (*RespuestaAceptarInvitacion, error)
}

type invitacionFacadeImpl struct {
	crearUseCase   *uc_crear.CrearInvitacionCasoDeUso
	aceptarUseCase *uc_aceptar.AceptarInvitacionCasoDeUso
}

func NewInvitacionFacade(
	crearUseCase *uc_crear.CrearInvitacionCasoDeUso,
	aceptarUseCase *uc_aceptar.AceptarInvitacionCasoDeUso,
) InvitacionFacade {
	return &invitacionFacadeImpl{
		crearUseCase:   crearUseCase,
		aceptarUseCase: aceptarUseCase,
	}
}

func (f *invitacionFacadeImpl) CrearInvitacion(ctx context.Context, cmd ComandoCrearInvitacion) (*RespuestaCrearInvitacion, error) {
	_, err := f.crearUseCase.Ejecutar(ctx, &uc_crear.ComandoCrearInvitacion{
		TenantID:  cmd.TenantID,
		RolID:     cmd.RolID,
		Correo:    cmd.Correo,
		CreadoPor: cmd.CreadoPor,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaCrearInvitacion{
		Mensaje: "Invitación enviada exitosamente",
	}, nil
}

func (f *invitacionFacadeImpl) AceptarInvitacion(ctx context.Context, cmd ComandoAceptarInvitacion) (*RespuestaAceptarInvitacion, error) {
	resp, err := f.aceptarUseCase.Ejecutar(ctx, &uc_aceptar.ComandoAceptarInvitacion{
		Token:     cmd.Token,
		UsuarioID: cmd.UsuarioID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaAceptarInvitacion{
		TenantID: resp.TenantID,
		RolID:    resp.RolID,
	}, nil
}

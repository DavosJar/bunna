package eliminarinvitacion

import (
	"context"
	"fmt"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type EliminarInvitacionCasoDeUso struct {
	invitacionRepo invitaciones.InvitacionRepositorio
	authSvc        rbac.AuthorizationService
}

func NewEliminarInvitacionCasoDeUso(
	invitacionRepo invitaciones.InvitacionRepositorio,
	authSvc rbac.AuthorizationService,
) *EliminarInvitacionCasoDeUso {
	return &EliminarInvitacionCasoDeUso{
		invitacionRepo: invitacionRepo,
		authSvc:        authSvc,
	}
}

func (uc *EliminarInvitacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoEliminarInvitacion) (*RespuestaEliminarInvitacion, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoUsuarioInvitar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	invitacion, err := uc.invitacionRepo.ObtenerPorID(ctx, cmd.InvitacionID)
	if err != nil {
		return nil, err
	}

	if invitacion.TenantID() != cmd.TenantID {
		return nil, invitaciones.ErrNoEncontrada
	}

	if invitacion.EstaAceptada() {
		return nil, invitaciones.ErrYaAceptada
	}

	if err := uc.invitacionRepo.Eliminar(ctx, cmd.InvitacionID); err != nil {
		return nil, fmt.Errorf("error al eliminar invitación: %w", err)
	}

	return &RespuestaEliminarInvitacion{
		Mensaje: "Invitación eliminada exitosamente",
	}, nil
}

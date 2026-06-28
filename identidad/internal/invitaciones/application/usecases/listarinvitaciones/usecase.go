package listarinvitaciones

import (
	"context"
	"time"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type ListarInvitacionesCasoDeUso struct {
	invitacionRepo invitaciones.InvitacionRepositorio
	rolRepo        rbac.RolRepositorio
}

func NewListarInvitacionesCasoDeUso(
	invitacionRepo invitaciones.InvitacionRepositorio,
	rolRepo rbac.RolRepositorio,
) *ListarInvitacionesCasoDeUso {
	return &ListarInvitacionesCasoDeUso{
		invitacionRepo: invitacionRepo,
		rolRepo:        rolRepo,
	}
}

func (uc *ListarInvitacionesCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoListarInvitaciones) (*RespuestaListarInvitaciones, error) {
	paginacion := shareddomain.Paginacion{
		Pagina:       cmd.Pagina,
		TamanoPagina: cmd.TamanoPagina,
	}

	items, total, err := uc.invitacionRepo.ListarPorTenant(ctx, cmd.TenantID, paginacion, cmd.Estado)
	if err != nil {
		return nil, err
	}

	dtos := make([]InvitacionDTO, len(items))
	for i, inv := range items {
		rolNombre := ""
		if rol, err := uc.rolRepo.ObtenerPorID(ctx, inv.RolID()); err == nil {
			rolNombre = rol.Nombre
		}

		estado := "pendiente"
		if inv.EstaAceptada() {
			estado = "aceptada"
		} else if inv.Expiro(time.Now()) {
			estado = "expirada"
		}

		dtos[i] = InvitacionDTO{
			ID:            inv.ID(),
			Email:         inv.Email(),
			Nombre:        inv.Nombre(),
			RolID:         inv.RolID(),
			RolNombre:     rolNombre,
			Estado:        estado,
			FechaCreacion: inv.FechaCreacion().Format(time.RFC3339),
			Expiracion:    inv.Expiracion().Format(time.RFC3339),
		}
	}

	return &RespuestaListarInvitaciones{
		Invitaciones: dtos,
		Total:        total,
	}, nil
}

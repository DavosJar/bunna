package aceptarInvitacion

import (
	"context"
	"fmt"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type AceptarInvitacionCasoDeUso struct {
	invitacionRepo     invitaciones.InvitacionRepositorio
	membresiaRepo      tenant.MembresiaRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
}

func NewAceptarInvitacionCasoDeUso(
	invitacionRepo invitaciones.InvitacionRepositorio,
	membresiaRepo tenant.MembresiaRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
) *AceptarInvitacionCasoDeUso {
	return &AceptarInvitacionCasoDeUso{
		invitacionRepo:       invitacionRepo,
		membresiaRepo:        membresiaRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
	}
}

func (uc *AceptarInvitacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoAceptarInvitacion) (*RespuestaAceptarInvitacion, error) {
	if cmd.Token == "" {
		return nil, invitaciones.ErrEnlaceInvalido
	}

	tokenHash := invitaciones.HashearTokenPublico(cmd.Token)

	invitacion, err := uc.invitacionRepo.ObtenerPorTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, invitaciones.ErrEnlaceInvalido
	}

	if err := invitacion.Aceptar(); err != nil {
		return nil, err
	}

	if err := uc.invitacionRepo.MarcarAceptada(ctx, invitacion.ID()); err != nil {
		return nil, fmt.Errorf("error al marcar invitación como aceptada: %w", err)
	}

	miembro, err := tenant.NuevaMembresia(cmd.UsuarioID, invitacion.TenantID())
	if err != nil {
		return nil, err
	}

	if err := uc.membresiaRepo.Crear(ctx, miembro); err != nil {
		return nil, fmt.Errorf("error al crear membresía: %w", err)
	}

	if err := uc.usuarioTenantRolRepo.Crear(ctx, cmd.UsuarioID, invitacion.TenantID(), invitacion.RolID()); err != nil {
		return nil, fmt.Errorf("error al asignar rol: %w", err)
	}

	return &RespuestaAceptarInvitacion{
		TenantID: invitacion.TenantID(),
		RolID:    invitacion.RolID(),
	}, nil
}

package aceptarInvitacion

import (
	"context"
	"fmt"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type AceptarInvitacionCasoDeUso struct {
	invitacionRepo       invitaciones.InvitacionRepositorio
	membresiaRepo        tenant.MembresiaRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
	usuarioRepo          usuario.UsuarioRepositorio
}

func NewAceptarInvitacionCasoDeUso(
	invitacionRepo invitaciones.InvitacionRepositorio,
	membresiaRepo tenant.MembresiaRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
	usuarioRepo usuario.UsuarioRepositorio,
) *AceptarInvitacionCasoDeUso {
	return &AceptarInvitacionCasoDeUso{
		invitacionRepo:       invitacionRepo,
		membresiaRepo:        membresiaRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
		usuarioRepo:          usuarioRepo,
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

	// Buscar usuario por el email de la invitación
	user, err := uc.usuarioRepo.ObtenerPorCorreo(ctx, invitacion.Email())
	if err != nil {
		return nil, invitaciones.ErrUsuarioNoRegistrado
	}

	if err := uc.invitacionRepo.MarcarAceptada(ctx, invitacion.ID()); err != nil {
		return nil, fmt.Errorf("error al marcar invitación como aceptada: %w", err)
	}

	miembro, err := tenant.NuevaMembresia(user.ID(), invitacion.TenantID())
	if err != nil {
		return nil, err
	}

	if err := uc.membresiaRepo.Crear(ctx, miembro); err != nil {
		return nil, fmt.Errorf("error al crear membresía: %w", err)
	}

	if err := uc.usuarioTenantRolRepo.Crear(ctx, user.ID(), invitacion.TenantID(), invitacion.RolID()); err != nil {
		return nil, fmt.Errorf("error al asignar rol: %w", err)
	}

	return &RespuestaAceptarInvitacion{
		TenantID: invitacion.TenantID(),
		RolID:    invitacion.RolID(),
	}, nil
}

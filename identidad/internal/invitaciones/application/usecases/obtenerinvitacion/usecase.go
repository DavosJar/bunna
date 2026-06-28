package obtenerinvitacion

import (
	"context"
	"time"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type ObtenerInvitacionCasoDeUso struct {
	invitacionRepo invitaciones.InvitacionRepositorio
	tenantRepo     tenant.TenantRepositorio
	rolRepo        rbac.RolRepositorio
}

func NewObtenerInvitacionCasoDeUso(
	invitacionRepo invitaciones.InvitacionRepositorio,
	tenantRepo tenant.TenantRepositorio,
	rolRepo rbac.RolRepositorio,
) *ObtenerInvitacionCasoDeUso {
	return &ObtenerInvitacionCasoDeUso{
		invitacionRepo: invitacionRepo,
		tenantRepo:     tenantRepo,
		rolRepo:        rolRepo,
	}
}

func (uc *ObtenerInvitacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoObtenerInvitacion) (*RespuestaObtenerInvitacion, error) {
	if cmd.Token == "" {
		return nil, invitaciones.ErrEnlaceInvalido
	}

	tokenHash := invitaciones.HashearTokenPublico(cmd.Token)
	inv, err := uc.invitacionRepo.ObtenerPorTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	tenantObj, err := uc.tenantRepo.ObtenerPorID(ctx, inv.TenantID())
	var tenantNombre string
	if err == nil {
		tenantNombre = tenantObj.Nombre()
	}

	rolObj, err := uc.rolRepo.ObtenerPorID(ctx, inv.RolID())
	var rolNombre string
	if err == nil {
		rolNombre = rolObj.Nombre
	}

	estado := "pendiente"
	if inv.EstaAceptada() {
		estado = "aceptada"
	} else if inv.Expiro(time.Now()) {
		estado = "expirada"
	}

	return &RespuestaObtenerInvitacion{
		ID:           inv.ID(),
		TenantID:     inv.TenantID(),
		TenantNombre: tenantNombre,
		RolID:        inv.RolID(),
		RolNombre:    rolNombre,
		Email:        inv.Email(),
		Estado:       estado,
		Expiracion:   inv.Expiracion().Format(time.RFC3339),
	}, nil
}

package expeluser

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type ExpulsarUsuarioCasoDeUso struct {
	membresiaRepo      tenant.MembresiaRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
	authSvc              rbac.AuthorizationService
}

func NewExpulsarUsuarioCasoDeUso(
	membresiaRepo tenant.MembresiaRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
	authSvc rbac.AuthorizationService,
) *ExpulsarUsuarioCasoDeUso {
	return &ExpulsarUsuarioCasoDeUso{
		membresiaRepo:        membresiaRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
		authSvc:              authSvc,
	}
}

func (uc *ExpulsarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoExpulsarUsuario) (*RespuestaExpulsarUsuario, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoUsuarioExpulsar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	if cmd.EjecutorID == cmd.UsuarioID {
		return nil, fmt.Errorf("no puedes expulsarte a ti mismo")
	}

	// 1. Quitar todos los roles del usuario en el tenant
	roles, err := uc.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, cmd.UsuarioID, cmd.TenantID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener roles del usuario en el tenant: %w", err)
	}
	for _, rol := range roles {
		if err := uc.usuarioTenantRolRepo.Eliminar(ctx, cmd.UsuarioID, cmd.TenantID, rol.ID); err != nil {
			return nil, fmt.Errorf("error al revocar rol %s: %w", rol.ID, err)
		}
	}

	// 2. Quitar membresía del tenant
	if err := uc.membresiaRepo.Eliminar(ctx, cmd.UsuarioID, cmd.TenantID); err != nil {
		return nil, fmt.Errorf("error al eliminar membresía: %w", err)
	}

	return &RespuestaExpulsarUsuario{
		UsuarioID:         cmd.UsuarioID,
		Estado:            "EXPULSADO",
		SesionesRevocadas: 0,
		ExpulsadoEn:       time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

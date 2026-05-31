package updatetenant

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type ConfigurarTenantCasoDeUso struct {
	tenantRepo tenant.TenantRepositorio
	authSvc    rbac.AuthorizationService
}

func NewConfigurarTenantCasoDeUso(
	tenantRepo tenant.TenantRepositorio,
	authSvc rbac.AuthorizationService,
) *ConfigurarTenantCasoDeUso {
	return &ConfigurarTenantCasoDeUso{tenantRepo: tenantRepo, authSvc: authSvc}
}

func (uc *ConfigurarTenantCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoConfigurarTenant) (*RespuestaConfigurarTenant, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoTenantConfigurar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	t, err := uc.tenantRepo.ObtenerPorID(ctx, cmd.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant no encontrado: %w", err)
	}

	_ = t

	return &RespuestaConfigurarTenant{
		TenantID:     cmd.TenantID,
		Nombre:       cmd.Nombre,
		Slug:         cmd.Slug,
		ModificadoEn: time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

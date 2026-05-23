package facades

import (
	"context"

	uc_updatetenant "github.com/davosjar/bunna/services/identidad/internal/tenants/application/usecases/updatetenant"
)

type ComandoConfigurarTenant struct {
	TenantID   string
	Nombre     string
	Slug       string
	EjecutorID string
}

type RespuestaConfigurarTenant struct {
	TenantID     string
	Nombre       string
	Slug         string
	ModificadoEn string
}

type TenantFacade interface {
	ConfigurarTenant(ctx context.Context, cmd ComandoConfigurarTenant) (*RespuestaConfigurarTenant, error)
}

type tenantFacadeImpl struct {
	configurarTenant *uc_updatetenant.ConfigurarTenantCasoDeUso
}

func NewTenantFacade(configurarTenant *uc_updatetenant.ConfigurarTenantCasoDeUso) TenantFacade {
	return &tenantFacadeImpl{configurarTenant: configurarTenant}
}

func (f *tenantFacadeImpl) ConfigurarTenant(ctx context.Context, cmd ComandoConfigurarTenant) (*RespuestaConfigurarTenant, error) {
	resp, err := f.configurarTenant.Ejecutar(ctx, &uc_updatetenant.ComandoConfigurarTenant{
		TenantID:   cmd.TenantID,
		Nombre:     cmd.Nombre,
		Slug:       cmd.Slug,
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaConfigurarTenant{
		TenantID:     resp.TenantID,
		Nombre:       resp.Nombre,
		Slug:         resp.Slug,
		ModificadoEn: resp.ModificadoEn,
	}, nil
}

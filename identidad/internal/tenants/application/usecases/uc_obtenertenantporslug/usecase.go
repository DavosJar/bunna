package uc_obtenertenantporslug

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type ObtenerTenantPorSlugCasoDeUso struct {
	tenantRepo tenant.TenantRepositorio
}

func NewObtenerTenantPorSlugCasoDeUso(tenantRepo tenant.TenantRepositorio) *ObtenerTenantPorSlugCasoDeUso {
	return &ObtenerTenantPorSlugCasoDeUso{
		tenantRepo: tenantRepo,
	}
}

func (uc *ObtenerTenantPorSlugCasoDeUso) Ejecutar(ctx context.Context, slug string) (*RespuestaObtenerTenantPorSlug, error) {
	t, err := uc.tenantRepo.ObtenerPorSlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return &RespuestaObtenerTenantPorSlug{
		ID:            t.ID(),
		Nombre:        t.Nombre(),
		Slug:          t.Slug(),
		Activo:        t.EstaActivo(),
		FechaCreacion: t.FechaCreacion(),
	}, nil
}

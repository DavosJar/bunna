package uc_obtenertenantporid

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type ObtenerTenantPorIDCasoDeUso struct {
	tenantRepo tenant.TenantRepositorio
}

func NewObtenerTenantPorIDCasoDeUso(tenantRepo tenant.TenantRepositorio) *ObtenerTenantPorIDCasoDeUso {
	return &ObtenerTenantPorIDCasoDeUso{
		tenantRepo: tenantRepo,
	}
}

func (uc *ObtenerTenantPorIDCasoDeUso) Ejecutar(ctx context.Context, id string) (*RespuestaObtenerTenantPorID, error) {
	t, err := uc.tenantRepo.ObtenerPorID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &RespuestaObtenerTenantPorID{
		ID:            t.ID(),
		Nombre:        t.Nombre(),
		Slug:          t.Slug(),
		Activo:        t.EstaActivo(),
		FechaCreacion: t.FechaCreacion(),
	}, nil
}

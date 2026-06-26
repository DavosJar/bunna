package uc_listarmistenants

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type ListarMisTenantsCasoDeUso struct {
	tenantRepo tenant.TenantRepositorio
}

func NewListarMisTenantsCasoDeUso(tenantRepo tenant.TenantRepositorio) *ListarMisTenantsCasoDeUso {
	return &ListarMisTenantsCasoDeUso{
		tenantRepo: tenantRepo,
	}
}

func (uc *ListarMisTenantsCasoDeUso) Ejecutar(ctx context.Context, usuarioID string) (*RespuestaListarMisTenants, error) {
	tenants, err := uc.tenantRepo.ListarPorUsuario(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	dtos := make([]DtoTenant, len(tenants))
	for i, t := range tenants {
		dtos[i] = DtoTenant{
			ID:            t.ID(),
			Nombre:        t.Nombre(),
			Slug:          t.Slug(),
			Activo:        t.EstaActivo(),
			FechaCreacion: t.FechaCreacion(),
		}
	}

	return &RespuestaListarMisTenants{Tenants: dtos}, nil
}

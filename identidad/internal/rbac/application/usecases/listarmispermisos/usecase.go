package listarmispermisos

import (
	"context"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type ListarMisPermisosCasoDeUso struct {
	rolRepo        rbac.RolRepositorio
	rolPermisoRepo rbac.RolPermisoRepositorio
}

func NewListarMisPermisosCasoDeUso(
	rolRepo rbac.RolRepositorio,
	rolPermisoRepo rbac.RolPermisoRepositorio,
) *ListarMisPermisosCasoDeUso {
	return &ListarMisPermisosCasoDeUso{
		rolRepo:        rolRepo,
		rolPermisoRepo: rolPermisoRepo,
	}
}

func (uc *ListarMisPermisosCasoDeUso) Ejecutar(ctx context.Context, rolNombre, tenantID string) ([]string, error) {
	rolDB, err := uc.rolRepo.ObtenerPorNombre(ctx, rolNombre)
	if err != nil {
		return nil, err
	}

	// Primero busca permisos del tenant específico
	permisos, err := uc.rolPermisoRepo.ListarPorRolYTenant(ctx, rolDB.ID, tenantID)
	if err != nil {
		return nil, err
	}

	// Si no encuentra, cae al tenant sistema
	if len(permisos) == 0 {
		permisos, err = uc.rolPermisoRepo.ListarPorRolYTenant(ctx, rolDB.ID, rbac.TenantIDSistema)
		if err != nil {
			return nil, err
		}
	}

	codigos := make([]string, len(permisos))
	for i, p := range permisos {
		codigos[i] = p.Codigo
	}
	return codigos, nil
}

package listarmispermisos

import (
	"context"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

// PermisoDTO representa la información de un permiso que se expone al exterior
type PermisoDTO struct {
	Codigo      string
	Nombre      string
	Descripcion string
	Modulo      string
}

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

func (uc *ListarMisPermisosCasoDeUso) Ejecutar(ctx context.Context, rolNombre, tenantID string) ([]PermisoDTO, error) {
	rolDB, err := uc.rolRepo.ObtenerPorNombreYTenant(ctx, rolNombre, tenantID)
	if err != nil {
		return nil, err
	}

	permisos, err := uc.rolPermisoRepo.ListarPorRolYTenant(ctx, rolDB.ID, tenantID)
	if err != nil {
		return nil, err
	}

	items := make([]PermisoDTO, len(permisos))
	for i, p := range permisos {
		items[i] = PermisoDTO{
			Codigo:      p.Codigo,
			Nombre:      p.Nombre,
			Descripcion: p.Descripcion,
			Modulo:      p.Modulo,
		}
	}
	return items, nil
}

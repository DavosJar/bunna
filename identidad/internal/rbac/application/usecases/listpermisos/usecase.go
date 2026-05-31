package listpermisos

import (
	"context"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type ListarPermisosOutput struct {
	Permisos []PermisoItem
	Total    int
}

type PermisoItem struct {
	ID          string
	Codigo      string
	Nombre      string
	Descripcion string
	Modulo      string
}

type ListarPermisosCasoDeUso struct {
	permisoRepo rbac.PermisoRepositorio
	authSvc     rbac.AuthorizationService
}

func NewListarPermisosCasoDeUso(
	permisoRepo rbac.PermisoRepositorio,
	authSvc rbac.AuthorizationService,
) *ListarPermisosCasoDeUso {
	return &ListarPermisosCasoDeUso{
		permisoRepo: permisoRepo,
		authSvc:     authSvc,
	}
}

func (uc *ListarPermisosCasoDeUso) Ejecutar(ctx context.Context, ejecutorID string) (*ListarPermisosOutput, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, ejecutorID, "", rbac.PermisoPermisoConsultar)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	permisos, err := uc.permisoRepo.Listar(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]PermisoItem, len(permisos))
	for i, p := range permisos {
		items[i] = PermisoItem{
			ID:          p.ID,
			Codigo:      p.Codigo,
			Nombre:      p.Nombre,
			Descripcion: p.Descripcion,
			Modulo:      p.Modulo,
		}
	}

	return &ListarPermisosOutput{
		Permisos: items,
		Total:    len(items),
	}, nil
}

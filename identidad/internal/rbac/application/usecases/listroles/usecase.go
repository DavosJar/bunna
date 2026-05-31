package listroles

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
)

type ListarRolesCasoDeUso struct {
	rolRepo     rbac.RolRepositorio
	permisoRepo rbac.PermisoRepositorio
	authSvc     rbac.AuthorizationService
}

func NewListarRolesCasoDeUso(
	rolRepo rbac.RolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
	authSvc rbac.AuthorizationService,
) *ListarRolesCasoDeUso {
	return &ListarRolesCasoDeUso{rolRepo: rolRepo, permisoRepo: permisoRepo, authSvc: authSvc}
}

func (uc *ListarRolesCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoListarRoles) (*RespuestaListarRoles, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoPermisoConsultar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	dbRoles, err := uc.rolRepo.Listar(ctx, rbac.EspecificacionRol{}, cmd.Paginacion)
	if err != nil {
		return nil, fmt.Errorf("error al listar roles: %w", err)
	}

	dtoList := make([]RolDTO, 0, len(dbRoles))
	for _, r := range dbRoles {
		permisos, err := uc.permisoRepo.ListarPorRol(ctx, r.ID)
		if err != nil {
			return nil, fmt.Errorf("error al listar permisos del rol %s: %w", r.ID, err)
		}
		codigos := make([]string, len(permisos))
		for i, p := range permisos {
			codigos[i] = p.Codigo
		}
		dtoList = append(dtoList, RolDTO{
			ID:          r.ID,
			Nombre:      r.Nombre,
			Descripcion: r.Descripcion,
			EsSistema:   r.EsSistema,
			Permisos:    codigos,
		})
	}

	return &RespuestaListarRoles{
		Roles:  dtoList,
		Total:  len(dtoList),
		Pagina: cmd.Paginacion.Pagina,
	}, nil
}

package listusers

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type ListarUsuariosCasoDeUso struct {
	userRepo usuario.UsuarioRepositorio
	authSvc  rbac.AuthorizationService
}

func NewListarUsuariosCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	authSvc rbac.AuthorizationService,
) *ListarUsuariosCasoDeUso {
	return &ListarUsuariosCasoDeUso{userRepo: userRepo, authSvc: authSvc}
}

func (uc *ListarUsuariosCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoListarUsuarios) (*RespuestaListarUsuarios, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoUsuarioConsultar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	spec := usuario.EspecificacionUsuario{ListaLiltros: cmd.Filtros}
	usuarios, err := uc.userRepo.Listar(ctx, spec, cmd.Paginacion)
	if err != nil {
		return nil, fmt.Errorf("error al listar usuarios: %w", err)
	}

	dtoList := make([]UsuarioDTO, 0, len(usuarios))
	for _, u := range usuarios {
		dtoList = append(dtoList, UsuarioDTO{
			ID:       u.ID(),
			Correo:   u.Correo(),
			Nombre:   u.Nombre(),
			Apellido: u.Apellido(),
			Estado:   string(u.Estado()),
			CreadoEn: u.FechaCreacion().Format("2006-01-02T15:04:05Z"),
		})
	}

	return &RespuestaListarUsuarios{
		Usuarios: dtoList,
		Total:    len(dtoList),
		Pagina:   cmd.Paginacion.Pagina,
		Tamano:   cmd.Paginacion.TamanoPagina,
	}, nil
}

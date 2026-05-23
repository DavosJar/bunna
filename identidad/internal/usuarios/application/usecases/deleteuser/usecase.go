package deleteuser

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type DarDeBajaUsuarioCasoDeUso struct {
	userRepo usuario.UsuarioRepositorio
	authSvc  rbac.AuthorizationService
}

func NewDarDeBajaUsuarioCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	authSvc rbac.AuthorizationService,
) *DarDeBajaUsuarioCasoDeUso {
	return &DarDeBajaUsuarioCasoDeUso{userRepo: userRepo, authSvc: authSvc}
}

func (uc *DarDeBajaUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoDarDeBajaUsuario) (*RespuestaDarDeBajaUsuario, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoUsuarioEliminar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	u, err := uc.userRepo.ObtenerPorID(ctx, cmd.UsuarioID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	if err := u.CambiarEstado(usuario.PENDIENTE_DE_ELIMINACION); err != nil {
		return nil, fmt.Errorf("error al cambiar estado: %w", err)
	}

	uActualizado, err := uc.userRepo.Actualizar(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("error al persistir cambio: %w", err)
	}

	return &RespuestaDarDeBajaUsuario{
		UsuarioID: uActualizado.ID(),
		Estado:    string(uActualizado.Estado()),
		BajaEn:    uActualizado.FechaActualizacion().Format("2006-01-02T15:04:05Z"),
	}, nil
}

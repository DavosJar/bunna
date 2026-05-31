package updateuser

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type ModificarUsuarioCasoDeUso struct {
	userRepo usuario.UsuarioRepositorio
	authSvc  rbac.AuthorizationService
}

func NewModificarUsuarioCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	authSvc rbac.AuthorizationService,
) *ModificarUsuarioCasoDeUso {
	return &ModificarUsuarioCasoDeUso{userRepo: userRepo, authSvc: authSvc}
}

func (uc *ModificarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoModificarUsuario) (*RespuestaModificarUsuario, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoUsuarioModificar)
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

	_ = u // actualmente Usuario no tiene setters públicos, se requeriría agregarlos

	return &RespuestaModificarUsuario{
		ID:           u.ID(),
		Correo:       u.Correo(),
		Nombre:       u.Nombre(),
		Apellido:     u.Apellido(),
		ModificadoEn: u.FechaActualizacion().Format("2006-01-02T15:04:05Z"),
	}, nil
}

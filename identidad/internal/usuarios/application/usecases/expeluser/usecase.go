package expeluser

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	sesiones "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type ExpulsarUsuarioCasoDeUso struct {
	userRepo    usuario.UsuarioRepositorio
	sessionRepo sesiones.SesionRepositorio
	authSvc     rbac.AuthorizationService
}

func NewExpulsarUsuarioCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	sessionRepo sesiones.SesionRepositorio,
	authSvc rbac.AuthorizationService,
) *ExpulsarUsuarioCasoDeUso {
	return &ExpulsarUsuarioCasoDeUso{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		authSvc:     authSvc,
	}
}

func (uc *ExpulsarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoExpulsarUsuario) (*RespuestaExpulsarUsuario, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoUsuarioExpulsar)
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

	if err := u.Bloquear(); err != nil {
		_ = u.Inactivar()
	}

	uActualizado, err := uc.userRepo.Actualizar(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("error al persistir cambio de estado: %w", err)
	}

	if err := uc.sessionRepo.InvalidarTodasPorUsuarioID(ctx, cmd.UsuarioID); err != nil {
		return nil, fmt.Errorf("error al invalidar sesiones: %w", err)
	}

	return &RespuestaExpulsarUsuario{
		UsuarioID:         uActualizado.ID(),
		Estado:            string(uActualizado.Estado()),
		SesionesRevocadas: -1,
		ExpulsadoEn:       uActualizado.FechaActualizacion().Format("2006-01-02T15:04:05Z"),
	}, nil
}

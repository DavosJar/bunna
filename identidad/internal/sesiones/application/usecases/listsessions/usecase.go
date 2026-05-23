package listsessions

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	sesiones "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

type ListarSesionesCasoDeUso struct {
	sessionRepo sesiones.SesionRepositorio
	authSvc     rbac.AuthorizationService
}

func NewListarSesionesCasoDeUso(
	sessionRepo sesiones.SesionRepositorio,
	authSvc rbac.AuthorizationService,
) *ListarSesionesCasoDeUso {
	return &ListarSesionesCasoDeUso{sessionRepo: sessionRepo, authSvc: authSvc}
}

func (uc *ListarSesionesCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoListarSesiones) (*RespuestaListarSesiones, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoSesionConsultar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	sesionesActivas, err := uc.sessionRepo.ListarActivasPorUsuarioID(ctx, cmd.UsuarioID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error al listar sesiones: %w", err)
	}

	dtoList := make([]SesionDTO, 0, len(sesionesActivas))
	for _, s := range sesionesActivas {
		dtoList = append(dtoList, SesionDTO{
			ID:              s.ID(),
			UsuarioID:       s.UsuarioID(),
			IPOrigen:        s.IPOrigen(),
			Estado:          string(s.Estado()),
			UltimaActividad: s.FechaActualizacion(),
		})
	}

	return &RespuestaListarSesiones{
		Sesiones: dtoList,
		Total:    len(dtoList),
		Pagina:   cmd.Paginacion.Pagina,
	}, nil
}

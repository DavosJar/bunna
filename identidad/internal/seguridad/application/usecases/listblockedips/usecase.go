package listblockedips

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

type ListarIPsBloqueadasCasoDeUso struct {
	intentoRepo seguridad.IntentoIPRepositorio
	authSvc     rbac.AuthorizationService
}

func NewListarIPsBloqueadasCasoDeUso(
	intentoRepo seguridad.IntentoIPRepositorio,
	authSvc rbac.AuthorizationService,
) *ListarIPsBloqueadasCasoDeUso {
	return &ListarIPsBloqueadasCasoDeUso{intentoRepo: intentoRepo, authSvc: authSvc}
}

func (uc *ListarIPsBloqueadasCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoListarIPsBloqueadas) (*RespuestaListarIPsBloqueadas, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoIPBloqueadaConsultar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	_ = cmd.Paginacion

	return &RespuestaListarIPsBloqueadas{
		IPs:    []IPBloqueadaDTO{},
		Total:  0,
		Pagina: cmd.Paginacion.Pagina,
	}, nil
}

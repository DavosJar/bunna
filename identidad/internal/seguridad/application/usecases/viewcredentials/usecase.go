package viewcredentials

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

type ConsultarCredencialesCasoDeUso struct {
	credRepo seguridad.CredencialesRepositorio
	authSvc  rbac.AuthorizationService
}

func NewConsultarCredencialesCasoDeUso(
	credRepo seguridad.CredencialesRepositorio,
	authSvc rbac.AuthorizationService,
) *ConsultarCredencialesCasoDeUso {
	return &ConsultarCredencialesCasoDeUso{credRepo: credRepo, authSvc: authSvc}
}

func (uc *ConsultarCredencialesCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoConsultarCredenciales) (*RespuestaConsultarCredenciales, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoCredencialesConsultar)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	creds, err := uc.credRepo.ObtenerPorUsuarioID(ctx, cmd.UsuarioID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener credenciales: %w", err)
	}

	bloqueadoHasta := ""
	if !creds.BloqueadoHasta().IsZero() {
		bloqueadoHasta = creds.BloqueadoHasta().Format("2006-01-02T15:04:05Z")
	}

	return &RespuestaConsultarCredenciales{
		UsuarioID:        creds.UsuarioID(),
		Activo:           creds.Activo(),
		CorreoVerificado: creds.CorreoVerificado(),
		IntentosFallidos: creds.IntentosFallidos(),
		BloqueadoHasta:   bloqueadoHasta,
	}, nil
}

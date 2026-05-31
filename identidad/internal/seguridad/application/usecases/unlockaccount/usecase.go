package unlockaccount

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

type DesbloquearCuentaCasoDeUso struct {
	credRepo seguridad.CredencialesRepositorio
	authSvc  rbac.AuthorizationService
}

func NewDesbloquearCuentaCasoDeUso(
	credRepo seguridad.CredencialesRepositorio,
	authSvc rbac.AuthorizationService,
) *DesbloquearCuentaCasoDeUso {
	return &DesbloquearCuentaCasoDeUso{credRepo: credRepo, authSvc: authSvc}
}

func (uc *DesbloquearCuentaCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoDesbloquearCuenta) (*RespuestaDesbloquearCuenta, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoCredencialesDesbloquear)
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

	_ = creds

	ahora := time.Now().Format("2006-01-02T15:04:05Z")
	return &RespuestaDesbloquearCuenta{
		UsuarioID:      cmd.UsuarioID,
		DesbloqueadoEn: ahora,
	}, nil
}

package resetpassword

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/shared/application"
	sesiones "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

type ResetearContrasenaCasoDeUso struct {
	credRepo   seguridad.CredencialesRepositorio
	sesionRepo sesiones.SesionRepositorio
	encSvc     seguridad.EncriptacionServicio
	authSvc    rbac.AuthorizationService
}

func NewResetearContrasenaCasoDeUso(
	credRepo seguridad.CredencialesRepositorio,
	sesionRepo sesiones.SesionRepositorio,
	encSvc seguridad.EncriptacionServicio,
	authSvc rbac.AuthorizationService,
) *ResetearContrasenaCasoDeUso {
	return &ResetearContrasenaCasoDeUso{
		credRepo:   credRepo,
		sesionRepo: sesionRepo,
		encSvc:     encSvc,
		authSvc:    authSvc,
	}
}

func (uc *ResetearContrasenaCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoResetearContrasena) (*RespuestaResetearContrasena, error) {
	if err := application.ValidarFormatoPassword(cmd.NuevaPassword, "nueva_password"); err != nil {
		return nil, err
	}

	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoUsuarioResetearPassword)
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

	nuevoHash, err := uc.encSvc.Hashear(cmd.NuevaPassword)
	if err != nil {
		return nil, fmt.Errorf("error al hashear password: %w", err)
	}

	creds.CambiarHash(nuevoHash)

	if _, err := uc.credRepo.Actualizar(ctx, creds); err != nil {
		return nil, fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	if err := uc.sesionRepo.InvalidarTodasPorUsuarioID(ctx, cmd.UsuarioID); err != nil {
		return nil, fmt.Errorf("error al invalidar sesiones: %w", err)
	}

	ahora := time.Now().Format("2006-01-02T15:04:05Z")
	return &RespuestaResetearContrasena{
		UsuarioID:    cmd.UsuarioID,
		ModificadoEn: ahora,
	}, nil
}

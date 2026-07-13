package createrole

import (
	"context"
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type CrearRolCasoDeUso struct {
	rolRepo        rbac.RolRepositorio
	permisoRepo    rbac.PermisoRepositorio
	rolPermisoRepo rbac.RolPermisoRepositorio
	idGen          shareddomain.GeneradorID
	authSvc        rbac.AuthorizationService
}

func NewCrearRolCasoDeUso(
	rolRepo rbac.RolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
	rolPermisoRepo rbac.RolPermisoRepositorio,
	idGen shareddomain.GeneradorID,
	authSvc rbac.AuthorizationService,
) *CrearRolCasoDeUso {
	return &CrearRolCasoDeUso{
		rolRepo:        rolRepo,
		permisoRepo:    permisoRepo,
		rolPermisoRepo: rolPermisoRepo,
		idGen:          idGen,
		authSvc:        authSvc,
	}
}

func (uc *CrearRolCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoCrearRol) (*RespuestaCrearRol, error) {
	ok, err := uc.authSvc.TienePermiso(ctx, cmd.EjecutorID, cmd.TenantID, rbac.PermisoRolCrear)
	if err != nil {
		return nil, fmt.Errorf("error al verificar permiso: %w", err)
	}
	if !ok {
		return nil, rbac.ErrPermisoDenegado
	}

	id, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID de rol: %w", err)
	}

	rol := &rbac.RolDB{
		ID:          id,
		Nombre:      cmd.Nombre,
		Descripcion: cmd.Descripcion,
		EsSistema:   false,
		TenantID:    cmd.TenantID,
	}

	if err := uc.rolRepo.Crear(ctx, rol); err != nil {
		return nil, fmt.Errorf("error al crear rol: %w", err)
	}

	return &RespuestaCrearRol{
		ID:          rol.ID,
		Nombre:      rol.Nombre,
		Descripcion: rol.Descripcion,
		EsSistema:   rol.EsSistema,
		CreadoEn:    time.Now().Format("2006-01-02T15:04:05Z"),
	}, nil
}

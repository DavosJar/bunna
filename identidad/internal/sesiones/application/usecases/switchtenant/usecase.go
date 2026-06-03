package switchtenant

import (
	"context"
	"errors"
	"time"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

var (
	ErrNoEresMiembro        = errors.New("no eres miembro del tenant solicitado")
	ErrSinRolEnTenant       = errors.New("no tienes un rol asignado en el tenant")
	ErrErrorGenerandoTokens = errors.New("error al generar tokens")
)

// CambiarTenantCasoDeUso maneja el cambio de tenant activo para un usuario.
type CambiarTenantCasoDeUso struct {
	membresiaRepo        tenant.MembresiaRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
	uow                  sesiones_domain.UnitOfWork
}

// NewCambiarTenantCasoDeUso construye el caso de uso con sus dependencias.
func NewCambiarTenantCasoDeUso(
	membresiaRepo tenant.MembresiaRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
	uow sesiones_domain.UnitOfWork,
) *CambiarTenantCasoDeUso {
	return &CambiarTenantCasoDeUso{
		membresiaRepo:        membresiaRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
		uow:                  uow,
	}
}

// Ejecutar cambia el tenant activo del usuario, generando nuevos tokens con el tenant y rol actualizados.
func (uc *CambiarTenantCasoDeUso) Ejecutar(ctx context.Context, cmd ComandoCambiarTenant) (*RespuestaCambiarTenant, error) {
	// Validar que el usuario sea miembro del tenant destino
	esMiembro, err := uc.membresiaRepo.ExisteMiembro(ctx, cmd.UsuarioID, cmd.TenantID)
	if err != nil {
		return nil, err
	}
	if !esMiembro {
		return nil, ErrNoEresMiembro
	}

	// Obtener el rol del usuario en el tenant destino
	roles, err := uc.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, cmd.UsuarioID, cmd.TenantID)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, ErrSinRolEnTenant
	}
	rol := roles[0].Nombre

	var respuesta *RespuestaCambiarTenant

	err = uc.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

		// Obtener la sesión actual para rotar tokens
		sesion, err := tx.SesionRepositorio().ObtenerPorID(ctx, cmd.SesionID)
		if err != nil {
			return err
		}

		// Generar nuevos tokens con el tenant y rol actualizados
		nuevoAccessToken, expiracionAccess, err := tx.TokenServicio().GenerarAccessToken(cmd.UsuarioID, sesion.ID(), cmd.TenantID, rol)
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		nuevoRefreshToken, expiracionRefresh, err := tx.TokenServicio().GenerarRefreshToken(cmd.UsuarioID, sesion.ID())
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		nuevoAccessHash := tx.TokenServicio().HashearToken(nuevoAccessToken)
		nuevoRefreshHash := tx.TokenServicio().HashearToken(nuevoRefreshToken)

		if err := sesion.RotarTokens(nuevoAccessHash, nuevoRefreshHash, expiracionAccess, expiracionRefresh, ahora); err != nil {
			return err
		}

		if _, err := tx.SesionRepositorio().Actualizar(ctx, sesion); err != nil {
			return err
		}

		respuesta = &RespuestaCambiarTenant{
			AccessToken:       nuevoAccessToken,
			RefreshToken:      nuevoRefreshToken,
			ExpiracionAccess:  expiracionAccess,
			ExpiracionRefresh: expiracionRefresh,
			UsuarioID:         cmd.UsuarioID,
			SesionID:          sesion.ID(),
			TenantID:          cmd.TenantID,
			Rol:               rol,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

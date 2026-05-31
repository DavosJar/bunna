package presentation

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	recuperacion "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
	rate_limiter "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/rate_limiter"
	login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	logout "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/logout"
	refresh "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/refresh"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	verificacion "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
)

// MapearError mapea un error de dominio/aplicación a un error HTTP con el status code adecuado.
func MapearError(err error) error {
	if err == nil {
		return nil
	}

	// 400 Bad Request
	if esError(err,
		login.ErrEmailRequerido, login.ErrPasswordRequerido, login.ErrEmailInvalido,
		usuario.ErrIDRequerido, usuario.ErrNombreRequerido, usuario.ErrApellidoRequerido,
		usuario.ErrCorreoRequerido, usuario.ErrTelefonoRequerido,
		usuario.ErrTransicionNoPermitida, usuario.ErrTransicionVerificacionNoPermitida,
		logout.ErrSesionIDRequerido, logout.ErrUsuarioIDRequerido,
		refresh.ErrRefreshTokenRequerido,
		seguridad.ErrIPRequerida,
		tenant.ErrNombreRequerido, tenant.ErrSlugRequerido, tenant.ErrSlugInvalido,
		recuperacion.ErrEmailRequerido, recuperacion.ErrEmailInvalido, recuperacion.ErrPasswordDebil,
	) {
		return huma.Error400BadRequest(err.Error())
	}

	// 401 Unauthorized
	if esError(err,
		login.ErrCredencialesInvalidas,
		refresh.ErrTokenInvalido, refresh.ErrRefreshTokenExpirado,
		logout.ErrNoAutorizado,
		rbac.ErrPasswordActualIncorrecto,
	) {
		return huma.Error401Unauthorized(err.Error())
	}

	// 403 Forbidden
	if esError(err,
		rbac.ErrPermisoDenegado,
		tenant.ErrSinPermiso,
		rbac.ErrSysAdminRequiereTenantVacio,
	) {
		return huma.Error403Forbidden(err.Error())
	}

	// 404 Not Found
	if esError(err,
		logout.ErrSesionNoEncontrada,
		rbac.ErrRolNoEncontrado,
		tenant.ErrTenantNoEncontrado,
		recuperacion.ErrUsuarioNoEncontrado, recuperacion.ErrEnlaceInvalido,
		verificacion.ErrUsuarioNoEncontrado, verificacion.ErrEnlaceInvalido,
	) {
		return huma.Error404NotFound(err.Error())
	}

	// 409 Conflict
	if esError(err,
		rbac.ErrRolYaAsignado, rbac.ErrRolNoAsignado, rbac.ErrRolInmutable,
		tenant.ErrSlugDuplicado,
		tenant.ErrUsuarioYaMiembro, tenant.ErrUsuarioNoEsMiembro, tenant.ErrUltimoAdministrador,
		verificacion.ErrCorreoYaVerificado,
		recuperacion.ErrEnlaceYaUtilizado,
	) {
		return huma.Error409Conflict(err.Error())
	}

	// 422 Unprocessable Entity
	if esError(err,
		login.ErrCuentaBloqueada, login.ErrCuentaInactiva,
		login.ErrCorreoNoVerificado, login.ErrErrorGenerandoTokens,
		refresh.ErrLimiteRefrescosAlcanzado, refresh.ErrSesionAbsolutaExpirada, refresh.ErrSesionNoValida,
		seguridad.ErrIPBloqueada,
		rate_limiter.ErrRateLimitExcedido,
		recuperacion.ErrDemasiadasSolicitudes,
		verificacion.ErrDemasiadosReenvios, verificacion.ErrVerificacionPendiente,
		tenant.ErrTenantInactivo, tenant.ErrTenantYaActivo, tenant.ErrTenantYaInactivo,
		recuperacion.ErrEnlaceExpirado,
		verificacion.ErrEnlaceExpirado,
	) {
		return huma.Error422UnprocessableEntity(err.Error())
	}

	// 500 Internal Server Error (default)
	return huma.Error500InternalServerError(err.Error())
}

// esError retorna true si err (o alguno de los errores en su cadena de wrapping)
// coincide con alguno de los targets.
func esError(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

package login

import (
	"context"
	"errors"
	"net/mail"
	"time"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/rate_limiter"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

var (
	ErrCredencialesInvalidas = errors.New("credenciales inválidas")
	ErrCuentaBloqueada       = errors.New("cuenta temporalmente bloqueada")
	ErrCuentaInactiva        = errors.New("cuenta inactiva")
	ErrEmailRequerido        = errors.New("el email es requerido")
	ErrEmailInvalido         = errors.New("el email no tiene un formato válido")
	ErrPasswordRequerido     = errors.New("la contraseña es requerida")
	ErrErrorGenerandoTokens  = errors.New("error al generar tokens")
	ErrCorreoNoVerificado    = errors.New("debes verificar tu correo electrónico antes de iniciar sesión")
)

type ConfigLogin struct {
	CuentaMaxIntentos     int
	CuentaBloqueoDuracion time.Duration
}

type IniciarSesionCasoDeUso struct {
	uow                  sesiones_domain.UnitOfWork
	bloqueoIP            *bloqueo_ip.ServicioBloqueoIP
	rateLimiter          *rate_limiter.ServicioRateLimit
	config               ConfigLogin
	membresiaRepo        tenant.MembresiaRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
}

func NewIniciarSesionCasoDeUso(
	uow sesiones_domain.UnitOfWork,
	bloqueoIP *bloqueo_ip.ServicioBloqueoIP,
	rateLimiter *rate_limiter.ServicioRateLimit,
	config ConfigLogin,
	membresiaRepo tenant.MembresiaRepositorio,
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio,
) *IniciarSesionCasoDeUso {
	return &IniciarSesionCasoDeUso{
		uow:                  uow,
		bloqueoIP:            bloqueoIP,
		rateLimiter:          rateLimiter,
		config:               config,
		membresiaRepo:        membresiaRepo,
		usuarioTenantRolRepo: usuarioTenantRolRepo,
	}
}

func (uc *IniciarSesionCasoDeUso) Ejecutar(ctx context.Context, cmd ComandoIniciarSesion) (*RespuestaIniciarSesion, error) {
	if err := validarComando(cmd); err != nil {
		return nil, err
	}

	if uc.rateLimiter != nil && cmd.IPOrigen != "" {
		if err := uc.rateLimiter.Verificar(ctx, cmd.IPOrigen); err != nil {
			return nil, err
		}
	}

	if uc.bloqueoIP != nil && cmd.IPOrigen != "" {
		if err := uc.bloqueoIP.Verificar(ctx, cmd.IPOrigen); err != nil {
			return nil, err
		}
	}

	var respuesta *RespuestaIniciarSesion
	var requiereRegistroIP bool

	err := uc.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		ahora := time.Now()

		var credsError error

		usuarios, err := tx.UsuarioRepositorio().Listar(ctx,
			usuario_domain.EspecificacionUsuario{
				ListaLiltros: []domain.CriterioFiltro{
					{Campo: "correo", Operador: "=", Valor: cmd.Email},
				},
			},
			domain.Paginacion{Pagina: 1, TamanoPagina: 1},
		)
		if err != nil || len(usuarios) == 0 {
			credsError = ErrCredencialesInvalidas
		}

		var usuarioID string
		var credenciales *seguridad_domain.CredencialesUsuario
		if credsError == nil {
			usuarioID = usuarios[0].ID()
			credenciales, err = tx.CredencialesRepositorio().ObtenerPorUsuarioID(ctx, usuarioID)
			if err != nil {
				credsError = ErrCredencialesInvalidas
			}
		}

		if credsError == nil && credenciales.EstaBloqueado(ahora) {
			credsError = ErrCuentaBloqueada
		}

		if credsError == nil && !credenciales.Activo() {
			credsError = ErrCuentaInactiva
		}

		// Prevenir user enumeration (ataque de timing) asegurando que siempre se gaste el mismo tiempo computacional
		var passwordCorrecto bool
		if credsError != nil {
			// Simular el tiempo de verificación hasheando el password ingresado
			_, _ = tx.EncriptacionServicio().Hashear(cmd.Password)
		} else {
			passwordCorrecto = tx.EncriptacionServicio().Verificar(cmd.Password, credenciales.PasswordHash())
		}

		if credsError != nil {
			requiereRegistroIP = true
			return credsError
		}

		if !passwordCorrecto {
			credenciales.IncrementarIntentoFallido()
			if credenciales.IntentosFallidos() >= uc.config.CuentaMaxIntentos {
				credenciales.BloquearHasta(ahora.Add(uc.config.CuentaBloqueoDuracion))
			}
			if _, err := tx.CredencialesRepositorio().Actualizar(ctx, credenciales); err != nil {
				return err
			}
			requiereRegistroIP = true
			return ErrCredencialesInvalidas
		}

		// Verificar que el correo esté verificado
		if usuarios[0].EstadoVerificacionCorreo() != usuario_domain.VERIFICADO {
			return ErrCorreoNoVerificado
		}

		tenantID := ""
		rol := ""
		if uc.membresiaRepo != nil && uc.usuarioTenantRolRepo != nil {
			tenants, err := uc.membresiaRepo.ListarTenantsPorUsuario(ctx, usuarioID)
			if err == nil && len(tenants) > 0 {
				tenantID = tenants[0]
				roles, err := uc.usuarioTenantRolRepo.ListarRolesPorUsuarioEnTenant(ctx, usuarioID, tenantID)
				if err == nil && len(roles) > 0 {
					rol = roles[0].Nombre
				}
			}
		}

		sesionID, err := tx.GeneradorID().NextID(ctx)
		if err != nil {
			return err
		}

		accessToken, expiracionAccess, err := tx.TokenServicio().GenerarAccessToken(usuarioID, sesionID, tenantID, rol)
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		refreshToken, expiracionRefresh, err := tx.TokenServicio().GenerarRefreshToken(usuarioID, sesionID)
		if err != nil {
			return ErrErrorGenerandoTokens
		}

		accessTokenHash := tx.TokenServicio().HashearToken(accessToken)
		refreshTokenHash := tx.TokenServicio().HashearToken(refreshToken)

		sesion, err := sesiones_domain.NuevaSesion(
			sesionID,
			usuarioID,
			accessTokenHash,
			refreshTokenHash,
			cmd.IPOrigen,
			ahora,
			expiracionAccess,
			expiracionRefresh,
		)
		if err != nil {
			return err
		}

		if _, err := tx.SesionRepositorio().Crear(ctx, sesion); err != nil {
			return err
		}

		credenciales.ResetearIntentos()
		if _, err := tx.CredencialesRepositorio().Actualizar(ctx, credenciales); err != nil {
			return err
		}

		respuesta = &RespuestaIniciarSesion{
			AccessToken:       accessToken,
			RefreshToken:      refreshToken,
			ExpiracionAccess:  expiracionAccess,
			ExpiracionRefresh: expiracionRefresh,
			UsuarioID:         usuarioID,
			SesionID:          sesionID,
			TenantID:          tenantID,
			Rol:               rol,
		}
		return nil
	})

	if requiereRegistroIP && uc.bloqueoIP != nil && cmd.IPOrigen != "" {
		_ = uc.bloqueoIP.RegistrarIntentoFallido(ctx, cmd.IPOrigen)
	}

	if err != nil {
		return nil, err
	}
	return respuesta, nil
}

func validarComando(cmd ComandoIniciarSesion) error {
	if cmd.Email == "" {
		return ErrEmailRequerido
	}
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return ErrEmailInvalido
	}
	if cmd.Password == "" {
		return ErrPasswordRequerido
	}
	return nil
}

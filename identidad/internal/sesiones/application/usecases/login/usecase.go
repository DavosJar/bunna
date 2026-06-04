package login

import (
	"context"
	"errors"
	"net/mail"
	"time"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/bloqueo_ip"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/application/services/rate_limiter"
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

	ahora := time.Now()
	var requiereRegistroIP bool

	// ============================================================
	// Phase 1: Credential Validation (NO transaction)
	// The attempt counter MUST persist even on failure, which is
	// why this runs outside the GORM transaction.
	// ============================================================

	usuarios, err := uc.uow.UsuarioRepositorio().Listar(ctx,
		usuario_domain.EspecificacionUsuario{
			ListaLiltros: []domain.CriterioFiltro{
				{Campo: "correo", Operador: "=", Valor: cmd.Email},
			},
		},
		domain.Paginacion{Pagina: 1, TamanoPagina: 1},
	)
	if err != nil || len(usuarios) == 0 {
		// Timing attack prevention: hash the password to keep
		// computation time consistent whether the user exists or not.
		_, _ = uc.uow.EncriptacionServicio().Hashear(cmd.Password)
		requiereRegistroIP = true
		return nil, ErrCredencialesInvalidas
	}

	usuario := usuarios[0]
	usuarioID := usuario.ID()

	credenciales, err := uc.uow.CredencialesRepositorio().ObtenerPorUsuarioID(ctx, usuarioID)
	if err != nil {
		// Timing attack prevention: same as above.
		_, _ = uc.uow.EncriptacionServicio().Hashear(cmd.Password)
		requiereRegistroIP = true
		return nil, ErrCredencialesInvalidas
	}

	// Check if the account is temporarily blocked.
	if credenciales.EstaBloqueado(ahora) {
		// Timing attack prevention.
		_, _ = uc.uow.EncriptacionServicio().Hashear(cmd.Password)
		requiereRegistroIP = true
		return nil, ErrCuentaBloqueada
	}

	// Check if the account is permanently inactive.
	if !credenciales.Activo() {
		// Timing attack prevention.
		_, _ = uc.uow.EncriptacionServicio().Hashear(cmd.Password)
		requiereRegistroIP = true
		return nil, ErrCuentaInactiva
	}

	// Verify the password.
	passwordCorrecto := uc.uow.EncriptacionServicio().Verificar(cmd.Password, credenciales.PasswordHash())

	if !passwordCorrecto {
		credenciales.IncrementarIntentoFallido()
		if credenciales.IntentosFallidos() >= uc.config.CuentaMaxIntentos {
			credenciales.BloquearHasta(ahora.Add(uc.config.CuentaBloqueoDuracion))
		}
		if _, err := uc.uow.CredencialesRepositorio().Actualizar(ctx, credenciales); err != nil {
			return nil, err
		}
		requiereRegistroIP = true
		if credenciales.EstaBloqueado(ahora) {
			return nil, ErrCuentaBloqueada
		}
		return nil, ErrCredencialesInvalidas
	}

	// ============================================================
	// Phase 2: Session Creation (WITH transaction)
	// Only reached after successful password verification.
	// If anything fails here the transaction rolls back, but the
	// failed-attempt counter from Phase 1 stays persisted.
	// ============================================================

	var respuesta *RespuestaIniciarSesion

	err = uc.uow.Transaccional(ctx, func(tx sesiones_domain.UnitOfWork) error {
		// Verify that the user's email is verified before creating a session.
		if usuario.EstadoVerificacionCorreo() != usuario_domain.VERIFICADO {
			return ErrCorreoNoVerificado
		}

		// Reload credenciales inside the transaction to get a consistent
		// view and avoid carrying over in-memory mutations from Phase 1.
		credencialesTx, err := tx.CredencialesRepositorio().ObtenerPorUsuarioID(ctx, usuarioID)
		if err != nil {
			return err
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

		credencialesTx.ResetearIntentos()
		if _, err := tx.CredencialesRepositorio().Actualizar(ctx, credencialesTx); err != nil {
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

package crearInvitacion

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type CrearInvitacionCasoDeUso struct {
	invitacionRepo   invitaciones.InvitacionRepositorio
	tenantRepo       tenant.TenantRepositorio
	rolRepo          rbac.RolRepositorio
	emailSvc         notificaciones.EmailServicio
	idGen            shareddomain.GeneradorID
	frontendURL      string
	tokenExpiracion  time.Duration
}

func NewCrearInvitacionCasoDeUso(
	invitacionRepo invitaciones.InvitacionRepositorio,
	tenantRepo tenant.TenantRepositorio,
	rolRepo rbac.RolRepositorio,
	emailSvc notificaciones.EmailServicio,
	idGen shareddomain.GeneradorID,
	frontendURL string,
	tokenExpiracion time.Duration,
) *CrearInvitacionCasoDeUso {
	return &CrearInvitacionCasoDeUso{
		invitacionRepo:  invitacionRepo,
		tenantRepo:      tenantRepo,
		rolRepo:         rolRepo,
		emailSvc:        emailSvc,
		idGen:           idGen,
		frontendURL:     frontendURL,
		tokenExpiracion: tokenExpiracion,
	}
}

func (uc *CrearInvitacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoCrearInvitacion) (*RespuestaCrearInvitacion, error) {
	if cmd.Correo == "" {
		return nil, invitaciones.ErrEmailRequerido
	}
	if _, err := mail.ParseAddress(cmd.Correo); err != nil {
		return nil, fmt.Errorf("formato de correo inválido: %w", err)
	}
	if cmd.RolID == "" {
		return nil, invitaciones.ErrRolRequerido
	}

	tenantObj, err := uc.tenantRepo.ObtenerPorID(ctx, cmd.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant no encontrado: %w", err)
	}

	if _, err := uc.rolRepo.ObtenerPorID(ctx, cmd.RolID); err != nil {
		return nil, fmt.Errorf("rol no encontrado: %w", err)
	}

	id, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID: %w", err)
	}

	token, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	expiracion := time.Now().Add(uc.tokenExpiracion)

	invitacion := invitaciones.NuevaInvitacion(id, cmd.TenantID, cmd.RolID, cmd.Correo, cmd.Nombre, token, expiracion)

	if err := uc.invitacionRepo.Crear(ctx, invitacion); err != nil {
		return nil, fmt.Errorf("error al crear invitación: %w", err)
	}

	nombreTenant := tenantObj.Nombre()
	urlInvitacion := fmt.Sprintf("%s/aceptar-invitacion?token=%s", uc.frontendURL, token)
	expiracionHoras := fmt.Sprintf("%.0f", uc.tokenExpiracion.Hours())

	if err := uc.emailSvc.EnviarTemplate(ctx, cmd.Correo, notificaciones.TipoInvitacion, map[string]string{
		"nombre_tenant":  nombreTenant,
		"url_invitacion": urlInvitacion,
		"expiracion_horas": expiracionHoras,
	}); err != nil {
		fmt.Printf("[CrearInvitacion] Error al enviar email: %v\n", err)
	}

	return &RespuestaCrearInvitacion{
		ID:    id,
		Token: token,
	}, nil
}

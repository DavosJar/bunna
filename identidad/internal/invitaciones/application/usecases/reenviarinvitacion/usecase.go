package reenviarinvitacion

import (
	"context"
	"fmt"
	"time"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	tenant "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type ReenviarInvitacionCasoDeUso struct {
	invitacionRepo  invitaciones.InvitacionRepositorio
	tenantRepo      tenant.TenantRepositorio
	emailSvc        notificaciones.EmailServicio
	idGen           shareddomain.GeneradorID
	frontendURL     string
	tokenExpiracion time.Duration
}

func NewReenviarInvitacionCasoDeUso(
	invitacionRepo invitaciones.InvitacionRepositorio,
	tenantRepo tenant.TenantRepositorio,
	emailSvc notificaciones.EmailServicio,
	idGen shareddomain.GeneradorID,
	frontendURL string,
	tokenExpiracion time.Duration,
) *ReenviarInvitacionCasoDeUso {
	return &ReenviarInvitacionCasoDeUso{
		invitacionRepo:  invitacionRepo,
		tenantRepo:      tenantRepo,
		emailSvc:        emailSvc,
		idGen:           idGen,
		frontendURL:     frontendURL,
		tokenExpiracion: tokenExpiracion,
	}
}

func (uc *ReenviarInvitacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoReenviarInvitacion) (*RespuestaReenviarInvitacion, error) {
	if cmd.InvitacionID == "" {
		return nil, invitaciones.ErrEnlaceInvalido
	}

	invitacion, err := uc.invitacionRepo.ObtenerPorID(ctx, cmd.InvitacionID)
	if err != nil {
		return nil, err
	}

	// Verificar que la invitación pertenezca al tenant del ejecutor
	if invitacion.TenantID() != cmd.TenantID {
		return nil, invitaciones.ErrNoEncontrada
	}

	// No se puede reenviar si ya fue aceptada
	if invitacion.EstaAceptada() {
		return nil, invitaciones.ErrYaAceptada
	}

	// Obtener datos del tenant para el nombre en el email
	tenantObj, err := uc.tenantRepo.ObtenerPorID(ctx, cmd.TenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant no encontrado: %w", err)
	}

	// Generar nuevo token
	nuevoToken, err := uc.idGen.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	// Actualizar token en BD
	nuevoHash := invitaciones.HashearTokenPublico(nuevoToken)
	if err := uc.invitacionRepo.ActualizarTokenHash(ctx, cmd.InvitacionID, nuevoHash); err != nil {
		return nil, fmt.Errorf("error al actualizar token: %w", err)
	}

	// Enviar email con el nuevo enlace
	urlInvitacion := fmt.Sprintf("%s/aceptar-invitacion?token=%s", uc.frontendURL, nuevoToken)
	expiracionHoras := fmt.Sprintf("%.0f", uc.tokenExpiracion.Hours())

	if err := uc.emailSvc.EnviarTemplate(ctx, invitacion.Email(), notificaciones.TipoInvitacion, map[string]string{
		"nombre_tenant":   tenantObj.Nombre(),
		"url_invitacion":  urlInvitacion,
		"expiracion_horas": expiracionHoras,
	}); err != nil {
		fmt.Printf("[ReenviarInvitacion] Error al enviar email: %v\n", err)
	}

	return &RespuestaReenviarInvitacion{
		Mensaje: "Invitación reenviada exitosamente",
	}, nil
}

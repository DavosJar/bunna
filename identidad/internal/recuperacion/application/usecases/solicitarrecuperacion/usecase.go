package solicitarrecuperacion

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	dominio "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type ConfigRecuperacion struct {
	TokenExpiracion     time.Duration
	RateLimitIPMax      int
	RateLimitUsuarioMax int
	RateLimitVentana    time.Duration
	FrontendURL         string
}

type SolicitarRecuperacionCasoDeUso struct {
	tokenRepo     dominio.TokenRecuperacionRepositorio
	usuarioRepo   dominio.UsuarioRecuperacionRepositorio
	emailServicio notificaciones.EmailServicio
	idGenerator   shareddomain.GeneradorID
	config        ConfigRecuperacion
}

func NewSolicitarRecuperacionCasoDeUso(
	tokenRepo dominio.TokenRecuperacionRepositorio,
	usuarioRepo dominio.UsuarioRecuperacionRepositorio,
	emailServicio notificaciones.EmailServicio,
	idGenerator shareddomain.GeneradorID,
	config ConfigRecuperacion,
) *SolicitarRecuperacionCasoDeUso {
	if config.TokenExpiracion == 0 {
		config.TokenExpiracion = time.Hour
	}
	if config.RateLimitIPMax == 0 {
		config.RateLimitIPMax = 3
	}
	if config.RateLimitUsuarioMax == 0 {
		config.RateLimitUsuarioMax = 1
	}
	if config.RateLimitVentana == 0 {
		config.RateLimitVentana = 15 * time.Minute
	}
	return &SolicitarRecuperacionCasoDeUso{
		tokenRepo:     tokenRepo,
		usuarioRepo:   usuarioRepo,
		emailServicio: emailServicio,
		idGenerator:   idGenerator,
		config:        config,
	}
}

func (uc *SolicitarRecuperacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoSolicitarRecuperacion) (*RespuestaSolicitarRecuperacion, error) {
	respuestaGenerica := &RespuestaSolicitarRecuperacion{
		Mensaje: "Si el email existe, recibirás un enlace de recuperación",
	}

	if cmd.Email == "" {
		return nil, dominio.ErrEmailRequerido
	}
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return nil, dominio.ErrEmailInvalido
	}

	usuario, err := uc.usuarioRepo.ObtenerPorCorreo(ctx, cmd.Email)
	if err != nil {
		return respuestaGenerica, nil
	}

	tokenID, err := uc.idGenerator.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID: %w", err)
	}
	tokenPlano, err := uc.idGenerator.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	expiraEn := time.Now().Add(uc.config.TokenExpiracion)
	token := dominio.NuevoTokenRecuperacion(tokenID, usuario.ID, tokenPlano, expiraEn)
	if err := uc.tokenRepo.Crear(ctx, token); err != nil {
		return nil, fmt.Errorf("error al persistir token: %w", err)
	}

	expiracionHoras := fmt.Sprintf("%.0f", uc.config.TokenExpiracion.Hours())
	urlRecuperacion := fmt.Sprintf("%s/reset-password?token=%s", uc.config.FrontendURL, tokenPlano)
	go func() {
		if err := uc.emailServicio.EnviarTemplate(ctx, usuario.Correo,
			notificaciones.TipoRecuperacionContrasena,
			map[string]string{
				"nombre":           usuario.Nombre,
				"token":            tokenPlano,
				"expiracion_horas": expiracionHoras,
				"url_recuperacion": urlRecuperacion,
			},
		); err != nil {
			fmt.Printf("[SolicitarRecuperacionCasoDeUso] Error al enviar email: %v\n", err)
		}
	}()

	return respuestaGenerica, nil
}

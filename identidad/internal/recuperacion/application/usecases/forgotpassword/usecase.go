package forgotpassword

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	dominio "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/shared/application"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type ConfigRecuperacion struct {
	TokenExpiracion     time.Duration
	RateLimitIPMax      int
	RateLimitUsuarioMax int
	RateLimitVentana    time.Duration
	FrontendURL         string
}

type RecuperarContrasenaCasoDeUso struct {
	tokenRepo     dominio.TokenRecuperacionRepositorio
	usuarioRepo   dominio.UsuarioRecuperacionRepositorio
	sesionRepo    sesiones_domain.SesionRepositorio
	credRepo      seguridad_domain.CredencialesRepositorio
	encriptacion  seguridad_domain.EncriptacionServicio
	emailServicio notificaciones.EmailServicio
	idGenerator   shareddomain.GeneradorID
	config        ConfigRecuperacion
}

func NewRecuperarContrasenaCasoDeUso(
	tokenRepo dominio.TokenRecuperacionRepositorio,
	usuarioRepo dominio.UsuarioRecuperacionRepositorio,
	sesionRepo sesiones_domain.SesionRepositorio,
	credRepo seguridad_domain.CredencialesRepositorio,
	encriptacion seguridad_domain.EncriptacionServicio,
	emailServicio notificaciones.EmailServicio,
	idGenerator shareddomain.GeneradorID,
	config ConfigRecuperacion,
) *RecuperarContrasenaCasoDeUso {
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
	return &RecuperarContrasenaCasoDeUso{
		tokenRepo:     tokenRepo,
		usuarioRepo:   usuarioRepo,
		sesionRepo:    sesionRepo,
		credRepo:      credRepo,
		encriptacion:  encriptacion,
		emailServicio: emailServicio,
		idGenerator:   idGenerator,
		config:        config,
	}
}

func (uc *RecuperarContrasenaCasoDeUso) Solicitar(ctx context.Context, cmd ComandoSolicitarRecuperacion) (*RespuestaSolicitarRecuperacion, error) {
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
			fmt.Printf("[RecuperarContrasenaCasoDeUso] Error al enviar email: %v\n", err)
		}
	}()

	return respuestaGenerica, nil
}

func (uc *RecuperarContrasenaCasoDeUso) ValidarToken(ctx context.Context, cmd ComandoValidarTokenRecuperacion) (*RespuestaValidarTokenRecuperacion, error) {
	if cmd.Token == "" {
		return nil, dominio.ErrEnlaceInvalido
	}

	hash := dominio.HashearToken(cmd.Token)
	token, err := uc.tokenRepo.ObtenerPorHash(ctx, hash)
	if err != nil {
		return nil, dominio.ErrEnlaceInvalido
	}

	if err := token.EsValido(time.Now()); err != nil {
		return nil, err
	}

	return &RespuestaValidarTokenRecuperacion{
		UsuarioID: token.UsuarioID(),
		Valido:    true,
	}, nil
}

func (uc *RecuperarContrasenaCasoDeUso) Confirmar(ctx context.Context, cmd ComandoConfirmarRestablecimiento) (*RespuestaConfirmarRestablecimiento, error) {
	if err := application.ValidarFormatoPassword(cmd.NuevaPassword, "nueva_password"); err != nil {
		return nil, dominio.ErrPasswordDebil
	}

	respuestaValidar, err := uc.ValidarToken(ctx, ComandoValidarTokenRecuperacion{Token: cmd.Token})
	if err != nil {
		return nil, err
	}

	nuevoHash, err := uc.encriptacion.Hashear(cmd.NuevaPassword)
	if err != nil {
		return nil, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	if err := uc.usuarioRepo.ActualizarPassword(ctx, respuestaValidar.UsuarioID, nuevoHash); err != nil {
		return nil, fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	hash := dominio.HashearToken(cmd.Token)
	token, _ := uc.tokenRepo.ObtenerPorHash(ctx, hash)
	token.Usar(time.Now())
	if err := uc.tokenRepo.Actualizar(ctx, token); err != nil {
		fmt.Printf("[RecuperarContrasenaCasoDeUso] Error al marcar token usado: %v\n", err)
	}

	if err := uc.sesionRepo.InvalidarTodasPorUsuarioID(ctx, respuestaValidar.UsuarioID); err != nil {
		fmt.Printf("[RecuperarContrasenaCasoDeUso] Error al invalidar sesiones: %v\n", err)
	}

	return &RespuestaConfirmarRestablecimiento{
		Mensaje: "Contraseña actualizada exitosamente",
	}, nil
}

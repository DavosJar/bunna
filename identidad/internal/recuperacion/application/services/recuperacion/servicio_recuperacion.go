package recuperacion

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	dominio "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

// ConfigRecuperacion contiene parámetros configurables
type ConfigRecuperacion struct {
	TokenExpiracion       time.Duration
	RateLimitIPMax        int
	RateLimitUsuarioMax   int
	RateLimitVentana      time.Duration
}

// ServicioRecuperacion maneja los casos de uso de recuperación de contraseña
type ServicioRecuperacion struct {
	tokenRepo     dominio.TokenRecuperacionRepositorio
	usuarioRepo   dominio.UsuarioRecuperacionRepositorio
	sesionRepo    sesiones_domain.SesionRepositorio
	credRepo      seguridad_domain.CredencialesRepositorio
	encriptacion  seguridad_domain.EncriptacionServicio
	emailServicio notificaciones.EmailServicio
	idGenerator   shareddomain.GeneradorID
	config        ConfigRecuperacion
}

// NuevoServicioRecuperacion crea una nueva instancia del servicio
func NuevoServicioRecuperacion(
	tokenRepo dominio.TokenRecuperacionRepositorio,
	usuarioRepo dominio.UsuarioRecuperacionRepositorio,
	sesionRepo sesiones_domain.SesionRepositorio,
	credRepo seguridad_domain.CredencialesRepositorio,
	encriptacion seguridad_domain.EncriptacionServicio,
	emailServicio notificaciones.EmailServicio,
	idGenerator shareddomain.GeneradorID,
	config ConfigRecuperacion,
) *ServicioRecuperacion {
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
	return &ServicioRecuperacion{
		tokenRepo:    tokenRepo,
		usuarioRepo:  usuarioRepo,
		sesionRepo:   sesionRepo,
		credRepo:     credRepo,
		encriptacion: encriptacion,
		emailServicio: emailServicio,
		idGenerator:  idGenerator,
		config:       config,
	}
}

// SolicitarRecuperacion genera token y envía email (respuesta genérica siempre)
func (s *ServicioRecuperacion) SolicitarRecuperacion(ctx context.Context, cmd ComandoSolicitarRecuperacion) (*RespuestaSolicitar, error) {
	respuestaGenerica := &RespuestaSolicitar{
		Mensaje: "Si el email existe, recibirás un enlace de recuperación",
	}

	// Validar email
	if cmd.Email == "" {
		return nil, dominio.ErrEmailRequerido
	}
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return nil, dominio.ErrEmailInvalido
	}

	// Buscar usuario (silencioso si no existe)
	usuario, err := s.usuarioRepo.ObtenerPorCorreo(ctx, cmd.Email)
	if err != nil {
		return respuestaGenerica, nil
	}

	// Generar token
	tokenID, err := s.idGenerator.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID: %w", err)
	}
	tokenPlano, err := s.idGenerator.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	// Crear y persistir token
	expiraEn := time.Now().Add(s.config.TokenExpiracion)
	token := dominio.NuevoTokenRecuperacion(tokenID, usuario.ID, tokenPlano, expiraEn)
	if err := s.tokenRepo.Crear(ctx, token); err != nil {
		return nil, fmt.Errorf("error al persistir token: %w", err)
	}

	// Enviar email asíncrono (best-effort)
	expiracionHoras := fmt.Sprintf("%.0f", s.config.TokenExpiracion.Hours())
	go func() {
		if err := s.emailServicio.EnviarTemplate(ctx, usuario.Correo,
			notificaciones.TipoRecuperacionContrasena,
			map[string]string{
				"nombre":           usuario.Nombre,
				"token":            tokenPlano,
				"expiracion_horas": expiracionHoras,
			},
		); err != nil {
			fmt.Printf("[RecuperacionServicio] Error al enviar email: %v\n", err)
		}
	}()

	return respuestaGenerica, nil
}

// ValidarToken verifica que el token sea válido
func (s *ServicioRecuperacion) ValidarToken(ctx context.Context, cmd ComandoValidarToken) (*RespuestaValidar, error) {
	if cmd.Token == "" {
		return nil, dominio.ErrEnlaceInvalido
	}

	hash := dominio.HashearToken(cmd.Token)
	token, err := s.tokenRepo.ObtenerPorHash(ctx, hash)
	if err != nil {
		return nil, dominio.ErrEnlaceInvalido
	}

	if err := token.EsValido(time.Now()); err != nil {
		return nil, err
	}

	return &RespuestaValidar{
		UsuarioID: token.UsuarioID(),
		Valido:    true,
	}, nil
}

// ConfirmarRestablecimiento cambia la contraseña usando el token
func (s *ServicioRecuperacion) ConfirmarRestablecimiento(ctx context.Context, cmd ComandoConfirmarRestablecimiento) (*RespuestaConfirmar, error) {
	// Validar password
	if cmd.NuevoPassword == "" {
		return nil, dominio.ErrPasswordDebil
	}
	if len(cmd.NuevoPassword) < 8 {
		return nil, dominio.ErrPasswordDebil
	}

	// Validar token
	respuestaValidar, err := s.ValidarToken(ctx, ComandoValidarToken{Token: cmd.Token})
	if err != nil {
		return nil, err
	}

	// Hashear nueva contraseña
	nuevoHash, err := s.encriptacion.Hashear(cmd.NuevoPassword)
	if err != nil {
		return nil, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	// Actualizar contraseña
	if err := s.usuarioRepo.ActualizarPassword(ctx, respuestaValidar.UsuarioID, nuevoHash); err != nil {
		return nil, fmt.Errorf("error al actualizar contraseña: %w", err)
	}

	// Marcar token como usado
	hash := dominio.HashearToken(cmd.Token)
	token, _ := s.tokenRepo.ObtenerPorHash(ctx, hash)
	token.Usar(time.Now())
	if err := s.tokenRepo.Actualizar(ctx, token); err != nil {
		fmt.Printf("[RecuperacionServicio] Error al marcar token usado: %v\n", err)
	}

	// Invalidar todas las sesiones del usuario
	if err := s.sesionRepo.InvalidarTodasPorUsuarioID(ctx, respuestaValidar.UsuarioID); err != nil {
		fmt.Printf("[RecuperacionServicio] Error al invalidar sesiones: %v\n", err)
	}

	return &RespuestaConfirmar{
		Mensaje: "Contraseña actualizada exitosamente",
	}, nil
}

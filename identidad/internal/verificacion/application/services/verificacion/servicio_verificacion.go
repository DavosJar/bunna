package verificacion

import (
	"context"
	"fmt"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	dominio "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
)

type ConfigVerificacion struct {
	TokenExpiracion time.Duration
	MaxReenvios     int
	VentanaReenvios time.Duration
}

type ServicioVerificacion struct {
	repo          dominio.VerificacionRepositorio
	emailServicio notificaciones.EmailServicio
	idGenerator   shareddomain.GeneradorID
	config        ConfigVerificacion
}

func NuevoServicioVerificacion(
	repo dominio.VerificacionRepositorio,
	emailServicio notificaciones.EmailServicio,
	idGenerator shareddomain.GeneradorID,
	config ConfigVerificacion,
) *ServicioVerificacion {
	if config.TokenExpiracion == 0 {
		config.TokenExpiracion = 24 * time.Hour
	}
	if config.MaxReenvios == 0 {
		config.MaxReenvios = 5
	}
	if config.VentanaReenvios == 0 {
		config.VentanaReenvios = 24 * time.Hour
	}
	return &ServicioVerificacion{
		repo:          repo,
		emailServicio: emailServicio,
		idGenerator:   idGenerator,
		config:        config,
	}
}

func (s *ServicioVerificacion) SolicitarVerificacion(ctx context.Context, cmd ComandoSolicitarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	usuario, err := s.repo.ObtenerPorID(ctx, cmd.UsuarioID)
	if err != nil {
		return nil, dominio.ErrUsuarioNoEncontrado
	}

	if usuario.EstadoVerificacion == "VERIFICADO" {
		return nil, dominio.ErrCorreoYaVerificado
	}

	token, err := s.idGenerator.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	expiraEn := time.Now().Add(s.config.TokenExpiracion)
	prueba := dominio.NuevaPruebaVerificacion(token, expiraEn)

	if err := s.repo.ActualizarPrueba(ctx, cmd.UsuarioID, prueba); err != nil {
		return nil, fmt.Errorf("error al persistir token: %w", err)
	}

	expiracionHoras := fmt.Sprintf("%.0f", s.config.TokenExpiracion.Hours())
	if err := s.emailServicio.EnviarTemplate(ctx, usuario.Correo,
		notificaciones.TipoVerificacionCorreo,
		map[string]string{
			"nombre":           usuario.Nombre,
			"token":            token,
			"expiracion_horas": expiracionHoras,
		},
	); err != nil {
		fmt.Printf("[VerificacionServicio] Error al enviar email: %v\n", err)
		return nil, fmt.Errorf("error al enviar email de verificación: %w", err)
	}

	return &RespuestaSolicitarVerificacion{
		Mensaje: "Email de verificación enviado",
	}, nil
}

func (s *ServicioVerificacion) ConfirmarVerificacion(ctx context.Context, cmd ComandoConfirmarVerificacion) (*RespuestaConfirmarVerificacion, error) {
	if cmd.Token == "" {
		return nil, dominio.ErrEnlaceInvalido
	}

	hash := dominio.HashearToken(cmd.Token)

	usuario, err := s.repo.ObtenerPorHashToken(ctx, hash)
	if err != nil {
		return nil, dominio.ErrEnlaceInvalido
	}

	if usuario.PruebaVerificacion.Expiro(time.Now()) {
		if err := s.repo.ActualizarPrueba(ctx, usuario.ID, dominio.PruebaVerificacionVacia()); err != nil {
			fmt.Printf("[VerificacionServicio] Error al limpiar prueba: %v\n", err)
		}
		return nil, dominio.ErrEnlaceExpirado
	}

	if err := s.repo.MarcarVerificado(ctx, usuario.ID); err != nil {
		return nil, fmt.Errorf("error al marcar verificado: %w", err)
	}

	return &RespuestaConfirmarVerificacion{
		Mensaje: "Correo verificado exitosamente",
	}, nil
}

func (s *ServicioVerificacion) ReenviarVerificacion(ctx context.Context, cmd ComandoReenviarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	usuario, err := s.repo.ObtenerPorID(ctx, cmd.UsuarioID)
	if err != nil {
		return nil, dominio.ErrUsuarioNoEncontrado
	}

	if usuario.EstadoVerificacion == "VERIFICADO" {
		return nil, dominio.ErrCorreoYaVerificado
	}

	if usuario.ContadorReenvios >= s.config.MaxReenvios {
		return nil, dominio.ErrDemasiadosReenvios
	}

	return s.SolicitarVerificacion(ctx, ComandoSolicitarVerificacion{
		UsuarioID: cmd.UsuarioID,
	})
}
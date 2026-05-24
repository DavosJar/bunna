package verificacion

import (
	"context"
	"fmt"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	dominio "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
)

// ConfigVerificacion contiene parámetros configurables del servicio
type ConfigVerificacion struct {
	TokenExpiracion time.Duration
	MaxReenvios     int
	VentanaReenvios time.Duration
}

// ServicioVerificacion maneja los casos de uso de verificación de correo
type ServicioVerificacion struct {
	repo          dominio.VerificacionRepositorio
	emailServicio notificaciones.EmailServicio
	idGenerator   shareddomain.GeneradorID
	config        ConfigVerificacion
}

// NuevoServicioVerificacion crea una nueva instancia del servicio
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

// SolicitarVerificacion genera token y envía email de verificación
func (s *ServicioVerificacion) SolicitarVerificacion(ctx context.Context, cmd ComandoSolicitarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	usuario, err := s.repo.ObtenerPorID(ctx, cmd.UsuarioID)
	if err != nil {
		return nil, dominio.ErrUsuarioNoEncontrado
	}

	if usuario.EstadoVerificacion == "VERIFICADO" {
		return nil, dominio.ErrCorreoYaVerificado
	}

	// Generar token único
	token, err := s.idGenerator.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	// Crear prueba de verificación
	expiraEn := time.Now().Add(s.config.TokenExpiracion)
	prueba := dominio.NuevaPruebaVerificacion(token, expiraEn)

	// Persistir hash
	if err := s.repo.ActualizarPrueba(ctx, cmd.UsuarioID, prueba); err != nil {
		return nil, fmt.Errorf("error al persistir token: %w", err)
	}

	// Enviar email de forma asíncrona (best-effort)
	expiracionHoras := fmt.Sprintf("%.0f", s.config.TokenExpiracion.Hours())
	go func() {
		if err := s.emailServicio.EnviarTemplate(ctx, usuario.Correo,
			notificaciones.TipoVerificacionCorreo,
			map[string]string{
				"nombre":           usuario.Nombre,
				"token":            token,
				"expiracion_horas": expiracionHoras,
			},
		); err != nil {
			// Log pero no falla la operación
			fmt.Printf("[VerificacionServicio] Error al enviar email: %v\n", err)
		}
	}()

	return &RespuestaSolicitarVerificacion{
		Mensaje: "Email de verificación enviado",
	}, nil
}

// ConfirmarVerificacion valida el token y marca el correo como verificado
func (s *ServicioVerificacion) ConfirmarVerificacion(ctx context.Context, cmd ComandoConfirmarVerificacion) (*RespuestaConfirmarVerificacion, error) {
	if cmd.Token == "" {
		return nil, dominio.ErrEnlaceInvalido
	}

	// Hashear token recibido
	hash := dominio.HashearToken(cmd.Token)

	// Buscar usuario por hash
	usuario, err := s.repo.ObtenerPorHashToken(ctx, hash)
	if err != nil {
		return nil, dominio.ErrEnlaceInvalido
	}

	// Verificar expiración
	if usuario.PruebaVerificacion.Expiro(time.Now()) {
		if err := s.repo.ActualizarPrueba(ctx, usuario.ID, dominio.PruebaVerificacionVacia()); err != nil {
			fmt.Printf("[VerificacionServicio] Error al limpiar prueba: %v\n", err)
		}
		return nil, dominio.ErrEnlaceExpirado
	}

	// Marcar como verificado
	if err := s.repo.MarcarVerificado(ctx, usuario.ID); err != nil {
		return nil, fmt.Errorf("error al marcar verificado: %w", err)
	}

	return &RespuestaConfirmarVerificacion{
		Mensaje: "Correo verificado exitosamente",
	}, nil
}

// ReenviarVerificacion genera nuevo token y reenvía el email
func (s *ServicioVerificacion) ReenviarVerificacion(ctx context.Context, cmd ComandoReenviarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	usuario, err := s.repo.ObtenerPorID(ctx, cmd.UsuarioID)
	if err != nil {
		return nil, dominio.ErrUsuarioNoEncontrado
	}

	if usuario.EstadoVerificacion == "VERIFICADO" {
		return nil, dominio.ErrCorreoYaVerificado
	}

	// Verificar límite de reenvíos
	if usuario.ContadorReenvios >= s.config.MaxReenvios {
		return nil, dominio.ErrDemasiadosReenvios
	}

	// Reutilizar SolicitarVerificacion
	return s.SolicitarVerificacion(ctx, ComandoSolicitarVerificacion{
		UsuarioID: cmd.UsuarioID,
	})
}

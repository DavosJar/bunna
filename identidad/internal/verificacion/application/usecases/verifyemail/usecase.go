package verifyemail

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

type VerificarCorreoCasoDeUso struct {
	repo          dominio.VerificacionRepositorio
	emailServicio notificaciones.EmailServicio
	idGenerator   shareddomain.GeneradorID
	config        ConfigVerificacion
}

func NewVerificarCorreoCasoDeUso(
	repo dominio.VerificacionRepositorio,
	emailServicio notificaciones.EmailServicio,
	idGenerator shareddomain.GeneradorID,
	config ConfigVerificacion,
) *VerificarCorreoCasoDeUso {
	if config.TokenExpiracion == 0 {
		config.TokenExpiracion = 24 * time.Hour
	}
	if config.MaxReenvios == 0 {
		config.MaxReenvios = 5
	}
	if config.VentanaReenvios == 0 {
		config.VentanaReenvios = 24 * time.Hour
	}
	return &VerificarCorreoCasoDeUso{
		repo:          repo,
		emailServicio: emailServicio,
		idGenerator:   idGenerator,
		config:        config,
	}
}

func (uc *VerificarCorreoCasoDeUso) Solicitar(ctx context.Context, cmd ComandoSolicitarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	usuario, err := uc.repo.ObtenerPorID(ctx, cmd.UsuarioID)
	if err != nil {
		return nil, dominio.ErrUsuarioNoEncontrado
	}

	if usuario.EstadoVerificacion == "VERIFICADO" {
		return nil, dominio.ErrCorreoYaVerificado
	}

	token, err := uc.idGenerator.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar token: %w", err)
	}

	expiraEn := time.Now().Add(uc.config.TokenExpiracion)
	prueba := dominio.NuevaPruebaVerificacion(token, expiraEn)

	if err := uc.repo.ActualizarPrueba(ctx, cmd.UsuarioID, prueba); err != nil {
		return nil, fmt.Errorf("error al persistir token: %w", err)
	}

	expiracionHoras := fmt.Sprintf("%.0f", uc.config.TokenExpiracion.Hours())
	go func() {
		if err := uc.emailServicio.EnviarTemplate(ctx, usuario.Correo,
			notificaciones.TipoVerificacionCorreo,
			map[string]string{
				"nombre":           usuario.Nombre,
				"token":            token,
				"expiracion_horas": expiracionHoras,
			},
		); err != nil {
			fmt.Printf("[VerificarCorreoCasoDeUso] Error al enviar email: %v\n", err)
		}
	}()

	return &RespuestaSolicitarVerificacion{
		Mensaje: "Email de verificación enviado",
	}, nil
}

func (uc *VerificarCorreoCasoDeUso) Confirmar(ctx context.Context, cmd ComandoConfirmarVerificacion) (*RespuestaConfirmarVerificacion, error) {
	if cmd.Token == "" {
		return nil, dominio.ErrEnlaceInvalido
	}

	hash := dominio.HashearToken(cmd.Token)

	usuario, err := uc.repo.ObtenerPorHashToken(ctx, hash)
	if err != nil {
		return nil, dominio.ErrEnlaceInvalido
	}

	if usuario.PruebaVerificacion.Expiro(time.Now()) {
		if err := uc.repo.ActualizarPrueba(ctx, usuario.ID, dominio.PruebaVerificacionVacia()); err != nil {
			fmt.Printf("[VerificarCorreoCasoDeUso] Error al limpiar prueba: %v\n", err)
		}
		return nil, dominio.ErrEnlaceExpirado
	}

	if err := uc.repo.MarcarVerificado(ctx, usuario.ID); err != nil {
		return nil, fmt.Errorf("error al marcar verificado: %w", err)
	}

	return &RespuestaConfirmarVerificacion{
		Mensaje: "Correo verificado exitosamente",
	}, nil
}

func (uc *VerificarCorreoCasoDeUso) Reenviar(ctx context.Context, cmd ComandoReenviarVerificacion) (*RespuestaSolicitarVerificacion, error) {
	usuario, err := uc.repo.ObtenerPorID(ctx, cmd.UsuarioID)
	if err != nil {
		return nil, dominio.ErrUsuarioNoEncontrado
	}

	if usuario.EstadoVerificacion == "VERIFICADO" {
		return nil, dominio.ErrCorreoYaVerificado
	}

	if usuario.ContadorReenvios >= uc.config.MaxReenvios {
		return nil, dominio.ErrDemasiadosReenvios
	}

	return uc.Solicitar(ctx, ComandoSolicitarVerificacion{
		UsuarioID: cmd.UsuarioID,
	})
}

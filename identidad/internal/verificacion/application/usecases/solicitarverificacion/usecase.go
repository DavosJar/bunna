package solicitarverificacion

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
	FrontendURL     string
}

type SolicitarVerificacionCasoDeUso struct {
	repo          dominio.VerificacionRepositorio
	emailServicio notificaciones.EmailServicio
	idGenerator   shareddomain.GeneradorID
	config        ConfigVerificacion
}

func NewSolicitarVerificacionCasoDeUso(
	repo dominio.VerificacionRepositorio,
	emailServicio notificaciones.EmailServicio,
	idGenerator shareddomain.GeneradorID,
	config ConfigVerificacion,
) *SolicitarVerificacionCasoDeUso {
	if config.TokenExpiracion == 0 {
		config.TokenExpiracion = 24 * time.Hour
	}
	if config.MaxReenvios == 0 {
		config.MaxReenvios = 5
	}
	if config.VentanaReenvios == 0 {
		config.VentanaReenvios = 24 * time.Hour
	}
	if config.FrontendURL == "" {
		config.FrontendURL = "http://localhost:5173"
	}
	return &SolicitarVerificacionCasoDeUso{
		repo:          repo,
		emailServicio: emailServicio,
		idGenerator:   idGenerator,
		config:        config,
	}
}

func (uc *SolicitarVerificacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoSolicitarVerificacion) (*RespuestaSolicitarVerificacion, error) {
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
	urlVerificacion := fmt.Sprintf("%s/verificar-correo?token=%s", uc.config.FrontendURL, token)

	go func() {
		if err := uc.emailServicio.EnviarTemplate(ctx, usuario.Correo,
			notificaciones.TipoVerificacionCorreo,
			map[string]string{
				"nombre":           usuario.Nombre,
				"url_verificacion": urlVerificacion,
				"expiracion_horas": expiracionHoras,
			},
		); err != nil {
			fmt.Printf("[SolicitarVerificacionCasoDeUso] Error al enviar email: %v\n", err)
		}
	}()

	return &RespuestaSolicitarVerificacion{
		Mensaje: "Email de verificación enviado",
	}, nil
}

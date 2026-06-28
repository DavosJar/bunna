package confirmarverificacion

import (
	"context"
	"fmt"
	"time"

	dominio "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
)

type ConfigVerificacion struct {
	TokenExpiracion time.Duration
	MaxReenvios     int
	VentanaReenvios time.Duration
	FrontendURL     string
}

type ConfirmarVerificacionCasoDeUso struct {
	repo   dominio.VerificacionRepositorio
	config ConfigVerificacion
}

func NewConfirmarVerificacionCasoDeUso(
	repo dominio.VerificacionRepositorio,
	config ConfigVerificacion,
) *ConfirmarVerificacionCasoDeUso {
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
	return &ConfirmarVerificacionCasoDeUso{
		repo:   repo,
		config: config,
	}
}

func (uc *ConfirmarVerificacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoConfirmarVerificacion) (*RespuestaConfirmarVerificacion, error) {
	if cmd.Token == "" {
		return nil, dominio.ErrEnlaceInvalido
	}

	hash := dominio.HashearToken(cmd.Token)

	usuario, err := uc.repo.ObtenerPorHashToken(ctx, hash)
	if err != nil {
		return nil, dominio.ErrEnlaceInvalido
	}

	// Si ya está verificado, es idempotente
	if usuario.EstadoVerificacion == "VERIFICADO" {
		return &RespuestaConfirmarVerificacion{
			Mensaje: "Correo ya verificado",
		}, nil
	}

	if usuario.PruebaVerificacion.Expiro(time.Now()) {
		if err := uc.repo.ActualizarPrueba(ctx, usuario.ID, dominio.PruebaVerificacionVacia()); err != nil {
			fmt.Printf("[ConfirmarVerificacionCasoDeUso] Error al limpiar prueba: %v\n", err)
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

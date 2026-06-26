package confirmarrecuperacion

import (
	"context"
	"fmt"
	"time"

	dominio "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	uc_validar "github.com/davosjar/bunna/services/identidad/internal/recuperacion/application/usecases/validartokenrecuperacion"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/davosjar/bunna/services/identidad/internal/shared/application"
)

type ConfirmarRecuperacionCasoDeUso struct {
	tokenRepo     dominio.TokenRecuperacionRepositorio
	usuarioRepo   dominio.UsuarioRecuperacionRepositorio
	sesionRepo    sesiones_domain.SesionRepositorio
	encriptacion  seguridad_domain.EncriptacionServicio
	validarTokenUC *uc_validar.ValidarTokenRecuperacionCasoDeUso
}

func NewConfirmarRecuperacionCasoDeUso(
	tokenRepo dominio.TokenRecuperacionRepositorio,
	usuarioRepo dominio.UsuarioRecuperacionRepositorio,
	sesionRepo sesiones_domain.SesionRepositorio,
	encriptacion seguridad_domain.EncriptacionServicio,
	validarTokenUC *uc_validar.ValidarTokenRecuperacionCasoDeUso,
) *ConfirmarRecuperacionCasoDeUso {
	return &ConfirmarRecuperacionCasoDeUso{
		tokenRepo:      tokenRepo,
		usuarioRepo:    usuarioRepo,
		sesionRepo:     sesionRepo,
		encriptacion:   encriptacion,
		validarTokenUC: validarTokenUC,
	}
}

func (uc *ConfirmarRecuperacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoConfirmarRecuperacion) (*RespuestaConfirmarRecuperacion, error) {
	if err := application.ValidarFormatoPassword(cmd.NuevaPassword, "nueva_password"); err != nil {
		return nil, dominio.ErrPasswordDebil
	}

	respuestaValidar, err := uc.validarTokenUC.Ejecutar(ctx, &uc_validar.ComandoValidarTokenRecuperacion{Token: cmd.Token})
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
		fmt.Printf("[ConfirmarRecuperacionCasoDeUso] Error al marcar token usado: %v\n", err)
	}

	if err := uc.sesionRepo.InvalidarTodasPorUsuarioID(ctx, respuestaValidar.UsuarioID); err != nil {
		fmt.Printf("[ConfirmarRecuperacionCasoDeUso] Error al invalidar sesiones: %v\n", err)
	}

	return &RespuestaConfirmarRecuperacion{
		Mensaje: "Contraseña actualizada exitosamente",
	}, nil
}

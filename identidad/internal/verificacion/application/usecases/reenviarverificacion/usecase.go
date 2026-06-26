package reenviarverificacion

import (
	"context"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	dominio "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
	uc_solicitar "github.com/davosjar/bunna/services/identidad/internal/verificacion/application/usecases/solicitarverificacion"
)

type ConfigVerificacion struct {
	TokenExpiracion time.Duration
	MaxReenvios     int
	VentanaReenvios time.Duration
	FrontendURL     string
}

type ReenviarVerificacionCasoDeUso struct {
	repo          dominio.VerificacionRepositorio
	emailServicio notificaciones.EmailServicio
	idGenerator   shareddomain.GeneradorID
	solicitarUC   *uc_solicitar.SolicitarVerificacionCasoDeUso
	config        ConfigVerificacion
}

func NewReenviarVerificacionCasoDeUso(
	repo dominio.VerificacionRepositorio,
	emailServicio notificaciones.EmailServicio,
	idGenerator shareddomain.GeneradorID,
	solicitarUC *uc_solicitar.SolicitarVerificacionCasoDeUso,
	config ConfigVerificacion,
) *ReenviarVerificacionCasoDeUso {
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
	return &ReenviarVerificacionCasoDeUso{
		repo:          repo,
		emailServicio: emailServicio,
		idGenerator:   idGenerator,
		solicitarUC:   solicitarUC,
		config:        config,
	}
}

func (uc *ReenviarVerificacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoReenviarVerificacion) (*RespuestaSolicitarVerificacion, error) {
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

	return uc.solicitarUC.Ejecutar(ctx, &uc_solicitar.ComandoSolicitarVerificacion{
		UsuarioID: cmd.UsuarioID,
	})
}

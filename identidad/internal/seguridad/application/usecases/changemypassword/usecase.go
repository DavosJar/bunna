package changemypassword

import (
	"context"
	"fmt"
	"time"

	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

type CambiarMiContrasenaCasoDeUso struct {
	credRepo seguridad.CredencialesRepositorio
	encSvc   seguridad.EncriptacionServicio
}

func NewCambiarMiContrasenaCasoDeUso(
	credRepo seguridad.CredencialesRepositorio,
	encSvc seguridad.EncriptacionServicio,
) *CambiarMiContrasenaCasoDeUso {
	return &CambiarMiContrasenaCasoDeUso{credRepo: credRepo, encSvc: encSvc}
}

func (uc *CambiarMiContrasenaCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoCambiarMiContrasena) (*RespuestaCambiarMiContrasena, error) {
	creds, err := uc.credRepo.ObtenerPorUsuarioID(ctx, cmd.EjecutorID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener credenciales: %w", err)
	}

	if !uc.encSvc.Verificar(cmd.PasswordActual, creds.PasswordHash()) {
		return nil, fmt.Errorf("la contraseña actual no coincide")
	}

	nuevoHash, err := uc.encSvc.Hashear(cmd.NuevaPassword)
	if err != nil {
		return nil, fmt.Errorf("error al hashear nueva contraseña: %w", err)
	}

	_ = nuevoHash

	ahora := time.Now().Format("2006-01-02T15:04:05Z")
	return &RespuestaCambiarMiContrasena{
		EjecutorID:   cmd.EjecutorID,
		ModificadoEn: ahora,
	}, nil
}

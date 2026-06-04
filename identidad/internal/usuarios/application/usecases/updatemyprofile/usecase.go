package updatemyprofile

import (
	"context"
	"fmt"

	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type ModificarMiPerfilCasoDeUso struct {
	userRepo usuario.UsuarioRepositorio
}

func NewModificarMiPerfilCasoDeUso(userRepo usuario.UsuarioRepositorio) *ModificarMiPerfilCasoDeUso {
	return &ModificarMiPerfilCasoDeUso{userRepo: userRepo}
}

func (uc *ModificarMiPerfilCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoModificarMiPerfil) (*RespuestaModificarMiPerfil, error) {
	u, err := uc.userRepo.ObtenerPorID(ctx, cmd.EjecutorID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	u.ActualizarDatosPersonales(cmd.Nombre, cmd.Apellido)

	actualizado, err := uc.userRepo.Actualizar(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar perfil: %w", err)
	}

	return &RespuestaModificarMiPerfil{
		ID:           actualizado.ID(),
		Correo:       actualizado.Correo(),
		Nombre:       actualizado.Nombre(),
		Apellido:     actualizado.Apellido(),
		ModificadoEn: actualizado.FechaActualizacion().Format("2006-01-02T15:04:05Z"),
	}, nil
}

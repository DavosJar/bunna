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

	// Por ahora solo lectura del perfil; los setters se agregarían al agregado
	_ = u

	return &RespuestaModificarMiPerfil{
		ID:           u.ID(),
		Correo:       u.Correo(),
		Nombre:       u.Nombre(),
		Apellido:     u.Apellido(),
		ModificadoEn: u.FechaActualizacion().Format("2006-01-02T15:04:05Z"),
	}, nil
}

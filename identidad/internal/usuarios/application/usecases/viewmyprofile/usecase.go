package viewmyprofile

import (
	"context"
	"fmt"

	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type VerMiPerfilCasoDeUso struct {
	userRepo usuario.UsuarioRepositorio
}

func NewVerMiPerfilCasoDeUso(userRepo usuario.UsuarioRepositorio) *VerMiPerfilCasoDeUso {
	return &VerMiPerfilCasoDeUso{userRepo: userRepo}
}

func (uc *VerMiPerfilCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoVerMiPerfil) (*RespuestaVerMiPerfil, error) {
	u, err := uc.userRepo.ObtenerPorID(ctx, cmd.EjecutorID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	return &RespuestaVerMiPerfil{
		ID:       u.ID(),
		Correo:   u.Correo(),
		Nombre:   u.Nombre(),
		Apellido: u.Apellido(),
		Telefono: u.Telefono(),
		Estado:   string(u.Estado()),
		CreadoEn: u.FechaCreacion().Format("2006-01-02T15:04:05Z"),
	}, nil
}

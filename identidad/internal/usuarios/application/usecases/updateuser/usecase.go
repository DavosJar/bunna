package updateuser

import (
	"context"
	"fmt"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type ModificarUsuarioCasoDeUso struct {
	userRepo usuario.UsuarioRepositorio
	authSvc  rbac.AuthorizationService
}

func NewModificarUsuarioCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	authSvc rbac.AuthorizationService,
) *ModificarUsuarioCasoDeUso {
	return &ModificarUsuarioCasoDeUso{userRepo: userRepo, authSvc: authSvc}
}

// DISABLED — Código inconexo, no eliminar.
// La modificación de datos de otros usuarios administrativamente está
// deshabilitada. Cada usuario modifica su propio perfil a través de
// ModificarMiPerfilCasoDeUso (PUT /api/v1/mi-perfil). Este caso de uso
// se conserva por si el negocio requiere reactivarlo.
func (uc *ModificarUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoModificarUsuario) (*RespuestaModificarUsuario, error) {
	return nil, fmt.Errorf("modificación administrativa de usuarios deshabilitada: cada usuario modifica su propio perfil")
}

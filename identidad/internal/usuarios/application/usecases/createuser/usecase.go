package createuser

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type CrearUsuarioCasoDeUso struct {
	userRepo usuario.UsuarioRepositorio
	credRepo seguridad.CredencialesRepositorio
	encSvc   seguridad.EncriptacionServicio
	authSvc  rbac.AuthorizationService
	idGen    shareddomain.GeneradorID
}

func NewCrearUsuarioCasoDeUso(
	userRepo usuario.UsuarioRepositorio,
	credRepo seguridad.CredencialesRepositorio,
	encSvc seguridad.EncriptacionServicio,
	authSvc rbac.AuthorizationService,
	idGen shareddomain.GeneradorID,
) *CrearUsuarioCasoDeUso {
	return &CrearUsuarioCasoDeUso{
		userRepo: userRepo,
		credRepo: credRepo,
		encSvc:   encSvc,
		authSvc:  authSvc,
		idGen:    idGen,
	}
}

// DISABLED — Código inconexo, no eliminar.
// Los usuarios no se crean administrativamente: se registran por sí mismos
// o son invitados a un tenant. Este caso de uso se deshabilita pero se
// conserva porque forma parte del flujo de creación directa que podría
// reactivarse si el negocio cambia. No es código muerto.
func (uc *CrearUsuarioCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoCrearUsuario) (*RespuestaCrearUsuario, error) {
	return nil, fmt.Errorf("creación administrativa de usuarios deshabilitada: los usuarios se registran o son invitados")
}

func emailTieneDominioValido(correo string) error {
	parts := strings.SplitN(correo, "@", 2)
	if len(parts) != 2 {
		return fmt.Errorf("formato de correo inválido")
	}
	dominio := parts[1]
	mxRecords, err := net.LookupMX(dominio)
	if err != nil || len(mxRecords) == 0 {
		return fmt.Errorf("el dominio del correo no existe o no acepta correos")
	}
	return nil
}

package usuario

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type UsuarioRepositorio interface {
	Crear(ctx context.Context, usuario *Usuario) (*Usuario, error)
	Actualizar(ctx context.Context, usuario *Usuario) (*Usuario, error)
	Eliminar(ctx context.Context, id string) error
	ObtenerPorID(ctx context.Context, id string) (*Usuario, error)
	ObtenerPorCorreo(ctx context.Context, correo string) (*Usuario, error)
	Listar(ctx context.Context, especificacion EspecificacionUsuario, paginacion domain.Paginacion) ([]*Usuario, error)
}

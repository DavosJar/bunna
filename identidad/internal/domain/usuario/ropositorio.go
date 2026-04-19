package usuario

import "context"

type UsuarioRepositorio interface {
	Crear(ctx context.Context, usuario *Usuario) (*Usuario, error)
	Actualizar(ctx context.Context, usuario *Usuario) (*Usuario, error)
	Eliminar(ctx context.Context, id string) error
	ObtenerPorID(ctx context.Context, id string) (*Usuario, error)
	Listar(ctx context.Context, especificacion EspecificacionUsuario, paginacion Paginacion) ([]*Usuario, error)
}

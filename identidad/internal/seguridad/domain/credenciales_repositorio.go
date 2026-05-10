package domain

import "context"

type CredencialesRepositorio interface {
	Crear(ctx context.Context, credenciales *CredencialesUsuario) (*CredencialesUsuario, error)
	Actualizar(ctx context.Context, credenciales *CredencialesUsuario) (*CredencialesUsuario, error)
	ObtenerPorUsuarioID(ctx context.Context, usuarioID string) (*CredencialesUsuario, error)
	Eliminar(ctx context.Context, usuarioID string) error
	Find(ctx context.Context, especificacion EspecificacionCredenciales, paginacion Paginacion) ([]*CredencialesUsuario, error)
}

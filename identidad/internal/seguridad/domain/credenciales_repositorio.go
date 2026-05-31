package domain

import (
	"context"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type CredencialesRepositorio interface {
	Crear(ctx context.Context, credenciales *CredencialesUsuario) (*CredencialesUsuario, error)
	Actualizar(ctx context.Context, credenciales *CredencialesUsuario) (*CredencialesUsuario, error)
	ObtenerPorUsuarioID(ctx context.Context, usuarioID string) (*CredencialesUsuario, error)
	Eliminar(ctx context.Context, usuarioID string) error
	Find(ctx context.Context, especificacion EspecificacionCredenciales, paginacion shareddomain.Paginacion) ([]*CredencialesUsuario, error)
}

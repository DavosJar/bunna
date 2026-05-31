package domain

import (
	"context"
	"time"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type SesionRepositorio interface {
	Crear(ctx context.Context, sesion *Sesion) (*Sesion, error)
	Actualizar(ctx context.Context, sesion *Sesion) (*Sesion, error)
	ObtenerPorID(ctx context.Context, sesionID string) (*Sesion, error)
	ObtenerPorRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*Sesion, error)
	ListarActivasPorUsuarioID(ctx context.Context, usuarioID string, ahora time.Time) ([]*Sesion, error)
	Listar(ctx context.Context, especificacion EspecificacionSesion, paginacion shareddomain.Paginacion) ([]*Sesion, error)
	InvalidarTodasPorUsuarioID(ctx context.Context, usuarioID string) error
	Eliminar(ctx context.Context, sesionID string) error
}

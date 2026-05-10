package domain

import (
	"context"
	"time"
)

type SesionRepositorio interface {
	Crear(ctx context.Context, sesion *Sesion) (*Sesion, error)
	Actualizar(ctx context.Context, sesion *Sesion) (*Sesion, error)
	ObtenerPorID(ctx context.Context, sesionID string) (*Sesion, error)
	ObtenerPorRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*Sesion, error)
	ListarActivasPorUsuarioID(ctx context.Context, usuarioID string, ahora time.Time) ([]*Sesion, error)
	InvalidarTodasPorUsuarioID(ctx context.Context, usuarioID string) error
	Eliminar(ctx context.Context, sesionID string) error
}
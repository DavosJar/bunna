package domain

import (
	"context"

	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

type FincaRepositorio interface {
	Crear(ctx context.Context, finca *Finca) error
	ObtenerPorID(ctx context.Context, id string) (*Finca, error)
	ListarPorUsuario(ctx context.Context, usuarioID string) ([]Finca, error)
	ListarTodas(ctx context.Context) ([]Finca, error)
	Listar(ctx context.Context, especificacion EspecificacionFinca, paginacion shared.Paginacion) ([]Finca, error)
	Actualizar(ctx context.Context, finca *Finca) error
	Eliminar(ctx context.Context, id string) error
	ContarLotes(ctx context.Context, fincaID string) (int, error)
}

type LoteRepositorio interface {
	Crear(ctx context.Context, lote *Lote) error
	ObtenerPorID(ctx context.Context, id string) (*Lote, error)
	ListarPorFinca(ctx context.Context, fincaID string) ([]Lote, error)
	Listar(ctx context.Context, especificacion EspecificacionLote, paginacion shared.Paginacion) ([]Lote, error)
	Actualizar(ctx context.Context, lote *Lote) error
	Eliminar(ctx context.Context, id string) error
}

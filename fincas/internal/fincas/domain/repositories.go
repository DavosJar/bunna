package domain

import "context"

type FincaRepositorio interface {
	Crear(ctx context.Context, finca *Finca) error
	ObtenerPorID(ctx context.Context, id string) (*Finca, error)
	ListarPorUsuario(ctx context.Context, usuarioID string) ([]Finca, error)
	ListarTodas(ctx context.Context) ([]Finca, error)
	Actualizar(ctx context.Context, finca *Finca) error
	Eliminar(ctx context.Context, id string) error
}

type LoteRepositorio interface {
	Crear(ctx context.Context, lote *Lote) error
	ObtenerPorID(ctx context.Context, id string) (*Lote, error)
	ListarPorFinca(ctx context.Context, fincaID string) ([]Lote, error)
	ContarPorFinca(ctx context.Context, fincaID string) (int, error)
	Actualizar(ctx context.Context, lote *Lote) error
	Eliminar(ctx context.Context, id string) error
	EliminarPorFinca(ctx context.Context, fincaID string) error
}

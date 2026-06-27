package domain

import (
	"context"

	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

type NodoRepositorio interface {
	Crear(ctx context.Context, nodo *Nodo) error
	ObtenerPorID(ctx context.Context, id string) (*Nodo, error)
	ObtenerPorNodeKey(ctx context.Context, nodeKey string) (*Nodo, error)
	Listar(ctx context.Context, especificacion EspecificacionNodo, paginacion shared.Paginacion) ([]Nodo, error)
	Actualizar(ctx context.Context, nodo *Nodo) error
	Eliminar(ctx context.Context, id string) error
}

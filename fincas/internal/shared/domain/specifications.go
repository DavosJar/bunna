package domain

import "context"

type CriterioFiltro struct {
	Campo    string
	Operador string
	Valor    any
}

type Paginacion struct {
	Pagina       int
	TamanoPagina int
	Ordenaciones []Ordenacion
}

type Ordenacion struct {
	Campo string
	Tipo  TipoOrdenacion
}

type TipoOrdenacion string

const (
	ASC  TipoOrdenacion = "ASC"
	DESC TipoOrdenacion = "DESC"
)

type GeneradorID interface {
	NextID(ctx context.Context) (string, error)
}

package domain

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

type TipoOrdenacion string

const (
	ASC  TipoOrdenacion = "ASC"
	DESC TipoOrdenacion = "DESC"
)

type Ordenacion struct {
	Campo string
	Tipo  TipoOrdenacion
}

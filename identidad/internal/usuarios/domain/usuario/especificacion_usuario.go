package usuario

type EspecificacionUsuario struct {
	ListaLiltros []CriterioFiltro
}

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

var ColumnasPermitidas = map[string]bool{
	"nombre":                   true,
	"apellido":                 true,
	"correo":                   true,
	"fechaCreacion":            true,
	"fechaActualizacion":       true,
	"estado":                   true,
	"telefono":                 true,
	"estadoVerificacionCorreo": true,
}

//en la capa de infraestructura se validará que los campos de filtro y ordenación
//con un if !ColumnasPermitidas[campo] para evitar inyecciones SQL

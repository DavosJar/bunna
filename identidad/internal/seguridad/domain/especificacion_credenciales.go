package domain

type EspecificacionCredenciales struct {
	ListaFiltros []CriterioFiltro
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

type Ordenacion struct {
	Campo string
	Tipo  TipoOrdenacion
}

type TipoOrdenacion string

const (
	ASC  TipoOrdenacion = "ASC"
	DESC TipoOrdenacion = "DESC"
)

var ColumnasPermitidas = map[string]bool{
	"usuarioID":        true,
	"activo":           true,
	"intentosFallidos": true,
	"bloqueadoHasta":   true,
	"correoVerificado": true,
}

// Nota: BETWEEN será implementado próximamente para permitir rangos en bloqueadoHasta
// Ejemplo de uso futuro:
// case "BETWEEN":
//     // Validar que filtro.Valor es un slice con dos elementos [min, max]
//     query = query.Where(columnaDB+" BETWEEN ? AND ?", filtro.Valor.([]interface{})[0], filtro.Valor.([]interface{})[1])

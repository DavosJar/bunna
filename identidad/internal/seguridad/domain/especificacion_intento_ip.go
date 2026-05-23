package domain

import shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type EspecificacionIntentoIP struct {
	ListaFiltros []shareddomain.CriterioFiltro
}

var ColumnasPermitidasIntentoIP = map[string]bool{
	"ip":             true,
	"contador":       true,
	"ventanaInicio":  true,
	"bloqueadoHasta": true,
}

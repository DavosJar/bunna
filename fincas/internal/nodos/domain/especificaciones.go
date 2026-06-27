package domain

import shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"

type EspecificacionNodo struct {
	Filtros []shared.CriterioFiltro
}

var ColumnasPermitidasNodos = map[string]bool{
	"nombre":   true,
	"nodeKey":  true,
	"fincaID":  true,
	"loteID":   true,
	"tenantID": true,
	"estado":   true,
}

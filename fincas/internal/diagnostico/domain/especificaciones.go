package domain

import shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"

type EspecificacionDiagnostico struct {
	Filtros []shared.CriterioFiltro
}

var ColumnasPermitidasDiagnostico = map[string]bool{
	"nombre":    true,
	"muestraID": true,
	"tenantID":  true,
}

type EspecificacionMuestra struct {
	Filtros []shared.CriterioFiltro
}

var ColumnasPermitidasMuestra = map[string]bool{
	"nombre":   true,
	"loteID":   true,
	"tenantID": true,
}

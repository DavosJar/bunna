package domain

import shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"

type EspecificacionFinca struct {
	Filtros []shared.CriterioFiltro
}

type EspecificacionLote struct {
	Filtros []shared.CriterioFiltro
}

var ColumnasPermitidasFincas = map[string]bool{
	"nombre":    true,
	"ubicacion": true,
	"estado":    true,
	"usuarioID": true,
	"tenantID":  true,
}

var ColumnasPermitidasLotes = map[string]bool{
	"nombre":  true,
	"area":    true,
	"estado":  true,
	"fincaID": true,
}

package domain

import shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type EspecificacionSesion struct {
	ListaFiltros []shareddomain.CriterioFiltro
}

var ColumnasPermitidas = map[string]bool{
	"usuarioID":              true,
	"estado":                 true,
	"ipOrigen":               true,
	"fechaCreacion":          true,
	"fechaActualizacion":     true,
	"fechaExpiracionAccess":  true,
	"fechaExpiracionRefresh": true,
	"ultimaActividad":        true,
	"contadorRefrescos":      true,
}

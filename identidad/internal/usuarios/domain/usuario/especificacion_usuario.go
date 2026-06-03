package usuario

import "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type EspecificacionUsuario struct {
	ListaLiltros []domain.CriterioFiltro
	TenantID     string
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

// en la capa de infraestructura se validará que los campos de filtro y ordenación
// con un if !ColumnasPermitidas[campo] para evitar inyecciones SQL

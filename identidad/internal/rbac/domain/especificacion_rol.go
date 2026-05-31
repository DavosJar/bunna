package rbac

import shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type EspecificacionRol struct {
	ListaFiltros []shareddomain.CriterioFiltro
}

var ColumnasPermitidasRol = map[string]bool{
	"nombre":    true,
	"esSistema": true,
}

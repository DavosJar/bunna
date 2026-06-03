package rbac

import shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type EspecificacionRol struct {
	ListaFiltros []shareddomain.CriterioFiltro
	TenantID     string // filtro opcional por tenant
}

var ColumnasPermitidasRol = map[string]bool{
	"nombre":    true,
	"esSistema": true,
	"tenant_id": true,
}

package listusers

import "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type ComandoListarUsuarios struct {
	Filtros    []domain.CriterioFiltro
	Paginacion domain.Paginacion
	TenantID   string
	EjecutorID string
}

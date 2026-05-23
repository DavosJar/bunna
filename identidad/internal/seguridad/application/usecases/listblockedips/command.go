package listblockedips

import "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type ComandoListarIPsBloqueadas struct {
	Paginacion domain.Paginacion
	TenantID   string
	EjecutorID string
}

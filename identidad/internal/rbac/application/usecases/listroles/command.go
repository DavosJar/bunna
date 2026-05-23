package listroles

import "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type ComandoListarRoles struct {
	Paginacion domain.Paginacion
	TenantID   string
	EjecutorID string
}

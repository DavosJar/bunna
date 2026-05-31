package listsessions

import "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type ComandoListarSesiones struct {
	UsuarioID  string
	Paginacion domain.Paginacion
	TenantID   string
	EjecutorID string
}

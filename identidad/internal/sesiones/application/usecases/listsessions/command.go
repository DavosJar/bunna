package listsessions

import "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type ComandoListarSesiones struct {
	UsuarioID  string
	Paginacion domain.Paginacion
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoListarSesiones) ToLog() map[string]any {
	return map[string]any{
		"usuario_id":  c.UsuarioID,
		"paginacion": map[string]any{
			"pagina":        c.Paginacion.Pagina,
			"tamano_pagina": c.Paginacion.TamanoPagina,
		},
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

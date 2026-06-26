package listroles

import "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type ComandoListarRoles struct {
	Paginacion domain.Paginacion
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoListarRoles) ToLog() map[string]any {
	return map[string]any{
		"paginacion": map[string]any{
			"pagina":        c.Paginacion.Pagina,
			"tamano_pagina": c.Paginacion.TamanoPagina,
		},
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

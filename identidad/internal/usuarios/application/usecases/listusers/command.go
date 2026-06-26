package listusers

import "github.com/davosjar/bunna/services/identidad/internal/shared/domain"

type ComandoListarUsuarios struct {
	Filtros    []domain.CriterioFiltro
	Paginacion domain.Paginacion
	TenantID   string
	EjecutorID string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoListarUsuarios) ToLog() map[string]any {
	return map[string]any{
		"filtros_count": len(c.Filtros),
		"paginacion": map[string]any{
			"pagina":        c.Paginacion.Pagina,
			"tamano_pagina": c.Paginacion.TamanoPagina,
		},
		"tenant_id":   c.TenantID,
		"ejecutor_id": c.EjecutorID,
	}
}

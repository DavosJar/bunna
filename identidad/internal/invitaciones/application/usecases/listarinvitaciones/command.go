package listarinvitaciones

type ComandoListarInvitaciones struct {
	TenantID     string
	Pagina       int
	TamanoPagina int
	Estado       string
}

func (c ComandoListarInvitaciones) ToLog() map[string]any {
	return map[string]any{
		"tenant_id":      c.TenantID,
		"pagina":         c.Pagina,
		"tamano_pagina":  c.TamanoPagina,
		"estado":         c.Estado,
	}
}

type RespuestaListarInvitaciones struct {
	Invitaciones []InvitacionDTO
	Total        int
}

type InvitacionDTO struct {
	ID            string
	Email         string
	Nombre        string
	RolID         string
	RolNombre     string
	Estado        string
	FechaCreacion string
	Expiracion    string
}

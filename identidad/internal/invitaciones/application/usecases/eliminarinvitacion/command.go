package eliminarinvitacion

type ComandoEliminarInvitacion struct {
	InvitacionID string
	TenantID     string
	EjecutorID   string
}

func (c ComandoEliminarInvitacion) ToLog() map[string]any {
	return map[string]any{
		"invitacion_id": c.InvitacionID,
		"tenant_id":     c.TenantID,
		"ejecutor_id":   c.EjecutorID,
	}
}

type RespuestaEliminarInvitacion struct {
	Mensaje string
}

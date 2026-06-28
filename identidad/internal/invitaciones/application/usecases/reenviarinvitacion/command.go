package reenviarinvitacion

type ComandoReenviarInvitacion struct {
	InvitacionID string
	TenantID     string
}

func (c ComandoReenviarInvitacion) ToLog() map[string]any {
	return map[string]any{
		"invitacion_id": c.InvitacionID,
		"tenant_id":     c.TenantID,
	}
}

type RespuestaReenviarInvitacion struct {
	Mensaje string
}

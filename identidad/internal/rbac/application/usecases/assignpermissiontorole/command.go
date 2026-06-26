package assignpermissiontorole

type ComandoAsignarPermisoARol struct {
	RolID         string
	PermisoCodigo string
	TenantID      string
	EjecutorID    string
	AsignadoPor   string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoAsignarPermisoARol) ToLog() map[string]any {
	return map[string]any{
		"rol_id":          c.RolID,
		"permiso_codigo":  c.PermisoCodigo,
		"tenant_id":       c.TenantID,
		"ejecutor_id":     c.EjecutorID,
		"asignado_por":    c.AsignadoPor,
	}
}

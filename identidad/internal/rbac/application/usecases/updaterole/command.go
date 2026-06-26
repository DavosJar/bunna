package updaterole

type ComandoModificarRol struct {
	RolID       string
	Nombre      string
	Descripcion string
	TenantID    string
	EjecutorID  string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoModificarRol) ToLog() map[string]any {
	return map[string]any{
		"rol_id":       c.RolID,
		"nombre":       c.Nombre,
		"descripcion":  c.Descripcion,
		"tenant_id":    c.TenantID,
		"ejecutor_id":  c.EjecutorID,
	}
}

package createrole

type ComandoCrearRol struct {
	Nombre      string
	Descripcion string
	Permisos    []string
	TenantID    string
	EjecutorID  string
}

// ToLog returns a safe representation — no sensitive fields.
func (c ComandoCrearRol) ToLog() map[string]any {
	return map[string]any{
		"nombre":       c.Nombre,
		"descripcion":  c.Descripcion,
		"permisos":     c.Permisos,
		"tenant_id":    c.TenantID,
		"ejecutor_id":  c.EjecutorID,
	}
}

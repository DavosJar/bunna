package switchtenant

// ComandoCambiarTenant contiene los datos necesarios para cambiar de tenant activo.
type ComandoCambiarTenant struct {
	UsuarioID string
	SesionID  string
	TenantID  string
}

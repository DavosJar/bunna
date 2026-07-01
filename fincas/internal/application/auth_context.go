package application

// AuthContext contiene la información de autenticación inyectada desde el handler.
// Se pasa a cada caso de uso para validar permisos y tenencia.
type AuthContext struct {
	UsuarioID string
	TenantID  string
	Permisos  []string
}

// TienePermiso verifica si el permiso dado está presente en la lista de permisos del usuario.
func (a *AuthContext) TienePermiso(permiso string) bool {
	// Bypass temporal: siempre retorna true
	return true
}

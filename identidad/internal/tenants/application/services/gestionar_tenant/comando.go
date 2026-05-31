package gestionar_tenant

// ComandoCrearTenant representa los datos para crear un tenant
type ComandoCrearTenant struct {
	SolicitanteID string
	Nombre        string
	Slug          string
}

// ComandoActivarTenant representa los datos para activar un tenant
type ComandoActivarTenant struct {
	SolicitanteID string
	TenantID      string
}

// ComandoDesactivarTenant representa los datos para desactivar un tenant
type ComandoDesactivarTenant struct {
	SolicitanteID string
	TenantID      string
}

// ComandoAgregarUsuario representa los datos para agregar un usuario a un tenant
type ComandoAgregarUsuario struct {
	SolicitanteID string
	UsuarioID     string
	TenantID      string
}

// ComandoRemoverUsuario representa los datos para remover un usuario de un tenant
type ComandoRemoverUsuario struct {
	SolicitanteID string
	UsuarioID     string
	TenantID      string
}

package rbac

// TenantClaims contiene los datos de autorización de un usuario en un tenant
type TenantClaims struct {
	Slug     string   `json:"slug"`
	Roles    []string `json:"roles"`
	Permisos []string `json:"permisos"`
}

// UsuarioClaims contiene todos los datos de autorización del usuario para el JWT
type UsuarioClaims struct {
	UsuarioID string                  `json:"sub"`
	Global    bool                    `json:"global"`
	Tenants   map[string]TenantClaims `json:"tenants"`
}
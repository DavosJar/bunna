package dto

// TenantConRolDTO representa un tenant con el rol del usuario en él.
type TenantConRolDTO struct {
	ID       string `json:"id"        doc:"ID del tenant"`
	Nombre   string `json:"nombre"    doc:"Nombre del tenant"`
	Slug     string `json:"slug"      doc:"Slug del tenant"`
	Rol      string `json:"rol"       doc:"Rol del usuario en este tenant"`
	EsPropio bool   `json:"es_propio" doc:"Indica si es el tenant propio del usuario"`
}

// ListarMisTenantsResponse es la respuesta del endpoint mis-tenants.
type ListarMisTenantsResponse struct {
	Tenants  []TenantConRolDTO `json:"tenants"   doc:"Lista de tenants del usuario"`
	PropioID string            `json:"propio_id"  doc:"ID del tenant propio del usuario"`
}

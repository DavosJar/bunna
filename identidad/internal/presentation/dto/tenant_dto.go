package dto

type ConfigurarTenantRequest struct {
	Nombre string `json:"nombre" doc:"Nuevo nombre del tenant" example:"Mi Empresa"`
	Slug   string `json:"slug"   doc:"Nuevo slug del tenant"   example:"mi-empresa"`
}

type ConfigurarTenantResponse struct {
	TenantID     string `json:"tenant_id"     doc:"ID del tenant"`
	Nombre       string `json:"nombre"        doc:"Nombre actualizado"`
	Slug         string `json:"slug"          doc:"Slug actualizado"`
	ModificadoEn string `json:"modificado_en" doc:"Fecha de modificación"`
}

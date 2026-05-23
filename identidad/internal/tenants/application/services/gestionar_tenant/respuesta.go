package gestionar_tenant

import "time"

// DtoTenant representa la respuesta con datos de un tenant
type DtoTenant struct {
	ID             string    `json:"id"`
	Nombre         string    `json:"nombre"`
	Slug           string    `json:"slug"`
	Activo         bool      `json:"activo"`
	FechaCreacion  time.Time `json:"fecha_creacion"`
}
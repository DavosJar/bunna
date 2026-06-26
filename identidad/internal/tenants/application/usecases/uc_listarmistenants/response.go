package uc_listarmistenants

import "time"

type DtoTenant struct {
	ID            string    `json:"id"`
	Nombre        string    `json:"nombre"`
	Slug          string    `json:"slug"`
	Activo        bool      `json:"activo"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

type RespuestaListarMisTenants struct {
	Tenants []DtoTenant `json:"tenants"`
}

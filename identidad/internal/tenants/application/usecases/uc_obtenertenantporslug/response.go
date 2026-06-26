package uc_obtenertenantporslug

import "time"

type RespuestaObtenerTenantPorSlug struct {
	ID            string    `json:"id"`
	Nombre        string    `json:"nombre"`
	Slug          string    `json:"slug"`
	Activo        bool      `json:"activo"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

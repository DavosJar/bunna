package uc_obtenertenantporid

import "time"

type RespuestaObtenerTenantPorID struct {
	ID            string    `json:"id"`
	Nombre        string    `json:"nombre"`
	Slug          string    `json:"slug"`
	Activo        bool      `json:"activo"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

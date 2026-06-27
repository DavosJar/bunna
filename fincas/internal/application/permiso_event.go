package application

import "time"

// CatalogoPermisosPublicado se publica en startup con todos los permisos del módulo.
type CatalogoPermisosPublicado struct {
	EventID   string        `json:"event_id"`
	Tipo      string        `json:"tipo"`
	Origen    string        `json:"origen"`
	Modulo    string        `json:"modulo"`
	Permisos  []PermisoInfo `json:"permisos"`
	Version   string        `json:"version"`
	OcurredAt time.Time     `json:"ocurred_at"`
}

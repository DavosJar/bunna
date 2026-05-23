package tenant

import "time"

// Membresia representa la relación usuario ↔ tenant
type Membresia struct {
	usuarioID     string
	tenantID      string
	fechaCreacion time.Time
}

// NuevaMembresia crea una membresía entre un usuario y un tenant.
func NuevaMembresia(usuarioID, tenantID string) (*Membresia, error) {
	if usuarioID == "" {
		return nil, ErrUsuarioNoEsMiembro
	}
	if tenantID == "" {
		return nil, ErrTenantNoEncontrado
	}
	return &Membresia{
		usuarioID:     usuarioID,
		tenantID:      tenantID,
		fechaCreacion: time.Now(),
	}, nil
}

// NuevaMembresiaDesdeBD reconstruye una membresía desde persistencia.
func NuevaMembresiaDesdeBD(usuarioID, tenantID string, fechaCreacion time.Time) *Membresia {
	return &Membresia{
		usuarioID:     usuarioID,
		tenantID:      tenantID,
		fechaCreacion: fechaCreacion,
	}
}

// Getters
func (m *Membresia) UsuarioID() string        { return m.usuarioID }
func (m *Membresia) TenantID() string         { return m.tenantID }
func (m *Membresia) FechaCreacion() time.Time { return m.fechaCreacion }
package domain

import "time"

// EstadoFinca representa los estados válidos de una finca
type EstadoFinca string

const (
	FincaActiva            EstadoFinca = "ACTIVA"
	FincaPendienteEliminar EstadoFinca = "PENDIENTE_ELIMINACION"
)

// transicionesFinca define las transiciones válidas de estado
var transicionesFinca = map[EstadoFinca]map[EstadoFinca]bool{
	FincaActiva: {
		FincaPendienteEliminar: true,
	},
	FincaPendienteEliminar: {},
}

// Finca es el aggregate raíz del dominio de fincas
type Finca struct {
	id          string
	nombre      string
	ubicacion   string
	descripcion string
	usuarioID   string
	tenantID    *string
	estado      EstadoFinca
	createdAt   time.Time
	updatedAt   time.Time
}

// NuevaFinca crea una nueva finca. Sin validaciones de formato.
func NuevaFinca(nombre, ubicacion, descripcion, usuarioID string) *Finca {
	return &Finca{
		nombre:      nombre,
		ubicacion:   ubicacion,
		descripcion: descripcion,
		usuarioID:   usuarioID,
		estado:      FincaActiva,
	}
}

// Getters
func (f *Finca) ID() string              { return f.id }
func (f *Finca) Nombre() string          { return f.nombre }
func (f *Finca) Ubicacion() string       { return f.ubicacion }
func (f *Finca) Descripcion() string     { return f.descripcion }
func (f *Finca) UsuarioID() string       { return f.usuarioID }
func (f *Finca) TenantID() *string       { return f.tenantID }
func (f *Finca) Estado() EstadoFinca     { return f.estado }
func (f *Finca) CreatedAt() time.Time    { return f.createdAt }
func (f *Finca) UpdatedAt() time.Time    { return f.updatedAt }

// EsPropietario verifica si un usuario es el dueño de la finca
func (f *Finca) EsPropietario(usuarioID string) bool {
	return f.usuarioID == usuarioID
}

// TieneLotes verifica si una cantidad dada indica que la finca tiene lotes
func (f *Finca) TieneLotes(cantidad int) bool {
	return cantidad > 0
}

// Actualizar actualiza los datos de la finca
func (f *Finca) Actualizar(nombre, ubicacion, descripcion string) {
	f.nombre = nombre
	f.ubicacion = ubicacion
	f.descripcion = descripcion
}

// CambiarEstado cambia el estado de la finca validando la transición
func (f *Finca) CambiarEstado(siguiente EstadoFinca) error {
	if !transicionesFinca[f.estado][siguiente] {
		return ErrTransicionEstadoNoPermitida
	}
	f.estado = siguiente
	return nil
}

// NewFincaFromPersistence reconstruye una finca desde persistencia
func NewFincaFromPersistence(
	id, nombre, ubicacion, descripcion, usuarioID string,
	tenantID *string, estado EstadoFinca, createdAt, updatedAt time.Time,
) *Finca {
	return &Finca{
		id:          id,
		nombre:      nombre,
		ubicacion:   ubicacion,
		descripcion: descripcion,
		usuarioID:   usuarioID,
		tenantID:    tenantID,
		estado:      estado,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

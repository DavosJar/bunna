package domain

import "time"

// EstadoLote representa los estados válidos de un lote
type EstadoLote string

const (
	LoteActivo   EstadoLote = "ACTIVO"
	LoteEliminado EstadoLote = "ELIMINADO"
)

var transicionesLote = map[EstadoLote]map[EstadoLote]bool{
	LoteActivo: {
		LoteEliminado: true,
	},
	LoteEliminado: {},
}

// Lote es una subdivisión de una Finca
type Lote struct {
	id          string
	fincaID     string
	nombre      string
	area        float64
	descripcion string
	estado      EstadoLote
	createdAt   time.Time
	updatedAt   time.Time
}

// NuevoLote crea un nuevo lote asociado a una finca. Sin validaciones de formato.
func NuevoLote(fincaID, nombre string, area float64, descripcion string) *Lote {
	return &Lote{
		fincaID:     fincaID,
		nombre:      nombre,
		area:        area,
		descripcion: descripcion,
		estado:      LoteActivo,
	}
}

// Getters
func (l *Lote) ID() string            { return l.id }
func (l *Lote) FincaID() string       { return l.fincaID }
func (l *Lote) Nombre() string        { return l.nombre }
func (l *Lote) Area() float64         { return l.area }
func (l *Lote) Descripcion() string   { return l.descripcion }
func (l *Lote) Estado() EstadoLote    { return l.estado }
func (l *Lote) CreatedAt() time.Time  { return l.createdAt }
func (l *Lote) UpdatedAt() time.Time  { return l.updatedAt }

// Actualizar actualiza los datos del lote
func (l *Lote) Actualizar(nombre string, area float64, descripcion string) {
	l.nombre = nombre
	l.area = area
	l.descripcion = descripcion
}

// CambiarEstado cambia el estado del lote validando la transición
func (l *Lote) CambiarEstado(siguiente EstadoLote) error {
	if !transicionesLote[l.estado][siguiente] {
		return ErrTransicionEstadoNoPermitida
	}
	l.estado = siguiente
	return nil
}

// NewLoteFromPersistence reconstruye un lote desde persistencia
func NewLoteFromPersistence(
	id, fincaID, nombre string, area float64, descripcion string,
	estado EstadoLote, createdAt, updatedAt time.Time,
) *Lote {
	return &Lote{
		id:          id,
		fincaID:     fincaID,
		nombre:      nombre,
		area:        area,
		descripcion: descripcion,
		estado:      estado,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

package domain

import "time"

type EstadoNodo string

const (
	NodoActivo         EstadoNodo = "ACTIVO"
	NodoInactivo       EstadoNodo = "INACTIVO"
	NodoMantenimiento  EstadoNodo = "MANTENIMIENTO"
)

var transicionesNodo = map[EstadoNodo]map[EstadoNodo]bool{
	NodoActivo: {
		NodoInactivo:      true,
		NodoMantenimiento: true,
	},
	NodoInactivo: {
		NodoActivo:       true,
		NodoMantenimiento: true,
	},
	NodoMantenimiento: {
		NodoActivo:   true,
		NodoInactivo: true,
	},
}

type Nodo struct {
	id        string
	fincaID   string
	loteID    *string
	tenantID  string
	nombre    string
	nodeKey   string
	estado    EstadoNodo
	creadoEn  time.Time
	actualizadoEn time.Time
}

func NuevoNodo(tenantID, fincaID, nodeKey string, loteID *string, nombre string) *Nodo {
	return &Nodo{
		tenantID:  tenantID,
		fincaID:   fincaID,
		nodeKey:   nodeKey,
		loteID:    loteID,
		nombre:    nombre,
		estado:    NodoActivo,
		creadoEn:  time.Now(),
		actualizadoEn: time.Now(),
	}
}

func (n *Nodo) ID() string              { return n.id }
func (n *Nodo) FincaID() string         { return n.fincaID }
func (n *Nodo) LoteID() *string         { return n.loteID }
func (n *Nodo) TenantID() string        { return n.tenantID }
func (n *Nodo) Nombre() string          { return n.nombre }
func (n *Nodo) NodeKey() string         { return n.nodeKey }
func (n *Nodo) Estado() EstadoNodo      { return n.estado }
func (n *Nodo) CreadoEn() time.Time     { return n.creadoEn }
func (n *Nodo) ActualizadoEn() time.Time { return n.actualizadoEn }

func (n *Nodo) IsActivo() bool {
	return n.estado == NodoActivo
}

func (n *Nodo) Actualizar(loteID *string, nombre string) {
	n.loteID = loteID
	n.nombre = nombre
	n.actualizadoEn = time.Now()
}

func (n *Nodo) CambiarEstado(siguiente EstadoNodo) error {
	if !transicionesNodo[n.estado][siguiente] {
		return ErrTransicionEstadoNoPermitida
	}
	n.estado = siguiente
	n.actualizadoEn = time.Now()
	return nil
}

func NewNodoFromPersistence(
	id, tenantID, fincaID, nodeKey string,
	loteID *string,
	nombre string,
	estado EstadoNodo,
	creadoEn, actualizadoEn time.Time,
) *Nodo {
	return &Nodo{
		id:        id,
		tenantID:  tenantID,
		fincaID:   fincaID,
		nodeKey:   nodeKey,
		loteID:    loteID,
		nombre:    nombre,
		estado:    estado,
		creadoEn:  creadoEn,
		actualizadoEn: actualizadoEn,
	}
}

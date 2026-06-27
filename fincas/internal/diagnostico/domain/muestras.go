package domain

import (
	"errors"
	"time"
)

// la muestra representa a un elemento que representa a una zona dentro de un lote
type Muestra struct {
	id        string
	fincaID   string
	loteID    string
	tenantID  string
	ubicacion Ubicacion
	createdAt time.Time
	updatedAt time.Time
}

type Ubicacion struct {
	latitud  float64
	longitud float64
}

func NewMuestraSinUbicacion(id, fincaID, loteID, tenantID string) (*Muestra, error) {
	if fincaID == "" {
		return nil, errors.New("el fincaID es requerido")
	}
	if tenantID == "" {
		return nil, errors.New("el tenantID es requerido")
	}
	return &Muestra{
		id:        id,
		fincaID:   fincaID,
		loteID:    loteID,
		tenantID:  tenantID,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func NewMuestra(id, fincaID, loteID string, ubicacion Ubicacion, tenantID string) (*Muestra, error) {
	if fincaID == "" {
		return nil, errors.New("el fincaID es requerido")
	}
	if tenantID == "" {
		return nil, errors.New("el tenantID es requerido")
	}
	if ubicacion.latitud == 0 || ubicacion.longitud == 0 {
		return nil, errors.New("la ubicacion es requerida")
	}
	return &Muestra{
		id:        id,
		fincaID:   fincaID,
		loteID:    loteID,
		ubicacion: ubicacion,
		tenantID:  tenantID,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// una vez fijada la ubicacion no se puede editar, Justiicacion, un cambio de ubicacion conlleva un cambio de
// suelo por tanto las historia previa se perderia
func NewUbicacion(latitud, longitud float64) (*Ubicacion, error) {
	if latitud < -90 || latitud > 90 {
		return nil, errors.New("latitud inválida")
	}
	if longitud < -180 || longitud > 180 {
		return nil, errors.New("longitud inválida")
	}
	return &Ubicacion{
		latitud:  latitud,
		longitud: longitud,
	}, nil
}

func (m *Muestra) ID() string           { return m.id }
func (m *Muestra) FincaID() string      { return m.fincaID }
func (m *Muestra) LoteID() string       { return m.loteID }
func (m *Muestra) TenantID() string     { return m.tenantID }
func (m *Muestra) Ubicacion() Ubicacion { return m.ubicacion }
func (m *Muestra) CreatedAt() time.Time { return m.createdAt }

func (u *Ubicacion) Latitud() float64  { return u.latitud }
func (u *Ubicacion) Longitud() float64 { return u.longitud }

func NewMusetraFromStorage(id, fincaID, loteID, tenantId string, ubicacion Ubicacion, createdAt, updatedAt time.Time) (*Muestra, error) {
	return &Muestra{
		id:        id,
		fincaID:   fincaID,
		loteID:    loteID,
		tenantID:  tenantId,
		ubicacion: ubicacion,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

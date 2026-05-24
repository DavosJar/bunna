package domain

import (
	"time"
)

type Lote struct {
	ID          string
	FincaID     string
	Nombre      string
	Area        float64
	Descripcion string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewLote(fincaID, nombre string, area float64, descripcion string) (*Lote, error) {
	l := &Lote{
		FincaID:     fincaID,
		Nombre:      nombre,
		Area:        area,
		Descripcion: descripcion,
	}

	if err := l.validar(); err != nil {
		return nil, err
	}

	return l, nil
}

func NewLoteFromPersistence(id, fincaID, nombre string, area float64, descripcion string, createdAt, updatedAt time.Time) *Lote {
	return &Lote{
		ID:          id,
		FincaID:     fincaID,
		Nombre:      nombre,
		Area:        area,
		Descripcion: descripcion,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (l *Lote) validar() error {
	if len(l.Nombre) < 3 {
		return ErrNombreLoteRequerido
	}
	if len(l.Nombre) > 200 {
		return ErrNombreLoteLargo
	}
	if l.Area <= 0 {
		return ErrAreaRequerida
	}
	if len(l.Descripcion) > 1000 {
		return ErrDescripcionLarga
	}
	if l.FincaID == "" {
		return ErrFincaNoEncontrada
	}
	return nil
}

func (l *Lote) Actualizar(nombre string, area float64, descripcion string) error {
	origNombre, origArea, origDescripcion := l.Nombre, l.Area, l.Descripcion

	l.Nombre = nombre
	l.Area = area
	l.Descripcion = descripcion

	if err := l.validar(); err != nil {
		l.Nombre = origNombre
		l.Area = origArea
		l.Descripcion = origDescripcion
		return err
	}

	l.UpdatedAt = time.Now()
	return nil
}

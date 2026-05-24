package domain

import (
	"time"
)

type Finca struct {
	ID          string
	Nombre      string
	Ubicacion   string
	Descripcion string
	UsuarioID   string
	TenantID    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewFinca(nombre, ubicacion, descripcion, usuarioID string) (*Finca, error) {
	f := &Finca{
		Nombre:      nombre,
		Ubicacion:   ubicacion,
		Descripcion: descripcion,
		UsuarioID:   usuarioID,
	}

	if err := f.validar(); err != nil {
		return nil, err
	}

	return f, nil
}

func NewFincaFromPersistence(id, nombre, ubicacion, descripcion, usuarioID string, tenantID *string, createdAt, updatedAt time.Time) *Finca {
	return &Finca{
		ID:          id,
		Nombre:      nombre,
		Ubicacion:   ubicacion,
		Descripcion: descripcion,
		UsuarioID:   usuarioID,
		TenantID:    tenantID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (f *Finca) validar() error {
	if len(f.Nombre) < 3 {
		return ErrNombreFincaRequerido
	}
	if len(f.Nombre) > 200 {
		return ErrNombreFincaLargo
	}
	if f.Ubicacion == "" {
		return ErrUbicacionRequerida
	}
	if len(f.Ubicacion) > 500 {
		return ErrUbicacionLarga
	}
	if len(f.Descripcion) > 1000 {
		return ErrDescripcionLarga
	}
	if f.UsuarioID == "" {
		return ErrNoPropietario
	}
	return nil
}

func (f *Finca) Actualizar(nombre, ubicacion, descripcion string) error {
	origNombre, origUbicacion, origDescripcion := f.Nombre, f.Ubicacion, f.Descripcion

	f.Nombre = nombre
	f.Ubicacion = ubicacion
	f.Descripcion = descripcion

	if err := f.validar(); err != nil {
		f.Nombre = origNombre
		f.Ubicacion = origUbicacion
		f.Descripcion = origDescripcion
		return err
	}

	f.UpdatedAt = time.Now()
	return nil
}

func (f *Finca) EsPropietario(usuarioID string) bool {
	return f.UsuarioID == usuarioID
}

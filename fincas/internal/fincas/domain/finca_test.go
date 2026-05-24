package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewFinca_OK(t *testing.T) {
	f, err := NewFinca("Mi Finca", "Calle 123", "una descripción", uuid.NewString())
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, "Mi Finca", f.Nombre)
	assert.Equal(t, "Calle 123", f.Ubicacion)
	assert.Equal(t, "una descripción", f.Descripcion)
	assert.Empty(t, f.ID)
	assert.True(t, f.CreatedAt.IsZero())
}

func TestNewFinca_ErrNombreCorto(t *testing.T) {
	f, err := NewFinca("ab", "ubicación", "", "user-1")
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrNombreFincaRequerido)
}

func TestNewFinca_ErrNombreLargo(t *testing.T) {
	nombre := string(make([]byte, 201))
	for i := range nombre {
		nombre = nombre[:i] + "a" + nombre[i+1:]
	}
	f, err := NewFinca(nombre, "ubicación", "", "user-1")
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrNombreFincaLargo)
}

func TestNewFinca_ErrUbicacionVacia(t *testing.T) {
	f, err := NewFinca("Mi Finca", "", "", "user-1")
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrUbicacionRequerida)
}

func TestNewFinca_ErrUbicacionLarga(t *testing.T) {
	ubicacion := string(make([]byte, 501))
	for i := range ubicacion {
		ubicacion = ubicacion[:i] + "a" + ubicacion[i+1:]
	}
	f, err := NewFinca("Mi Finca", ubicacion, "", "user-1")
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrUbicacionLarga)
}

func TestNewFinca_ErrDescripcionLarga(t *testing.T) {
	desc := string(make([]byte, 1001))
	for i := range desc {
		desc = desc[:i] + "a" + desc[i+1:]
	}
	f, err := NewFinca("Mi Finca", "ubicación", desc, "user-1")
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrDescripcionLarga)
}

func TestNewFinca_ErrUsuarioIDVacio(t *testing.T) {
	f, err := NewFinca("Mi Finca", "ubicación", "", "")
	assert.Nil(t, f)
	assert.ErrorIs(t, err, ErrNoPropietario)
}

func TestFinca_Actualizar_OK(t *testing.T) {
	f, _ := NewFinca("Mi Finca", "Calle 123", "desc", "user-1")
	now := time.Now()
	f.UpdatedAt = now

	err := f.Actualizar("Nuevo Nombre", "Nueva Ubicación", "nueva desc")
	assert.NoError(t, err)
	assert.Equal(t, "Nuevo Nombre", f.Nombre)
	assert.Equal(t, "Nueva Ubicación", f.Ubicacion)
	assert.Equal(t, "nueva desc", f.Descripcion)
	assert.True(t, f.UpdatedAt.After(now))
}

func TestFinca_Actualizar_Rollback(t *testing.T) {
	f, _ := NewFinca("Mi Finca", "Calle 123", "desc", "user-1")

	err := f.Actualizar("ab", "Nueva Ubicación", "")
	assert.Error(t, err)
	assert.Equal(t, "Mi Finca", f.Nombre)
	assert.Equal(t, "Calle 123", f.Ubicacion)
	assert.Equal(t, "desc", f.Descripcion)
}

func TestFinca_EsPropietario(t *testing.T) {
	f, _ := NewFinca("Mi Finca", "ubicación", "", "user-1")
	assert.True(t, f.EsPropietario("user-1"))
	assert.False(t, f.EsPropietario("user-2"))
	assert.False(t, f.EsPropietario(""))
}

func TestNewFincaFromPersistence(t *testing.T) {
	now := time.Now()
	f := NewFincaFromPersistence(
		"fin-1", "Mi Finca", "ubicación", "desc", "user-1",
		nil, now, now,
	)
	assert.Equal(t, "fin-1", f.ID)
	assert.Equal(t, "Mi Finca", f.Nombre)
	assert.Equal(t, now, f.CreatedAt)
	assert.Nil(t, f.TenantID)
}

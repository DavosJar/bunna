package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLote_OK(t *testing.T) {
	l, err := NewLote("fin-1", "Lote A", 15.5, "lote frente al río")
	assert.NoError(t, err)
	assert.NotNil(t, l)
	assert.Equal(t, "fin-1", l.FincaID)
	assert.Equal(t, "Lote A", l.Nombre)
	assert.Equal(t, 15.5, l.Area)
	assert.Equal(t, "lote frente al río", l.Descripcion)
	assert.Empty(t, l.ID)
	assert.True(t, l.CreatedAt.IsZero())
}

func TestNewLote_ErrNombreCorto(t *testing.T) {
	l, err := NewLote("fin-1", "ab", 10, "")
	assert.Nil(t, l)
	assert.ErrorIs(t, err, ErrNombreLoteRequerido)
}

func TestNewLote_ErrNombreLargo(t *testing.T) {
	nombre := string(make([]byte, 201))
	for i := range nombre {
		nombre = nombre[:i] + "a" + nombre[i+1:]
	}
	l, err := NewLote("fin-1", nombre, 10, "")
	assert.Nil(t, l)
	assert.ErrorIs(t, err, ErrNombreLoteLargo)
}

func TestNewLote_ErrAreaCero(t *testing.T) {
	l, err := NewLote("fin-1", "Lote A", 0, "")
	assert.Nil(t, l)
	assert.ErrorIs(t, err, ErrAreaRequerida)
}

func TestNewLote_ErrAreaNegativa(t *testing.T) {
	l, err := NewLote("fin-1", "Lote A", -5, "")
	assert.Nil(t, l)
	assert.ErrorIs(t, err, ErrAreaRequerida)
}

func TestNewLote_ErrDescripcionLarga(t *testing.T) {
	desc := string(make([]byte, 1001))
	for i := range desc {
		desc = desc[:i] + "a" + desc[i+1:]
	}
	l, err := NewLote("fin-1", "Lote A", 10, desc)
	assert.Nil(t, l)
	assert.ErrorIs(t, err, ErrDescripcionLarga)
}

func TestNewLote_ErrFincaIDVacio(t *testing.T) {
	l, err := NewLote("", "Lote A", 10, "")
	assert.Nil(t, l)
	assert.ErrorIs(t, err, ErrFincaIDRequerido)
}

func TestLote_Actualizar_OK(t *testing.T) {
	l, _ := NewLote("fin-1", "Lote A", 10, "desc")
	now := time.Now()
	l.UpdatedAt = now

	err := l.Actualizar("Lote B", 20.5, "nueva desc")
	assert.NoError(t, err)
	assert.Equal(t, "Lote B", l.Nombre)
	assert.Equal(t, 20.5, l.Area)
	assert.Equal(t, "nueva desc", l.Descripcion)
	assert.True(t, l.UpdatedAt.After(now))
}

func TestLote_Actualizar_Rollback(t *testing.T) {
	l, _ := NewLote("fin-1", "Lote A", 10, "desc")

	err := l.Actualizar("ab", 20, "")
	assert.Error(t, err)
	assert.Equal(t, "Lote A", l.Nombre)
	assert.Equal(t, 10.0, l.Area)
	assert.Equal(t, "desc", l.Descripcion)
}

func TestNewLoteFromPersistence(t *testing.T) {
	now := time.Now()
	l := NewLoteFromPersistence("lot-1", "fin-1", "Lote A", 15.5, "desc", now, now)
	assert.Equal(t, "lot-1", l.ID)
	assert.Equal(t, "fin-1", l.FincaID)
	assert.Equal(t, 15.5, l.Area)
	assert.Equal(t, now, l.CreatedAt)
}

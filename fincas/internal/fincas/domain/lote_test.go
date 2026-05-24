package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNuevoLote_OK(t *testing.T) {
	l := NuevoLote("fin-1", "Lote A", 15.5, "lote frente al río")
	assert.NotNil(t, l)
	assert.Equal(t, "fin-1", l.FincaID())
	assert.Equal(t, "Lote A", l.Nombre())
	assert.Equal(t, 15.5, l.Area())
	assert.Equal(t, "lote frente al río", l.Descripcion())
	assert.Equal(t, LoteActivo, l.Estado())
	assert.Empty(t, l.ID())
}

func TestLote_Actualizar(t *testing.T) {
	l := NuevoLote("fin-1", "Lote A", 10, "desc")

	l.Actualizar("Lote B", 20.5, "nueva desc")
	assert.Equal(t, "Lote B", l.Nombre())
	assert.Equal(t, 20.5, l.Area())
	assert.Equal(t, "nueva desc", l.Descripcion())
}

func TestLote_CambiarEstado_Valido(t *testing.T) {
	l := NuevoLote("fin-1", "Lote A", 10, "")
	assert.Equal(t, LoteActivo, l.Estado())

	err := l.CambiarEstado(LoteEliminado)
	assert.NoError(t, err)
	assert.Equal(t, LoteEliminado, l.Estado())
}

func TestLote_CambiarEstado_Invalido(t *testing.T) {
	l := NuevoLote("fin-1", "Lote A", 10, "")

	err := l.CambiarEstado(LoteActivo)
	assert.ErrorIs(t, err, ErrTransicionEstadoNoPermitida)
}

func TestNewLoteFromPersistence(t *testing.T) {
	now := time.Now()
	l := NewLoteFromPersistence("lot-1", "fin-1", "Lote A", 15.5, "desc", LoteActivo, now, now)
	assert.Equal(t, "lot-1", l.ID())
	assert.Equal(t, "fin-1", l.FincaID())
	assert.Equal(t, 15.5, l.Area())
	assert.Equal(t, LoteActivo, l.Estado())
	assert.Equal(t, now, l.CreatedAt())
}

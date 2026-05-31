package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEliminarFincaConLotes_SinLotes(t *testing.T) {
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")
	svc := NewFincaService()

	err := svc.EliminarFincaConLotes(f, 0, false)
	assert.NoError(t, err)
	assert.Equal(t, FincaPendienteEliminar, f.Estado())
}

func TestEliminarFincaConLotes_Confirmado(t *testing.T) {
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")
	svc := NewFincaService()

	err := svc.EliminarFincaConLotes(f, 3, true)
	assert.NoError(t, err)
	assert.Equal(t, FincaPendienteEliminar, f.Estado())
}

func TestEliminarFincaConLotes_ErrorConLotes(t *testing.T) {
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")
	svc := NewFincaService()

	err := svc.EliminarFincaConLotes(f, 3, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "3")
	assert.Contains(t, err.Error(), "lote(s)")
	assert.Equal(t, FincaActiva, f.Estado())
}

package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNuevaFinca_OK(t *testing.T) {
	f := NuevaFinca("Mi Finca", "Calle 123", "una descripción", "user-1")
	assert.NotNil(t, f)
	assert.Equal(t, "Mi Finca", f.Nombre())
	assert.Equal(t, "Calle 123", f.Ubicacion())
	assert.Equal(t, "una descripción", f.Descripcion())
	assert.Equal(t, "user-1", f.UsuarioID())
	assert.Equal(t, FincaActiva, f.Estado())
	assert.Empty(t, f.ID())
}

func TestFinca_EsPropietario(t *testing.T) {
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")
	assert.True(t, f.EsPropietario("user-1"))
	assert.False(t, f.EsPropietario("user-2"))
	assert.False(t, f.EsPropietario(""))
}

func TestFinca_TieneLotes(t *testing.T) {
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")
	assert.True(t, f.TieneLotes(1))
	assert.True(t, f.TieneLotes(5))
	assert.False(t, f.TieneLotes(0))
}

func TestFinca_Actualizar(t *testing.T) {
	f := NuevaFinca("Mi Finca", "Calle 123", "desc", "user-1")

	f.Actualizar("Nuevo Nombre", "Nueva Ubicación", "nueva desc")
	assert.Equal(t, "Nuevo Nombre", f.Nombre())
	assert.Equal(t, "Nueva Ubicación", f.Ubicacion())
	assert.Equal(t, "nueva desc", f.Descripcion())
}

func TestFinca_CambiarEstado_Valido(t *testing.T) {
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")
	assert.Equal(t, FincaActiva, f.Estado())

	err := f.CambiarEstado(FincaPendienteEliminar)
	assert.NoError(t, err)
	assert.Equal(t, FincaPendienteEliminar, f.Estado())
}

func TestFinca_CambiarEstado_Invalido(t *testing.T) {
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")

	// No se puede ir de ACTIVA al mismo estado
	err := f.CambiarEstado(FincaActiva)
	assert.ErrorIs(t, err, ErrTransicionEstadoNoPermitida)
}

func TestFinca_PendienteEliminarEsTerminal(t *testing.T) {
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")
	f.CambiarEstado(FincaPendienteEliminar)

	// No debería poder salir de PENDIENTE_ELIMINACION
	err := f.CambiarEstado(FincaActiva)
	assert.ErrorIs(t, err, ErrTransicionEstadoNoPermitida)
}

func TestFinca_TenantIDGetter(t *testing.T) {
	tenant := "tenant-1"
	f := NuevaFinca("Mi Finca", "ubicación", "", "user-1")
	// Asignamos tenantID directamente (como pasaría desde persistencia)
	// Verificamos que el getter existe y funciona
	f = NewFincaFromPersistence(
		"fin-1", "Mi Finca", "ubicación", "desc", "user-1",
		&tenant, FincaActiva, time.Now(), time.Now(),
	)
	assert.NotNil(t, f.TenantID())
	assert.Equal(t, "tenant-1", *f.TenantID())
}

func TestNewFincaFromPersistence(t *testing.T) {
	now := time.Now()
	f := NewFincaFromPersistence(
		"fin-1", "Mi Finca", "ubicación", "desc", "user-1",
		nil, FincaActiva, now, now,
	)
	assert.Equal(t, "fin-1", f.ID())
	assert.Equal(t, "Mi Finca", f.Nombre())
	assert.Equal(t, FincaActiva, f.Estado())
	assert.Equal(t, now, f.CreatedAt())
	assert.Nil(t, f.TenantID())
}

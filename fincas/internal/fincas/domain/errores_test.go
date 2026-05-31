package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrFincaConLotes(t *testing.T) {
	err := ErrFincaConLotes(3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "3")
	assert.Contains(t, err.Error(), "lote(s)")
}

func TestErroresSentinel(t *testing.T) {
	assert.True(t, errors.Is(ErrFincaNoEncontrada, ErrFincaNoEncontrada))
	assert.True(t, errors.Is(ErrLoteNoEncontrado, ErrLoteNoEncontrado))
	assert.True(t, errors.Is(ErrNoPropietario, ErrNoPropietario))
	assert.True(t, errors.Is(ErrTransicionEstadoNoPermitida, ErrTransicionEstadoNoPermitida))
}

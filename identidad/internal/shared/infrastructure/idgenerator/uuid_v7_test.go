package idgenerator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUUIDv7GeneratorNextID(t *testing.T) {
	gen := NewUUIDv7Generator()

	id, err := gen.NextID(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.Len(t, id, 36) // UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
}

func TestUUIDv7GeneratorNextIDMultiple(t *testing.T) {
	gen := NewUUIDv7Generator()

	id1, _ := gen.NextID(context.Background())
	id2, _ := gen.NextID(context.Background())
	id3, _ := gen.NextID(context.Background())

	// Los IDs deben ser únicos
	require.NotEqual(t, id1, id2)
	require.NotEqual(t, id2, id3)
	require.NotEqual(t, id1, id3)
}

func TestUUIDv7GeneratorContextCancel(t *testing.T) {
	gen := NewUUIDv7Generator()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// El generador debe funcionar incluso si ctx está cancelado
	// (la generación de UUID es local, no IO-bound)
	id, err := gen.NextID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

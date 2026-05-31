package idgenerator

import (
	"context"

	"github.com/google/uuid"
)

type GeneradorUUIDV7 struct{}

func NewGeneradorUUIDV7() *GeneradorUUIDV7 {
	return &GeneradorUUIDV7{}
}

func (g *GeneradorUUIDV7) NextID(_ context.Context) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

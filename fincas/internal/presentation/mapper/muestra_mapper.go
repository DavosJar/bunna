package mapper

import (
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/listarmuestrasporlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/tomarmuestra"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
)

// MuestraMapper convierte entre capa de aplicación y DTOs de muestra.
type MuestraMapper struct{}

func (MuestraMapper) TomarSalidaToResponse(s *tomarmuestra.Salida) dto.MuestraResponse {
	return dto.MuestraResponse{
		ID:        s.ID,
		LoteID:    s.LoteID,
		Latitud:   s.Latitud,
		Longitud:  s.Longitud,
		CreatedAt: s.CreatedAt,
	}
}

func (MuestraMapper) MuestraItemToResponse(item listarmuestrasporlote.MuestraItem) dto.MuestraItemResponse {
	return dto.MuestraItemResponse{
		ID:        item.ID,
		LoteID:    item.LoteID,
		Latitud:   item.Latitud,
		Longitud:  item.Longitud,
		CreatedAt: item.CreatedAt,
	}
}

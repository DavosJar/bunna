package mapper

import (
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/desactivarfinca"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarfinca"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
)

// FincaMapper convierte entre capa de aplicación y DTOs de finca.
type FincaMapper struct{}

func (FincaMapper) RegistrarSalidaToResponse(s *registrarfinca.Salida) dto.FincaResponse {
	return dto.FincaResponse{
		ID:          s.ID,
		Nombre:      s.Nombre,
		Ubicacion:   s.Ubicacion,
		Descripcion: s.Descripcion,
		Estado:      s.Estado,
		CreatedAt:   s.CreatedAt,
	}
}

func (FincaMapper) DesactivarSalidaToEstadoCambio(s *desactivarfinca.Salida) dto.EstadoCambioResponse {
	return dto.EstadoCambioResponse{
		ID:        s.ID,
		Estado:    s.Estado,
		UpdatedAt: s.UpdatedAt,
	}
}

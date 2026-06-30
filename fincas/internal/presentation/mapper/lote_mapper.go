package mapper

import (
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/agregarlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/eliminarlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/listarlotes"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
)

// LoteMapper convierte entre capa de aplicación y DTOs de lote.
type LoteMapper struct{}

func (LoteMapper) AgregarSalidaToResponse(s *agregarlote.Salida) dto.LoteResponse {
	return dto.LoteResponse{
		ID:          s.ID,
		FincaID:     s.FincaID,
		Nombre:      s.Nombre,
		Area:        s.Area,
		Descripcion: s.Descripcion,
		Estado:      s.Estado,
		CreatedAt:   s.CreatedAt,
	}
}

func (LoteMapper) EliminarSalidaToEstadoCambio(s *eliminarlote.Salida) dto.EstadoCambioResponse {
	return dto.EstadoCambioResponse{
		ID:        s.ID,
		Estado:    s.Estado,
		UpdatedAt: s.UpdatedAt,
	}
}

func (LoteMapper) ListarSalidaToResponse(s *listarlotes.Salida) []dto.LoteResponse {
	resp := make([]dto.LoteResponse, len(s.Lotes))
	for i, l := range s.Lotes {
		resp[i] = dto.LoteResponse{
			ID:          l.ID,
			FincaID:     l.FincaID,
			Nombre:      l.Nombre,
			Area:        l.Area,
			Descripcion: l.Descripcion,
			Estado:      l.Estado,
			CreatedAt:   l.CreatedAt,
		}
	}
	return resp
}

package mapper

import (
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/desactivarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/editarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/listarnodos"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/obtenernodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarinferenciadesdenodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/validarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
)

type NodoMapper struct{}

func (NodoMapper) RegistrarSalidaToResponse(s *registrarnodo.Salida) dto.NodoResponse {
	return dto.NodoResponse{
		ID:       s.ID,
		FincaID:  s.FincaID,
		LoteID:   s.LoteID,
		TenantID: s.TenantID,
		Nombre:   s.Nombre,
		NodeKey:  s.NodeKey,
		Estado:   s.Estado,
		CreadoEn: s.CreadoEn,
	}
}

func (NodoMapper) ObtenerSalidaToResponse(s *obtenernodo.Salida) dto.NodoResponse {
	return dto.NodoResponse{
		ID:            s.ID,
		FincaID:       s.FincaID,
		LoteID:        s.LoteID,
		TenantID:      s.TenantID,
		Nombre:        s.Nombre,
		NodeKey:       s.NodeKey,
		Estado:        s.Estado,
		CreadoEn:      s.CreadoEn,
		ActualizadoEn: s.ActualizadoEn,
	}
}

func (NodoMapper) ListarItemToResponse(item listarnodos.NodoItem) dto.NodoResponse {
	return dto.NodoResponse{
		ID:       item.ID,
		FincaID:  item.FincaID,
		LoteID:   item.LoteID,
		TenantID: item.TenantID,
		Nombre:   item.Nombre,
		NodeKey:  item.NodeKey,
		Estado:   item.Estado,
		CreadoEn: item.CreadoEn,
	}
}

func (NodoMapper) DesactivarSalidaToEstadoCambio(s *desactivarnodo.Salida) dto.EstadoCambioResponse {
	return dto.EstadoCambioResponse{
		ID:        s.ID,
		Estado:    s.Estado,
		UpdatedAt: s.ActualizadoEn,
	}
}

func (NodoMapper) EditarSalidaToResponse(s *editarnodo.Salida) dto.NodoResponse {
	return dto.NodoResponse{
		ID:            s.ID,
		FincaID:       s.FincaID,
		LoteID:        s.LoteID,
		TenantID:      s.TenantID,
		Nombre:        s.Nombre,
		NodeKey:       s.NodeKey,
		Estado:        s.Estado,
		ActualizadoEn: s.ActualizadoEn,
	}
}

func (NodoMapper) ValidarSalidaToResponse(s *validarnodo.Salida) dto.ValidarNodoResponse {
	return dto.ValidarNodoResponse{
		NodoID:   s.NodoID,
		FincaID:  s.FincaID,
		LoteID:   s.LoteID,
		TenantID: s.TenantID,
	}
}

func (NodoMapper) InferenciaSalidaToResponse(s *registrarinferenciadesdenodo.Salida) dto.InferenciaResponse {
	return dto.InferenciaResponse{
		MuestraID:     s.MuestraID,
		DiagnosticoID: s.DiagnosticoID,
		Estado:        s.Estado,
		TieneClorosis: s.TieneClorosis,
		Confianza:     s.Confianza,
		ImageURL:      s.ImageURL,
		CreatedAt:     s.CreatedAt,
	}
}

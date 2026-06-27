package facades

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/listarmuestrasporlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/tomarmuestra"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/mapper"
)

type (
	TomarMuestraUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd tomarmuestra.Command) (*tomarmuestra.Salida, error)
	}
	ListarMuestrasPorLoteUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd listarmuestrasporlote.Command) ([]listarmuestrasporlote.MuestraItem, error)
	}
)

type MuestrasFacade interface {
	Tomar(ctx context.Context, auth *application.AuthContext, fincaID, loteID string, req dto.TomarMuestraRequest) (*dto.MuestraResponse, error)
	ListarPorLote(ctx context.Context, auth *application.AuthContext, fincaID, loteID string) ([]dto.MuestraItemResponse, error)
}

type muestrasFacade struct {
	tomarUC  TomarMuestraUseCase
	listarUC ListarMuestrasPorLoteUseCase
	mapper   mapper.MuestraMapper
}

func NewMuestrasFacade(tomarUC TomarMuestraUseCase, listarUC ListarMuestrasPorLoteUseCase) MuestrasFacade {
	return &muestrasFacade{
		tomarUC:  tomarUC,
		listarUC: listarUC,
		mapper:   mapper.MuestraMapper{},
	}
}

func (f *muestrasFacade) Tomar(ctx context.Context, auth *application.AuthContext, fincaID, loteID string, req dto.TomarMuestraRequest) (*dto.MuestraResponse, error) {
	cmd := tomarmuestra.Command{
		FincaID:  fincaID,
		LoteID:   loteID,
		Latitud:  req.Latitud,
		Longitud: req.Longitud,
	}
	salida, err := f.tomarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.TomarSalidaToResponse(salida)
	return &resp, nil
}

func (f *muestrasFacade) ListarPorLote(ctx context.Context, auth *application.AuthContext, fincaID, loteID string) ([]dto.MuestraItemResponse, error) {
	cmd := listarmuestrasporlote.Command{FincaID: fincaID, LoteID: loteID}
	items, err := f.listarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	result := make([]dto.MuestraItemResponse, len(items))
	for i, item := range items {
		result[i] = f.mapper.MuestraItemToResponse(item)
	}
	return result, nil
}

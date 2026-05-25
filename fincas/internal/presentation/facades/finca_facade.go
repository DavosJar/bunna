package facades

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/desactivarfinca"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarfinca"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/mapper"
)

// Use case interfaces
type (
	RegistrarFincaUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd registrarfinca.Command) (*registrarfinca.Salida, error)
	}
	DesactivarFincaUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd desactivarfinca.Command) (*desactivarfinca.Salida, error)
	}
)

// FincasFacade define la API de presentación para operaciones sobre fincas.
type FincasFacade interface {
	Registrar(ctx context.Context, auth *application.AuthContext, req dto.RegistrarFincaRequest) (*dto.FincaResponse, error)
	Desactivar(ctx context.Context, auth *application.AuthContext, fincaID string, req dto.DesactivarFincaRequest) (*dto.EstadoCambioResponse, error)
}

type fincasFacade struct {
	registrarUC  RegistrarFincaUseCase
	desactivarUC DesactivarFincaUseCase
	mapper       mapper.FincaMapper
}

func NewFincasFacade(registrarUC RegistrarFincaUseCase, desactivarUC DesactivarFincaUseCase) FincasFacade {
	return &fincasFacade{
		registrarUC:  registrarUC,
		desactivarUC: desactivarUC,
		mapper:       mapper.FincaMapper{},
	}
}

func (f *fincasFacade) Registrar(ctx context.Context, auth *application.AuthContext, req dto.RegistrarFincaRequest) (*dto.FincaResponse, error) {
	cmd := registrarfinca.Command{
		Nombre:      req.Nombre,
		Ubicacion:   req.Ubicacion,
		Descripcion: req.Descripcion,
	}
	salida, err := f.registrarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.RegistrarSalidaToResponse(salida)
	return &resp, nil
}

func (f *fincasFacade) Desactivar(ctx context.Context, auth *application.AuthContext, fincaID string, req dto.DesactivarFincaRequest) (*dto.EstadoCambioResponse, error) {
	cmd := desactivarfinca.Command{
		FincaID:   fincaID,
		Confirmar: req.Confirmar,
	}
	salida, err := f.desactivarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.DesactivarSalidaToEstadoCambio(salida)
	return &resp, nil
}

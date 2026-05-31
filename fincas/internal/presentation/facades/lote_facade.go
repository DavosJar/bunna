package facades

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/agregarlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/eliminarlote"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/mapper"
)

type (
	AgregarLoteUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd agregarlote.Command) (*agregarlote.Salida, error)
	}
	EliminarLoteUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd eliminarlote.Command) (*eliminarlote.Salida, error)
	}
)

type LotesFacade interface {
	Agregar(ctx context.Context, auth *application.AuthContext, fincaID string, req dto.AgregarLoteRequest) (*dto.LoteResponse, error)
	Eliminar(ctx context.Context, auth *application.AuthContext, loteID string) (*dto.EstadoCambioResponse, error)
}

type lotesFacade struct {
	agregarUC  AgregarLoteUseCase
	eliminarUC EliminarLoteUseCase
	mapper     mapper.LoteMapper
}

func NewLotesFacade(agregarUC AgregarLoteUseCase, eliminarUC EliminarLoteUseCase) LotesFacade {
	return &lotesFacade{
		agregarUC:  agregarUC,
		eliminarUC: eliminarUC,
		mapper:     mapper.LoteMapper{},
	}
}

func (f *lotesFacade) Agregar(ctx context.Context, auth *application.AuthContext, fincaID string, req dto.AgregarLoteRequest) (*dto.LoteResponse, error) {
	cmd := agregarlote.Command{
		FincaID:     fincaID,
		Nombre:      req.Nombre,
		Area:        req.Area,
		Descripcion: req.Descripcion,
	}
	salida, err := f.agregarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.AgregarSalidaToResponse(salida)
	return &resp, nil
}

func (f *lotesFacade) Eliminar(ctx context.Context, auth *application.AuthContext, loteID string) (*dto.EstadoCambioResponse, error) {
	cmd := eliminarlote.Command{LoteID: loteID}
	salida, err := f.eliminarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.EliminarSalidaToEstadoCambio(salida)
	return &resp, nil
}

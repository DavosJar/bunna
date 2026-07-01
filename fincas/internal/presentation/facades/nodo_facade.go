package facades

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/desactivarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/editarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/listarnodos"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/obtenernodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarinferenciadesdenodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/validarnodo"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/mapper"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

type (
	RegistrarNodoUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd registrarnodo.Command) (*registrarnodo.Salida, error)
	}
	ListarNodosUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd listarnodos.Command, paginacion shared.Paginacion) (*listarnodos.Salida, error)
	}
	ObtenerNodoUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd obtenernodo.Command) (*obtenernodo.Salida, error)
	}
	EditarNodoUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd editarnodo.Command) (*editarnodo.Salida, error)
	}
	DesactivarNodoUseCase interface {
		Ejecutar(ctx context.Context, auth *application.AuthContext, cmd desactivarnodo.Command) (*desactivarnodo.Salida, error)
	}
	ValidarNodoUseCase interface {
		Ejecutar(ctx context.Context, cmd validarnodo.Command) (*validarnodo.Salida, error)
	}
	RegistrarInferenciaDesdeNodoUseCase interface {
		Ejecutar(ctx context.Context, cmd registrarinferenciadesdenodo.Command) (*registrarinferenciadesdenodo.Salida, error)
	}
)

type NodosFacade interface {
	Registrar(ctx context.Context, auth *application.AuthContext, req dto.RegistrarNodoRequest) (*dto.NodoResponse, error)
	Listar(ctx context.Context, auth *application.AuthContext, paginacion shared.Paginacion) (*dto.ListaResponse[dto.NodoResponse], error)
	Obtener(ctx context.Context, auth *application.AuthContext, nodoID string) (*dto.NodoResponse, error)
	Editar(ctx context.Context, auth *application.AuthContext, nodoID string, req dto.EditarNodoRequest) (*dto.NodoResponse, error)
	Desactivar(ctx context.Context, auth *application.AuthContext, nodoID string, req dto.DesactivarNodoRequest) (*dto.EstadoCambioResponse, error)
	Validar(ctx context.Context, nodeKey string) (*dto.ValidarNodoResponse, error)
	RegistrarInferencia(ctx context.Context, req dto.RegistrarInferenciaDesdeNodoRequest) (*dto.InferenciaResponse, error)
}

type nodosFacade struct {
	registrarUC  RegistrarNodoUseCase
	listarUC     ListarNodosUseCase
	obtenerUC    ObtenerNodoUseCase
	editarUC     EditarNodoUseCase
	desactivarUC DesactivarNodoUseCase
	validarUC    ValidarNodoUseCase
	inferenciaUC RegistrarInferenciaDesdeNodoUseCase
	mapper       mapper.NodoMapper
}

func NewNodosFacade(
	registrarUC RegistrarNodoUseCase,
	listarUC ListarNodosUseCase,
	obtenerUC ObtenerNodoUseCase,
	editarUC EditarNodoUseCase,
	desactivarUC DesactivarNodoUseCase,
	validarUC ValidarNodoUseCase,
	inferenciaUC RegistrarInferenciaDesdeNodoUseCase,
) NodosFacade {
	return &nodosFacade{
		registrarUC:  registrarUC,
		listarUC:     listarUC,
		obtenerUC:    obtenerUC,
		editarUC:     editarUC,
		desactivarUC: desactivarUC,
		validarUC:    validarUC,
		inferenciaUC: inferenciaUC,
		mapper:       mapper.NodoMapper{},
	}
}

func (f *nodosFacade) Registrar(ctx context.Context, auth *application.AuthContext, req dto.RegistrarNodoRequest) (*dto.NodoResponse, error) {
	cmd := registrarnodo.Command{
		TenantID: auth.TenantID,
		FincaID:  req.FincaID,
		NodeKey:  req.NodeKey,
		LoteID:   req.LoteID,
		Nombre:   req.Nombre,
	}
	salida, err := f.registrarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.RegistrarSalidaToResponse(salida)
	return &resp, nil
}

func (f *nodosFacade) Listar(ctx context.Context, auth *application.AuthContext, paginacion shared.Paginacion) (*dto.ListaResponse[dto.NodoResponse], error) {
	cmd := listarnodos.Command{
		TenantID: auth.TenantID,
	}
	salida, err := f.listarUC.Ejecutar(ctx, auth, cmd, paginacion)
	if err != nil {
		return nil, err
	}

	items := make([]dto.NodoResponse, len(salida.Nodos))
	for i, item := range salida.Nodos {
		items[i] = f.mapper.ListarItemToResponse(item)
	}
	return &dto.ListaResponse[dto.NodoResponse]{
		Data:         items,
		Total:        salida.Total,
		Pagina:       salida.Pagina,
		TotalPaginas: salida.TotalPaginas,
	}, nil
}

func (f *nodosFacade) Obtener(ctx context.Context, auth *application.AuthContext, nodoID string) (*dto.NodoResponse, error) {
	cmd := obtenernodo.Command{NodoID: nodoID}
	salida, err := f.obtenerUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.ObtenerSalidaToResponse(salida)
	return &resp, nil
}

func (f *nodosFacade) Editar(ctx context.Context, auth *application.AuthContext, nodoID string, req dto.EditarNodoRequest) (*dto.NodoResponse, error) {
	cmd := editarnodo.Command{
		NodoID: nodoID,
		LoteID: req.LoteID,
		Nombre: req.Nombre,
	}
	salida, err := f.editarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.EditarSalidaToResponse(salida)
	return &resp, nil
}

func (f *nodosFacade) Desactivar(ctx context.Context, auth *application.AuthContext, nodoID string, req dto.DesactivarNodoRequest) (*dto.EstadoCambioResponse, error) {
	cmd := desactivarnodo.Command{
		NodoID: nodoID,
		Estado: req.Estado,
	}
	salida, err := f.desactivarUC.Ejecutar(ctx, auth, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.DesactivarSalidaToEstadoCambio(salida)
	return &resp, nil
}

func (f *nodosFacade) Validar(ctx context.Context, nodeKey string) (*dto.ValidarNodoResponse, error) {
	cmd := validarnodo.Command{NodeKey: nodeKey}
	salida, err := f.validarUC.Ejecutar(ctx, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.ValidarSalidaToResponse(salida)
	return &resp, nil
}

func (f *nodosFacade) RegistrarInferencia(ctx context.Context, req dto.RegistrarInferenciaDesdeNodoRequest) (*dto.InferenciaResponse, error) {
	cmd := registrarinferenciadesdenodo.Command{
		NodoID:        req.NodoID,
		FincaID:       req.FincaID,
		LoteID:        req.LoteID,
		TenantID:      req.TenantID,
		ImageURL:      req.ImageURL,
		ImageBase64:   req.ImageBase64,
		TieneClorosis: req.TieneClorosis,
		Confianza:     req.Confianza,
		ProcesadoAt:   req.ProcesadoAt,
	}
	salida, err := f.inferenciaUC.Ejecutar(ctx, cmd)
	if err != nil {
		return nil, err
	}
	resp := f.mapper.InferenciaSalidaToResponse(salida)
	return &resp, nil
}

package listarfincas

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

type UseCase struct {
	fincaRepo domain.FincaRepositorio
}

func NewUseCase(fincaRepo domain.FincaRepositorio) *UseCase {
	return &UseCase{
		fincaRepo: fincaRepo,
	}
}

func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, q Query) (*Salida, error) {
	// Listamos fincas usando EspecificacionFinca para filtrar por el tenant actual (del AuthContext)
	var filtros []shared.CriterioFiltro
	if auth.TenantID != "" {
		filtros = append(filtros, shared.CriterioFiltro{
			Campo:    "tenantID",
			Operador: "=",
			Valor:    auth.TenantID,
		})
	} else if auth.UsuarioID != "" {
		// Fallback para usuarios sin tenant (por ej. si no estuviera configurado IAM)
		filtros = append(filtros, shared.CriterioFiltro{
			Campo:    "usuarioID",
			Operador: "=",
			Valor:    auth.UsuarioID,
		})
	}

	especificacion := domain.EspecificacionFinca{Filtros: filtros}
	paginacion := shared.Paginacion{Pagina: 1, TamanoPagina: 1000} // Sin paginar de momento, traer 1000 max

	fincas, err := uc.fincaRepo.Listar(ctx, especificacion, paginacion)
	if err != nil {
		return nil, err
	}

	salida := &Salida{
		Fincas: make([]FincaSalida, len(fincas)),
	}

	for i, f := range fincas {
		salida.Fincas[i] = FincaSalida{
			ID:          f.ID(),
			Nombre:      f.Nombre(),
			Ubicacion:   f.Ubicacion(),
			Descripcion: f.Descripcion(),
			Estado:      string(f.Estado()),
			CreatedAt:   f.CreatedAt(),
		}
	}

	return salida, nil
}

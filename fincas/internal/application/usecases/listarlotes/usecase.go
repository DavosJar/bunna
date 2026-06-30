package listarlotes

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

type UseCase struct {
	loteRepo domain.LoteRepositorio
}

func NewUseCase(loteRepo domain.LoteRepositorio) *UseCase {
	return &UseCase{
		loteRepo: loteRepo,
	}
}

func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, q Query) (*Salida, error) {
	if err := q.Validar(); err != nil {
		return nil, err
	}

	filtros := []shared.CriterioFiltro{
		{
			Campo:    "fincaID",
			Operador: "=",
			Valor:    q.FincaID,
		},
	}

	especificacion := domain.EspecificacionLote{Filtros: filtros}
	paginacion := shared.Paginacion{Pagina: 1, TamanoPagina: 1000}

	lotes, err := uc.loteRepo.Listar(ctx, especificacion, paginacion)
	if err != nil {
		return nil, err
	}

	salida := &Salida{
		Lotes: make([]LoteSalida, len(lotes)),
	}

	for i, l := range lotes {
		salida.Lotes[i] = LoteSalida{
			ID:          l.ID(),
			FincaID:     l.FincaID(),
			Nombre:      l.Nombre(),
			Area:        l.Area(),
			Descripcion: l.Descripcion(),
			Estado:      string(l.Estado()),
			CreatedAt:   l.CreatedAt(),
		}
	}

	return salida, nil
}

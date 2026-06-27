package listarnodos

import (
	"context"
	"math"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	nodosdomain "github.com/davosjar/bunna/services/fincas/internal/nodos/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = application.PermisoConsultarNodos

type UseCase struct {
	nodoRepo nodosdomain.NodoRepositorio
}

func NewUseCase(nodoRepo nodosdomain.NodoRepositorio) *UseCase {
	return &UseCase{nodoRepo: nodoRepo}
}

func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command, paginacion shared.Paginacion) (*Salida, error) {
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	espec := nodosdomain.EspecificacionNodo{
		Filtros: []shared.CriterioFiltro{
			{Campo: "tenantID", Operador: "=", Valor: auth.TenantID},
		},
	}

	nodos, err := uc.nodoRepo.Listar(ctx, espec, paginacion)
	if err != nil {
		return nil, err
	}

	items := make([]NodoItem, len(nodos))
	for i, n := range nodos {
		items[i] = NodoItem{
			ID:       n.ID(),
			FincaID:  n.FincaID(),
			LoteID:   n.LoteID(),
			TenantID: n.TenantID(),
			Nombre:   n.Nombre(),
			NodeKey:  n.NodeKey(),
			Estado:   string(n.Estado()),
			CreadoEn: n.CreadoEn(),
		}
	}

	tamano := paginacion.TamanoPagina
	if tamano < 1 {
		tamano = 10
	}

	return &Salida{
		Nodos:        items,
		Total:        len(nodos),
		Pagina:       paginacion.Pagina,
		TotalPaginas: int(math.Ceil(float64(len(nodos)) / float64(tamano))),
	}, nil
}

package validarnodo

import (
	"context"
	"errors"

	nodosdomain "github.com/davosjar/bunna/services/fincas/internal/nodos/domain"
)

type UseCase struct {
	nodoRepo nodosdomain.NodoRepositorio
}

func NewUseCase(nodoRepo nodosdomain.NodoRepositorio) *UseCase {
	return &UseCase{nodoRepo: nodoRepo}
}

func (uc *UseCase) Ejecutar(ctx context.Context, cmd Command) (*Salida, error) {
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	nodo, err := uc.nodoRepo.ObtenerPorNodeKey(ctx, cmd.NodeKey)
	if err != nil {
		if errors.Is(err, nodosdomain.ErrNodoNoEncontrado) {
			return nil, nodosdomain.ErrNodoNoEncontrado
		}
		return nil, err
	}

	return &Salida{
		NodoID:   nodo.ID(),
		FincaID:  nodo.FincaID(),
		LoteID:   nodo.LoteID(),
		TenantID: nodo.TenantID(),
	}, nil
}

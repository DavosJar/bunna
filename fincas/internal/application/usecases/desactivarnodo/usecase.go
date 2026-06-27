package desactivarnodo

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	nodosdomain "github.com/davosjar/bunna/services/fincas/internal/nodos/domain"
)

const permisoRequerido = application.PermisoDesactivarNodo

type UseCase struct {
	nodoRepo nodosdomain.NodoRepositorio
}

func NewUseCase(nodoRepo nodosdomain.NodoRepositorio) *UseCase {
	return &UseCase{nodoRepo: nodoRepo}
}

func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	nodo, err := uc.nodoRepo.ObtenerPorID(ctx, cmd.NodoID)
	if err != nil {
		if errors.Is(err, nodosdomain.ErrNodoNoEncontrado) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	if auth.TenantID != "" && nodo.TenantID() != auth.TenantID {
		return nil, application.ErrNotFound
	}

	if err := nodo.CambiarEstado(nodosdomain.EstadoNodo(cmd.Estado)); err != nil {
		return nil, application.ErrConflictoEstado(err.Error())
	}

	if err := uc.nodoRepo.Actualizar(ctx, nodo); err != nil {
		return nil, err
	}

	return &Salida{
		ID:            nodo.ID(),
		Estado:        string(nodo.Estado()),
		ActualizadoEn: nodo.ActualizadoEn(),
	}, nil
}

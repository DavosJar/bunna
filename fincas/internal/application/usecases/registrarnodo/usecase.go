package registrarnodo

import (
	"context"
	"errors"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	fincasdomain "github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	nodosdomain "github.com/davosjar/bunna/services/fincas/internal/nodos/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = application.PermisoCrearNodo

type UseCase struct {
	nodoRepo  nodosdomain.NodoRepositorio
	fincaRepo fincasdomain.FincaRepositorio
	generador shared.GeneradorID
}

func NewUseCase(
	nodoRepo nodosdomain.NodoRepositorio,
	fincaRepo fincasdomain.FincaRepositorio,
	generador shared.GeneradorID,
) *UseCase {
	return &UseCase{
		nodoRepo:  nodoRepo,
		fincaRepo: fincaRepo,
		generador: generador,
	}
}

func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	existing, err := uc.nodoRepo.ObtenerPorNodeKey(ctx, cmd.NodeKey)
	if err != nil && !errors.Is(err, nodosdomain.ErrNodoNoEncontrado) {
		return nil, err
	}
	if existing != nil {
		return nil, application.ErrValidacion("el nodeKey ya está registrado")
	}

	finca, err := uc.fincaRepo.ObtenerPorID(ctx, cmd.FincaID)
	if err != nil {
		if errors.Is(err, fincasdomain.ErrFincaNoEncontrada) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	tenantID := auth.TenantID
	if tenantID == "" {
		tenantID = cmd.TenantID
	}

	if finca.TenantID() != nil && *finca.TenantID() != tenantID {
		return nil, application.ErrNotFound
	}

	id, err := uc.generador.NextID(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	nodo := nodosdomain.NewNodoFromPersistence(
		id, tenantID, cmd.FincaID, cmd.NodeKey, cmd.LoteID, cmd.Nombre,
		nodosdomain.NodoActivo, now, now,
	)

	if err := uc.nodoRepo.Crear(ctx, nodo); err != nil {
		return nil, err
	}

	return &Salida{
		ID:       id,
		FincaID:  cmd.FincaID,
		LoteID:   cmd.LoteID,
		TenantID: tenantID,
		Nombre:   cmd.Nombre,
		NodeKey:  cmd.NodeKey,
		Estado:   string(nodo.Estado()),
		CreadoEn: now,
	}, nil
}

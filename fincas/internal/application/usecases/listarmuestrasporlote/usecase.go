package listarmuestrasporlote

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	fincasdomain "github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = application.PermisoVerMuestras

// UseCase orquesta la consulta de muestras asociadas a un lote.
type UseCase struct {
	loteRepo    fincasdomain.LoteRepositorio
	muestraRepo diagnosticodomain.MuestraRepositorio
}

// NewUseCase crea un nuevo caso de uso con las dependencias inyectadas.
func NewUseCase(
	loteRepo fincasdomain.LoteRepositorio,
	muestraRepo diagnosticodomain.MuestraRepositorio,
) *UseCase {
	return &UseCase{
		loteRepo:    loteRepo,
		muestraRepo: muestraRepo,
	}
}

// Ejecutar valida permisos, carga el lote, verifica tenencia, consulta las muestras
// y las retorna como slice de MuestraItem. Es un caso de uso de solo lectura:
// no publica eventos ni usa Unit of Work.
func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) ([]MuestraItem, error) {
	// 1. Validar permisos
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	// 2. Validar campos del comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 3. Cargar lote (opcional)
	if cmd.LoteID != "" {
		lote, err := uc.loteRepo.ObtenerPorID(ctx, cmd.LoteID)
		if err != nil {
			if errors.Is(err, fincasdomain.ErrLoteNoEncontrado) {
				return nil, application.ErrNotFound
			}
			return nil, err
		}

		// 4. Validar tenencia
		if auth.TenantID != "" && lote.TenantID() != auth.TenantID {
			return nil, application.ErrNotFound
		}
	}

	// 5. Construir especificación de búsqueda con filtros
	filtros := []shared.CriterioFiltro{
		{Campo: "fincaID", Operador: "=", Valor: cmd.FincaID},
	}
	
	if cmd.LoteID != "" {
		filtros = append(filtros, shared.CriterioFiltro{Campo: "loteID", Operador: "=", Valor: cmd.LoteID})
	}

	if auth.TenantID != "" {
		filtros = append(filtros, shared.CriterioFiltro{
			Campo: "tenantID", Operador: "=", Valor: auth.TenantID,
		})
	}

	especificacion := diagnosticodomain.EspecificacionMuestra{
		Filtros: filtros,
	}

	// 6. Consultar muestras (sin paginación específica)
	muestras, err := uc.muestraRepo.Buscar(ctx, especificacion, shared.Paginacion{})
	if err != nil {
		return nil, err
	}

	// 7. Mapear resultados — retornar slice vacío en lugar de nil si no hay muestras
	if len(muestras) == 0 {
		return []MuestraItem{}, nil
	}

	resultado := make([]MuestraItem, len(muestras))
	for i, m := range muestras {
		ubicacion := m.Ubicacion()
		resultado[i] = MuestraItem{
			ID:        m.ID(),
			FincaID:   m.FincaID(),
			LoteID:    m.LoteID(),
			Latitud:   ubicacion.Latitud(),
			Longitud:  ubicacion.Longitud(),
			CreatedAt: m.CreatedAt(),
		}
	}

	return resultado, nil
}

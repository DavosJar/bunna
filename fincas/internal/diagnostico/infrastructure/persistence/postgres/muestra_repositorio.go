package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"gorm.io/gorm"
)

type muestraRepositorio struct {
	db *gorm.DB
}

func NewMuestraRepositorio(db *gorm.DB) domain.MuestraRepositorio {
	return &muestraRepositorio{db: db}
}

func (r *muestraRepositorio) Crear(ctx context.Context, muestra *domain.Muestra) error {
	model := FromDomainMuestra(muestra)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *muestraRepositorio) ObtenerPorID(ctx context.Context, id string) (*domain.Muestra, error) {
	var model MuestraModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrMuestraNoEncontrada
		}
		return nil, err
	}
	return model.ToDomain()
}

// ListarPorDiagnostico obtiene la muestra asociada a un diagnóstico.
// La relación es 1:1 (un diagnóstico tiene una muestra), pero retorna slice
// para mantener la interfaz consistente.
func (r *muestraRepositorio) ListarPorDiagnostico(ctx context.Context, diagnosticoID string) ([]domain.Muestra, error) {
	var models []MuestraModel
	err := r.db.WithContext(ctx).
		Where("id IN (SELECT muestras_id FROM diagnosticos WHERE id = ?)", diagnosticoID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	muestras := make([]domain.Muestra, 0, len(models))
	for _, m := range models {
		mu, err := m.ToDomain()
		if err != nil {
			return nil, err
		}
		muestras = append(muestras, *mu)
	}
	return muestras, nil
}

func (r *muestraRepositorio) Actualizar(ctx context.Context, muestra *domain.Muestra) error {
	model := FromDomainMuestra(muestra)
	return r.db.WithContext(ctx).Model(&MuestraModel{}).Where("id = ?", model.ID).Updates(model).Error
}

func (r *muestraRepositorio) Eliminar(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&MuestraModel{}, "id = ?", id).Error
}

func (r *muestraRepositorio) Buscar(
	ctx context.Context,
	especificacion domain.EspecificacionMuestra,
	paginacion shared.Paginacion,
) ([]domain.Muestra, error) {
	query := r.db.WithContext(ctx).Model(&MuestraModel{})

	mapeoColumnas := map[string]string{
		"loteID":   "lote_id",
		"tenantID": "tenant_id",
	}

	for _, filtro := range especificacion.Filtros {
		if !domain.ColumnasPermitidasMuestra[filtro.Campo] {
			continue
		}

		columnaDB, ok := mapeoColumnas[filtro.Campo]
		if !ok {
			continue
		}

		switch filtro.Operador {
		case "=":
			query = query.Where(columnaDB+" = ?", filtro.Valor)
		case "!=":
			query = query.Where(columnaDB+" != ?", filtro.Valor)
		case "LIKE":
			query = query.Where(columnaDB+" LIKE ?", filtro.Valor)
		}
	}

	for _, ord := range paginacion.Ordenaciones {
		if !domain.ColumnasPermitidasMuestra[ord.Campo] {
			continue
		}
		columnaDB, ok := mapeoColumnas[ord.Campo]
		if !ok {
			continue
		}
		orden := "ASC"
		if ord.Tipo == shared.DESC {
			orden = "DESC"
		}
		query = query.Order(columnaDB + " " + orden)
	}

	pagina := paginacion.Pagina
	if pagina < 1 {
		pagina = 1
	}
	tamano := paginacion.TamanoPagina
	if tamano < 1 {
		tamano = 10
	}
	offset := (pagina - 1) * tamano

	var models []MuestraModel
	result := query.Offset(offset).Limit(tamano).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	muestras := make([]domain.Muestra, 0, len(models))
	for _, m := range models {
		mu, err := m.ToDomain()
		if err != nil {
			return nil, err
		}
		muestras = append(muestras, *mu)
	}
	return muestras, nil
}

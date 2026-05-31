package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"gorm.io/gorm"
)

type loteRepositorio struct {
	db *gorm.DB
}

func NewLoteRepositorio(db *gorm.DB) domain.LoteRepositorio {
	return &loteRepositorio{db: db}
}

func (r *loteRepositorio) Crear(ctx context.Context, lote *domain.Lote) error {
	model := FromDomainLote(lote)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *loteRepositorio) ObtenerPorID(ctx context.Context, id string) (*domain.Lote, error) {
	var model LoteModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrLoteNoEncontrado
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *loteRepositorio) ListarPorFinca(ctx context.Context, fincaID string) ([]domain.Lote, error) {
	var models []LoteModel
	err := r.db.WithContext(ctx).Where("finca_id = ?", fincaID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	lotes := make([]domain.Lote, len(models))
	for i, m := range models {
		lotes[i] = *m.ToDomain()
	}
	return lotes, nil
}

func (r *loteRepositorio) Listar(
	ctx context.Context,
	especificacion domain.EspecificacionLote,
	paginacion shared.Paginacion,
) ([]domain.Lote, error) {
	query := r.db.WithContext(ctx).Model(&LoteModel{})

	mapeoColumnas := map[string]string{
		"nombre":  "nombre",
		"area":    "area",
		"estado":  "estado",
		"fincaID": "finca_id",
	}

	for _, filtro := range especificacion.Filtros {
		if !domain.ColumnasPermitidasLotes[filtro.Campo] {
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
		if !domain.ColumnasPermitidasLotes[ord.Campo] {
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

	var models []LoteModel
	result := query.Offset(offset).Limit(tamano).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	lotes := make([]domain.Lote, len(models))
	for i, m := range models {
		lotes[i] = *m.ToDomain()
	}
	return lotes, nil
}

func (r *loteRepositorio) Actualizar(ctx context.Context, lote *domain.Lote) error {
	model := FromDomainLote(lote)
	return r.db.WithContext(ctx).Model(&LoteModel{}).Where("id = ?", model.ID).Updates(model).Error
}

func (r *loteRepositorio) Eliminar(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&LoteModel{}, "id = ?", id).Error
}

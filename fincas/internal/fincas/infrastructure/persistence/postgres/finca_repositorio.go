package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"gorm.io/gorm"
)

type fincaRepositorio struct {
	db *gorm.DB
}

func NewFincaRepositorio(db *gorm.DB) domain.FincaRepositorio {
	return &fincaRepositorio{db: db}
}

func (r *fincaRepositorio) Crear(ctx context.Context, finca *domain.Finca) error {
	model := FromDomainFinca(finca)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *fincaRepositorio) ObtenerPorID(ctx context.Context, id string) (*domain.Finca, error) {
	var model FincaModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFincaNoEncontrada
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *fincaRepositorio) ListarPorUsuario(ctx context.Context, usuarioID string) ([]domain.Finca, error) {
	var models []FincaModel
	err := r.db.WithContext(ctx).Where("usuario_id = ?", usuarioID).Find(&models).Error
	if err != nil {
		return nil, err
	}
	fincas := make([]domain.Finca, len(models))
	for i, m := range models {
		fincas[i] = *m.ToDomain()
	}
	return fincas, nil
}

func (r *fincaRepositorio) ListarTodas(ctx context.Context) ([]domain.Finca, error) {
	var models []FincaModel
	err := r.db.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}
	fincas := make([]domain.Finca, len(models))
	for i, m := range models {
		fincas[i] = *m.ToDomain()
	}
	return fincas, nil
}

func (r *fincaRepositorio) Listar(
	ctx context.Context,
	especificacion domain.EspecificacionFinca,
	paginacion shared.Paginacion,
) ([]domain.Finca, error) {
	query := r.db.WithContext(ctx).Model(&FincaModel{})

	mapeoColumnas := map[string]string{
		"nombre":    "nombre",
		"ubicacion": "ubicacion",
		"estado":    "estado",
		"usuarioID": "usuario_id",
		"tenantID":  "tenant_id",
	}

	for _, filtro := range especificacion.Filtros {
		if !domain.ColumnasPermitidasFincas[filtro.Campo] {
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
		if !domain.ColumnasPermitidasFincas[ord.Campo] {
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

	var models []FincaModel
	result := query.Offset(offset).Limit(tamano).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	fincas := make([]domain.Finca, len(models))
	for i, m := range models {
		fincas[i] = *m.ToDomain()
	}
	return fincas, nil
}

func (r *fincaRepositorio) Actualizar(ctx context.Context, finca *domain.Finca) error {
	model := FromDomainFinca(finca)
	return r.db.WithContext(ctx).Model(&FincaModel{}).Where("id = ?", model.ID).Updates(model).Error
}

func (r *fincaRepositorio) Eliminar(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&FincaModel{}, "id = ?", id).Error
}

func (r *fincaRepositorio) ContarLotes(ctx context.Context, fincaID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&LoteModel{}).Where("finca_id = ?", fincaID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

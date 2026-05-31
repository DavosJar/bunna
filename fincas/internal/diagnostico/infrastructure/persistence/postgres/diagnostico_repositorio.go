package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"gorm.io/gorm"
)

type diagnosticoRepositorio struct {
	db *gorm.DB
}

func NewDiagnosticoRepositorio(db *gorm.DB) domain.DiagnosticoRepositorio {
	return &diagnosticoRepositorio{db: db}
}

func (r *diagnosticoRepositorio) Crear(ctx context.Context, diagnostico *domain.Diagnostico) error {
	model := FromDomainDiagnostico(diagnostico)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *diagnosticoRepositorio) ObtenerPorID(ctx context.Context, id string) (*domain.Diagnostico, error) {
	var model DiagnosticoModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDiagnosticoNoEncontrado
		}
		return nil, err
	}
	return model.ToDomain()
}

// ListarPorFinca busca diagnósticos cuyo muestras_id pertenezca a muestras
// de lotes de la finca indicada. Se resuelve con un sub-select en dos pasos
// para no importar el paquete de fincas (aislamiento de módulos).
func (r *diagnosticoRepositorio) ListarPorFinca(ctx context.Context, fincaID string) ([]domain.Diagnostico, error) {
	var models []DiagnosticoModel
	err := r.db.WithContext(ctx).
		Where("muestras_id IN (SELECT id FROM muestras WHERE lote_id IN (SELECT id FROM lotes WHERE finca_id = ?))", fincaID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	diagnosticos := make([]domain.Diagnostico, 0, len(models))
	for _, m := range models {
		d, err := m.ToDomain()
		if err != nil {
			return nil, err
		}
		diagnosticos = append(diagnosticos, *d)
	}
	return diagnosticos, nil
}

func (r *diagnosticoRepositorio) Actualizar(ctx context.Context, diagnostico *domain.Diagnostico) error {
	model := FromDomainDiagnostico(diagnostico)
	return r.db.WithContext(ctx).Model(&DiagnosticoModel{}).Where("id = ?", model.ID).Updates(model).Error
}

func (r *diagnosticoRepositorio) Eliminar(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&DiagnosticoModel{}, "id = ?", id).Error
}

func (r *diagnosticoRepositorio) Buscar(
	ctx context.Context,
	especificacion domain.EspecificacionDiagnostico,
	paginacion shared.Paginacion,
) ([]domain.Diagnostico, error) {
	query := r.db.WithContext(ctx).Model(&DiagnosticoModel{})

	mapeoColumnas := map[string]string{
		"nombre":    "nombre",
		"muestraID": "muestras_id",
		"tenantID":  "tenant_id",
		"estado":    "estado",
	}

	for _, filtro := range especificacion.Filtros {
		if !domain.ColumnasPermitidasDiagnostico[filtro.Campo] {
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
		if !domain.ColumnasPermitidasDiagnostico[ord.Campo] {
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

	var models []DiagnosticoModel
	result := query.Offset(offset).Limit(tamano).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	diagnosticos := make([]domain.Diagnostico, 0, len(models))
	for _, m := range models {
		d, err := m.ToDomain()
		if err != nil {
			return nil, err
		}
		diagnosticos = append(diagnosticos, *d)
	}
	return diagnosticos, nil
}

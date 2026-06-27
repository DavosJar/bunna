package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/fincas/internal/nodos/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"gorm.io/gorm"
)

type nodoRepositorio struct {
	db *gorm.DB
}

func NewNodoRepositorio(db *gorm.DB) domain.NodoRepositorio {
	return &nodoRepositorio{db: db}
}

func (r *nodoRepositorio) Crear(ctx context.Context, nodo *domain.Nodo) error {
	model := FromDomainNodo(nodo)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *nodoRepositorio) ObtenerPorID(ctx context.Context, id string) (*domain.Nodo, error) {
	var model NodoModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNodoNoEncontrado
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *nodoRepositorio) ObtenerPorNodeKey(ctx context.Context, nodeKey string) (*domain.Nodo, error) {
	var model NodoModel
	err := r.db.WithContext(ctx).Where("node_key = ?", nodeKey).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNodoNoEncontrado
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *nodoRepositorio) Listar(
	ctx context.Context,
	especificacion domain.EspecificacionNodo,
	paginacion shared.Paginacion,
) ([]domain.Nodo, error) {
	query := r.db.WithContext(ctx).Model(&NodoModel{})

	mapeoColumnas := map[string]string{
		"nodeKey":  "node_key",
		"fincaID":  "finca_id",
		"loteID":   "lote_id",
		"tenantID": "tenant_id",
		"estado":   "estado",
		"nombre":   "nombre",
	}

	for _, filtro := range especificacion.Filtros {
		if !domain.ColumnasPermitidasNodos[filtro.Campo] {
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
		if !domain.ColumnasPermitidasNodos[ord.Campo] {
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

	var models []NodoModel
	result := query.Offset(offset).Limit(tamano).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	nodos := make([]domain.Nodo, len(models))
	for i, m := range models {
		nodos[i] = *m.ToDomain()
	}
	return nodos, nil
}

func (r *nodoRepositorio) Actualizar(ctx context.Context, nodo *domain.Nodo) error {
	model := FromDomainNodo(nodo)
	return r.db.WithContext(ctx).Model(&NodoModel{}).Where("id = ?", model.ID).Updates(model).Error
}

func (r *nodoRepositorio) Eliminar(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&NodoModel{}, "id = ?", id).Error
}

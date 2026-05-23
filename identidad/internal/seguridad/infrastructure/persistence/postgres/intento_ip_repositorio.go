package postgres

import (
	"context"
	"errors"
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"gorm.io/gorm"
)

type intentoIPRepositorio struct {
	db *gorm.DB
}

// NewIntentoIPRepositorio crea una nueva instancia del repositorio de intentos por IP.
func NewIntentoIPRepositorio(db *gorm.DB) seguridad_domain.IntentoIPRepositorio {
	return &intentoIPRepositorio{db: db}
}

// ObtenerPorIP retorna el registro de intentos más reciente para una IP.
func (r *intentoIPRepositorio) ObtenerPorIP(ctx context.Context, ip string) (*seguridad_domain.IntentoPorIP, error) {
	var model IntentoIPModel
	result := r.db.WithContext(ctx).
		Where("ip = ?", ip).
		Order("ventana_inicio DESC").
		First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("registro no encontrado")
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

// Crear persiste un nuevo registro de intento por IP.
func (r *intentoIPRepositorio) Crear(ctx context.Context, i *seguridad_domain.IntentoPorIP) (*seguridad_domain.IntentoPorIP, error) {
	model := IntentoIPFromDomain(i)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

// Actualizar actualiza un registro existente de intentos por IP.
func (r *intentoIPRepositorio) Actualizar(ctx context.Context, i *seguridad_domain.IntentoPorIP) (*seguridad_domain.IntentoPorIP, error) {
	result := r.db.WithContext(ctx).Exec(
		`UPDATE intentos_por_ip SET
			contador = ?,
			bloqueado_hasta = ?,
			fecha_actualizacion = ?
		WHERE id = ?`,
		i.Contador(),
		i.BloqueadoHasta(),
		time.Now(),
		i.ID(),
	)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.ObtenerPorIP(ctx, i.IP())
}

// Listar retorna registros de intentos por IP filtrados y paginados.
func (r *intentoIPRepositorio) Listar(ctx context.Context, spec seguridad_domain.EspecificacionIntentoIP, pag shareddomain.Paginacion) ([]*seguridad_domain.IntentoPorIP, error) {
	query := r.db.WithContext(ctx).Model(&IntentoIPModel{})

	mapeoColumnas := map[string]string{
		"ip":             "ip",
		"contador":       "contador",
		"ventanaInicio":  "ventana_inicio",
		"bloqueadoHasta": "bloqueado_hasta",
	}

	for _, filtro := range spec.ListaFiltros {
		if !seguridad_domain.ColumnasPermitidasIntentoIP[filtro.Campo] {
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
		case ">":
			query = query.Where(columnaDB+" > ?", filtro.Valor)
		case "<":
			query = query.Where(columnaDB+" < ?", filtro.Valor)
		case ">=":
			query = query.Where(columnaDB+" >= ?", filtro.Valor)
		case "<=":
			query = query.Where(columnaDB+" <= ?", filtro.Valor)
		}
	}

	for _, ord := range pag.Ordenaciones {
		if !seguridad_domain.ColumnasPermitidasIntentoIP[ord.Campo] {
			continue
		}

		columnaDB, ok := mapeoColumnas[ord.Campo]
		if !ok {
			continue
		}

		orden := "ASC"
		if ord.Tipo == shareddomain.DESC {
			orden = "DESC"
		}
		query = query.Order(columnaDB + " " + orden)
	}

	offset := (pag.Pagina - 1) * pag.TamanoPagina
	if pag.Pagina < 1 {
		offset = 0
	}
	if pag.TamanoPagina < 1 {
		pag.TamanoPagina = 10
	}

	var models []IntentoIPModel
	result := query.Offset(offset).Limit(pag.TamanoPagina).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	intentos := make([]*seguridad_domain.IntentoPorIP, len(models))
	for i, m := range models {
		intentos[i] = m.ToDomain()
	}
	return intentos, nil
}

// EliminarExpirados elimina registros cuya ventana de tiempo ya expiró.
func (r *intentoIPRepositorio) EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error {
	limite := ahora.Add(-ventana)
	result := r.db.WithContext(ctx).
		Where("ventana_inicio < ? AND bloqueado_hasta < ?", limite, ahora).
		Delete(&IntentoIPModel{})
	return result.Error
}
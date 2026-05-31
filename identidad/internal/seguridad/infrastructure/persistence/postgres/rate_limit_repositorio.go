package postgres

import (
	"context"
	"errors"
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"gorm.io/gorm"
)

type rateLimitRepositorio struct {
	db *gorm.DB
}

// NewRateLimitRepositorio crea una nueva instancia del repositorio de rate limiting.
// Usa su propia tabla (rate_limit_ip) para no interferir con el contador de
// intentos fallidos del servicio de bloqueo por IP.
func NewRateLimitRepositorio(db *gorm.DB) seguridad_domain.IntentoIPRepositorio {
	return &rateLimitRepositorio{db: db}
}

func (r *rateLimitRepositorio) ObtenerPorIP(ctx context.Context, ip string) (*seguridad_domain.IntentoPorIP, error) {
	var model RateLimitIPModel
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

func (r *rateLimitRepositorio) Crear(ctx context.Context, i *seguridad_domain.IntentoPorIP) (*seguridad_domain.IntentoPorIP, error) {
	model := RateLimitIPFromDomain(i)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

func (r *rateLimitRepositorio) Actualizar(ctx context.Context, i *seguridad_domain.IntentoPorIP) (*seguridad_domain.IntentoPorIP, error) {
	result := r.db.WithContext(ctx).Exec(
		`UPDATE rate_limit_ip SET
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

func (r *rateLimitRepositorio) Listar(ctx context.Context, spec seguridad_domain.EspecificacionIntentoIP, pag shareddomain.Paginacion) ([]*seguridad_domain.IntentoPorIP, error) {
	query := r.db.WithContext(ctx).Model(&RateLimitIPModel{})

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

	var models []RateLimitIPModel
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

func (r *rateLimitRepositorio) EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error {
	limite := ahora.Add(-ventana)
	result := r.db.WithContext(ctx).
		Where("ventana_inicio < ? AND bloqueado_hasta < ?", limite, ahora).
		Delete(&RateLimitIPModel{})
	return result.Error
}

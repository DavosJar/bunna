package postgres

import (
	"context"
	"errors"
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
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

func (r *rateLimitRepositorio) EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error {
	limite := ahora.Add(-ventana)
	result := r.db.WithContext(ctx).
		Where("ventana_inicio < ? AND bloqueado_hasta < ?", limite, ahora).
		Delete(&RateLimitIPModel{})
	return result.Error
}

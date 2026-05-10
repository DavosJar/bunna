package postgres

import (
	"context"
	"errors"
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
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

// EliminarExpirados elimina registros cuya ventana de tiempo ya expiró.
func (r *intentoIPRepositorio) EliminarExpirados(ctx context.Context, ahora time.Time, ventana time.Duration) error {
	limite := ahora.Add(-ventana)
	result := r.db.WithContext(ctx).
		Where("ventana_inicio < ? AND bloqueado_hasta < ?", limite, ahora).
		Delete(&IntentoIPModel{})
	return result.Error
}
package postgres

import (
	"context"
	"errors"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
	"gorm.io/gorm"
)

type invitacionRepositorio struct {
	db *gorm.DB
}

func NewInvitacionRepositorio(db *gorm.DB) invitaciones.InvitacionRepositorio {
	return &invitacionRepositorio{db: db}
}

func (r *invitacionRepositorio) Crear(ctx context.Context, invitacion *invitaciones.Invitacion) error {
	model := InvitacionFromDomain(invitacion)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *invitacionRepositorio) ObtenerPorTokenHash(ctx context.Context, tokenHash string) (*invitaciones.Invitacion, error) {
	var model InvitacionModel
	result := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, invitaciones.ErrNoEncontrada
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

func (r *invitacionRepositorio) MarcarAceptada(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Model(&InvitacionModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"aceptada":    true,
			"aceptada_at": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return invitaciones.ErrNoEncontrada
	}
	return nil
}

func (r *invitacionRepositorio) ObtenerPorID(ctx context.Context, id string) (*invitaciones.Invitacion, error) {
	var model InvitacionModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, invitaciones.ErrNoEncontrada
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

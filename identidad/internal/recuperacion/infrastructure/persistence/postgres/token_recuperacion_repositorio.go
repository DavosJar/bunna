package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	"gorm.io/gorm"
)

type tokenRecuperacionRepositorio struct {
	db *gorm.DB
}

func NewTokenRecuperacionRepositorio(db *gorm.DB) recuperacion.TokenRecuperacionRepositorio {
	return &tokenRecuperacionRepositorio{db: db}
}

func (r *tokenRecuperacionRepositorio) Crear(ctx context.Context, token *recuperacion.TokenRecuperacion) error {
	model := TokenRecuperacionFromDomain(token)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *tokenRecuperacionRepositorio) ObtenerPorHash(ctx context.Context, hash string) (*recuperacion.TokenRecuperacion, error) {
	var model TokenRecuperacionModel
	err := r.db.WithContext(ctx).First(&model, "token_hash = ?", hash).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, recuperacion.ErrEnlaceInvalido
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *tokenRecuperacionRepositorio) Actualizar(ctx context.Context, token *recuperacion.TokenRecuperacion) error {
	model := TokenRecuperacionFromDomain(token)
	return r.db.WithContext(ctx).Save(model).Error
}

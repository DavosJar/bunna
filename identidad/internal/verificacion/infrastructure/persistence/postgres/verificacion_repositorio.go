package postgres

import (
	"context"
	"errors"
	"time"

	dominio "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
	"gorm.io/gorm"
)

type verificacionRepositorio struct {
	db *gorm.DB
}

func NewVerificacionRepositorio(db *gorm.DB) dominio.VerificacionRepositorio {
	return &verificacionRepositorio{db: db}
}

func (r *verificacionRepositorio) ObtenerPorHashToken(ctx context.Context, hash string) (*dominio.UsuarioVerificacion, error) {
	var model VerificacionUsuarioModel
	result := r.db.WithContext(ctx).
		Where("verificacion_token_hash = ?", hash).
		First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, dominio.ErrEnlaceInvalido
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

func (r *verificacionRepositorio) ActualizarPrueba(ctx context.Context, usuarioID string, prueba dominio.PruebaVerificacion) error {
	ahora := time.Now()
	return r.db.WithContext(ctx).Model(&VerificacionUsuarioModel{}).
		Where("id = ?", usuarioID).
		Updates(map[string]interface{}{
			"verificacion_token_hash": prueba.SecretoHash(),
			"verificacion_expiracion": prueba.ExpiraEn(),
			"contador_reenvios":       gorm.Expr("contador_reenvios + 1"),
			"ultimo_reenvio":          ahora,
		}).Error
}

func (r *verificacionRepositorio) MarcarVerificado(ctx context.Context, usuarioID string) error {
	return r.db.WithContext(ctx).Model(&VerificacionUsuarioModel{}).
		Where("id = ?", usuarioID).
		Updates(map[string]interface{}{
			"estado_verificacion_correo": "VERIFICADO",
			"estado":                     "ACTIVO",
		}).Error
}

func (r *verificacionRepositorio) ObtenerPorID(ctx context.Context, usuarioID string) (*dominio.UsuarioVerificacion, error) {
	var model VerificacionUsuarioModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", usuarioID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, dominio.ErrUsuarioNoEncontrado
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
	"gorm.io/gorm"
)

var _ recuperacion.UsuarioRecuperacionRepositorio = (*usuarioRecuperacionRepositorio)(nil)

type usuarioRecuperacionModel struct {
	ID     string `gorm:"column:id;type:uuid;primaryKey"`
	Nombre string `gorm:"column:nombre"`
	Correo string `gorm:"column:correo"`
}

func (usuarioRecuperacionModel) TableName() string {
	return "usuarios"
}

type usuarioRecuperacionRepositorio struct {
	db *gorm.DB
}

func NewUsuarioRecuperacionRepositorio(db *gorm.DB) recuperacion.UsuarioRecuperacionRepositorio {
	return &usuarioRecuperacionRepositorio{db: db}
}

func (r *usuarioRecuperacionRepositorio) ObtenerPorCorreo(ctx context.Context, correo string) (*recuperacion.UsuarioRecuperacion, error) {
	var model usuarioRecuperacionModel
	result := r.db.WithContext(ctx).
		Where("correo = ?", correo).
		First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, recuperacion.ErrUsuarioNoEncontrado
		}
		return nil, result.Error
	}
	return &recuperacion.UsuarioRecuperacion{
		ID:     model.ID,
		Nombre: model.Nombre,
		Correo: model.Correo,
	}, nil
}

func (r *usuarioRecuperacionRepositorio) ActualizarPassword(ctx context.Context, usuarioID, nuevoHash string) error {
	result := r.db.WithContext(ctx).Exec(
		`UPDATE credenciales_usuarios SET
			password_hash = ?,
			intentos_fallidos = 0,
			bloqueado_hasta = NULL,
			fecha_actualizacion = ?
		WHERE usuario_id = ?`,
		nuevoHash,
		time.Now(),
		usuarioID,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return recuperacion.ErrUsuarioNoEncontrado
	}
	return nil
}

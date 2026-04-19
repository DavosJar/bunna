package postgres

import (
	"fmt"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"github.com/google/uuid"
)

type UsuarioModel struct {
	ID                       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Nombre                   string    `gorm:"column:nombre"`
	Apellido                 string    `gorm:"column:apellido"`
	Correo                   string    `gorm:"column:correo"`
	Telefono                 string    `gorm:"column:telefono"`
	Estado                   string    `gorm:"column:estado"`
	EstadoVerificacionCorreo string    `gorm:"column:estado_verificacion_correo"`
	FechaCreacion            time.Time `gorm:"column:fecha_creacion;autoCreateTime:milli"`
	FechaActualizacion       time.Time `gorm:"column:fecha_actualizacion;autoUpdateTime:milli"`
}

func (UsuarioModel) TableName() string {
	return "usuarios"
}

func (m *UsuarioModel) ToDomain() *usuario.Usuario {
	return usuario.NewUsuarioFromPersistence(
		m.ID.String(),
		m.Nombre,
		m.Apellido,
		m.Correo,
		m.Telefono,
		usuario.EstadoUsuario(m.Estado),
		usuario.EstadoVerificacionCorreo(m.EstadoVerificacionCorreo),
		m.FechaCreacion,
		m.FechaActualizacion,
	)
}

func FromDomain(u *usuario.Usuario) (*UsuarioModel, error) {
	uuidValue, err := uuid.Parse(u.ID())
	if err != nil {
		return nil, fmt.Errorf("ID de usuario inválido '%s': %w", u.ID(), err)
	}
	return &UsuarioModel{
		ID:                       uuidValue,
		Nombre:                   u.Nombre(),
		Apellido:                 u.Apellido(),
		Correo:                   u.Correo(),
		Telefono:                 u.Telefono(),
		Estado:                   string(u.Estado()),
		EstadoVerificacionCorreo: string(u.EstadoVerificacionCorreo()),
		FechaCreacion:            u.FechaCreacion(),
		FechaActualizacion:       u.FechaActualizacion(),
	}, nil
}

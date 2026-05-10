// Package postgres implementa los repositorios de sesiones usando GORM y PostgreSQL.
package postgres

import (
	"time"

	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

// SesionModel es el modelo GORM para la tabla sesiones.
type SesionModel struct {
	ID                    string    `gorm:"column:id;primaryKey;type:varchar(36)"`
	UsuarioID             string    `gorm:"column:usuario_id;type:varchar(36);not null"`
	AccessTokenHash       string    `gorm:"column:access_token_hash;type:varchar(64);not null"`
	RefreshTokenHash      string    `gorm:"column:refresh_token_hash;type:varchar(64);not null"`
	Estado                string    `gorm:"column:estado;type:varchar(20);not null;default:'ACTIVA'"`
	IPOrigen              string    `gorm:"column:ip_origen;type:varchar(45)"`
	FechaCreacion         time.Time `gorm:"column:fecha_creacion;not null"`
	FechaActualizacion    time.Time `gorm:"column:fecha_actualizacion;not null"`
	FechaExpiracionAccess  time.Time `gorm:"column:fecha_expiracion_access;not null"`
	FechaExpiracionRefresh time.Time `gorm:"column:fecha_expiracion_refresh;not null"`
	UltimaActividad       time.Time `gorm:"column:ultima_actividad;not null"`
	ContadorRefrescos     int       `gorm:"column:contador_refrescos;type:int;not null;default:0"`
}

// TableName retorna el nombre de la tabla en PostgreSQL.
func (SesionModel) TableName() string {
	return "sesiones"
}

// ToDomain convierte el modelo GORM a la entidad de dominio.
func (m *SesionModel) ToDomain() *sesiones_domain.Sesion {
	return sesiones_domain.NuevaSesionDesdeBD(
		m.ID,
		m.UsuarioID,
		m.AccessTokenHash,
		m.RefreshTokenHash,
		sesiones_domain.EstadoSesion(m.Estado),
		m.IPOrigen,
		m.FechaCreacion,
		m.FechaActualizacion,
		m.FechaExpiracionAccess,
		m.FechaExpiracionRefresh,
		m.UltimaActividad,
		m.ContadorRefrescos,
	)
}

// SesionFromDomain convierte la entidad de dominio al modelo GORM.
func SesionFromDomain(s *sesiones_domain.Sesion) *SesionModel {
	return &SesionModel{
		ID:                    s.ID(),
		UsuarioID:             s.UsuarioID(),
		AccessTokenHash:       s.AccessTokenHash(),
		RefreshTokenHash:      s.RefreshTokenHash(),
		Estado:                string(s.Estado()),
		IPOrigen:              s.IPOrigen(),
		FechaCreacion:         s.FechaCreacion(),
		FechaActualizacion:    s.FechaActualizacion(),
		FechaExpiracionAccess:  s.FechaExpiracionAccess(),
		FechaExpiracionRefresh: s.FechaExpiracionRefresh(),
		UltimaActividad:       s.UltimaActividad(),
		ContadorRefrescos:     s.ContadorRefrescos(),
	}
}
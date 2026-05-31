package postgres

import (
	"time"

	dominio "github.com/davosjar/bunna/services/identidad/internal/verificacion/domain"
)

type VerificacionUsuarioModel struct {
	ID                       string     `gorm:"column:id;type:uuid;primaryKey"`
	Nombre                   string     `gorm:"column:nombre"`
	Correo                   string     `gorm:"column:correo"`
	EstadoVerificacionCorreo string     `gorm:"column:estado_verificacion_correo"`
	VerificacionTokenHash    string     `gorm:"column:verificacion_token_hash"`
	VerificacionExpiracion   time.Time  `gorm:"column:verificacion_expiracion"`
	ContadorReenvios         int        `gorm:"column:contador_reenvios;default:0"`
	UltimoReenvio            *time.Time `gorm:"column:ultimo_reenvio"`
}

func (VerificacionUsuarioModel) TableName() string {
	return "usuarios"
}

func (m *VerificacionUsuarioModel) ToDomain() *dominio.UsuarioVerificacion {
	prueba := dominio.NuevaPruebaVerificacionDesdeBD(m.VerificacionTokenHash, m.VerificacionExpiracion)
	return &dominio.UsuarioVerificacion{
		ID:                 m.ID,
		Nombre:             m.Nombre,
		Correo:             m.Correo,
		EstadoVerificacion: m.EstadoVerificacionCorreo,
		PruebaVerificacion: prueba,
		ContadorReenvios:   m.ContadorReenvios,
		UltimoReenvio:      m.UltimoReenvio,
	}
}

package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

type CredencialesModel struct {
	UsuarioID        string    `gorm:"column:usuario_id;primaryKey;type:varchar(36)"`
	PasswordHash     string    `gorm:"column:password_hash;type:varchar(255)"`
	Activo           bool      `gorm:"column:activo;type:boolean;default:true"`
	CorreoVerificado bool      `gorm:"column:correo_verificado;type:boolean;default:false"`
	IntentosFallidos int       `gorm:"column:intentos_fallidos;type:int;default:0"`
	BloqueadoHasta   time.Time `gorm:"column:bloqueado_hasta;type:timestamptz;default:null"`
}

func (CredencialesModel) TableName() string {
	return "credenciales_usuarios"
}

func (m *CredencialesModel) ToDomain() *domain.CredencialesUsuario {
	return domain.NuevaCredencialesUsuarioDesdeBD(
		m.UsuarioID,
		m.PasswordHash,
		m.Activo,
		m.CorreoVerificado,
		m.IntentosFallidos,
		m.BloqueadoHasta,
	)
}

func CredencialesFromDomain(c *domain.CredencialesUsuario) (*CredencialesModel, error) {
	// Usar los métodos getter públicos para acceder a los valores
	return &CredencialesModel{
		UsuarioID:        c.UsuarioID(),
		PasswordHash:     c.PasswordHash(),
		Activo:           c.Activo(),
		CorreoVerificado: c.CorreoVerificado(),
		IntentosFallidos: c.IntentosFallidos(),
		BloqueadoHasta:   c.BloqueadoHasta(),
	}, nil
}

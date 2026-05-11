package postgres

import (
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

// RateLimitIPModel es el modelo GORM para la tabla rate_limit_ip.
// Separado de IntentoIPModel para evitar conflictos entre el rate limiter
// (cuenta requests totales) y el bloqueo por IP (cuenta intentos fallidos).
type RateLimitIPModel struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(36)"`
	IP             string    `gorm:"column:ip;type:varchar(45);not null;index"`
	Contador       int       `gorm:"column:contador;type:int;not null;default:1"`
	VentanaInicio  time.Time `gorm:"column:ventana_inicio;not null"`
	BloqueadoHasta time.Time `gorm:"column:bloqueado_hasta"`
	FechaCreacion  time.Time `gorm:"column:fecha_creacion;not null"`
	FechaActualizacion time.Time `gorm:"column:fecha_actualizacion;not null"`
}

func (RateLimitIPModel) TableName() string {
	return "rate_limit_ip"
}

func (m *RateLimitIPModel) ToDomain() *seguridad_domain.IntentoPorIP {
	return seguridad_domain.NuevoIntentoPorIPDesdeBD(
		m.ID,
		m.IP,
		m.Contador,
		m.VentanaInicio,
		m.BloqueadoHasta,
	)
}

func RateLimitIPFromDomain(i *seguridad_domain.IntentoPorIP) *RateLimitIPModel {
	ahora := time.Now()
	return &RateLimitIPModel{
		ID:             i.ID(),
		IP:             i.IP(),
		Contador:       i.Contador(),
		VentanaInicio:  i.VentanaInicio(),
		BloqueadoHasta: i.BloqueadoHasta(),
		FechaCreacion:  ahora,
		FechaActualizacion: ahora,
	}
}

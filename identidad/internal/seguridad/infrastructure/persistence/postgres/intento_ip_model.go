package postgres

import (
	"time"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

// IntentoIPModel es el modelo GORM para la tabla intentos_por_ip.
type IntentoIPModel struct {
	ID             string    `gorm:"column:id;primaryKey;type:varchar(36)"`
	IP             string    `gorm:"column:ip;type:varchar(45);not null;index"`
	Contador       int       `gorm:"column:contador;type:int;not null;default:1"`
	VentanaInicio  time.Time `gorm:"column:ventana_inicio;not null"`
	BloqueadoHasta time.Time `gorm:"column:bloqueado_hasta"`
	FechaCreacion  time.Time `gorm:"column:fecha_creacion;not null"`
	FechaActualizacion time.Time `gorm:"column:fecha_actualizacion;not null"`
}

// TableName retorna el nombre de la tabla en PostgreSQL.
func (IntentoIPModel) TableName() string {
	return "intentos_por_ip"
}

// ToDomain convierte el modelo GORM a la entidad de dominio.
func (m *IntentoIPModel) ToDomain() *seguridad_domain.IntentoPorIP {
	return seguridad_domain.NuevoIntentoPorIPDesdeBD(
		m.ID,
		m.IP,
		m.Contador,
		m.VentanaInicio,
		m.BloqueadoHasta,
	)
}

// IntentoIPFromDomain convierte la entidad de dominio al modelo GORM.
func IntentoIPFromDomain(i *seguridad_domain.IntentoPorIP) *IntentoIPModel {
	ahora := time.Now()
	return &IntentoIPModel{
		ID:             i.ID(),
		IP:             i.IP(),
		Contador:       i.Contador(),
		VentanaInicio:  i.VentanaInicio(),
		BloqueadoHasta: i.BloqueadoHasta(),
		FechaCreacion:  ahora,
		FechaActualizacion: ahora,
	}
}
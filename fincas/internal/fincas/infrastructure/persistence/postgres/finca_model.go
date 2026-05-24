package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
)

type FincaModel struct {
	ID          string    `gorm:"type:varchar(36);primaryKey;column:id"`
	Nombre      string    `gorm:"column:nombre"`
	Ubicacion   string    `gorm:"column:ubicacion"`
	Descripcion string    `gorm:"column:descripcion"`
	UsuarioID   string    `gorm:"column:usuario_id;index"`
	TenantID    *string   `gorm:"column:tenant_id;index"`
	Estado      string    `gorm:"column:estado;default:ACTIVA"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (FincaModel) TableName() string { return "fincas" }

func (m *FincaModel) ToDomain() *domain.Finca {
	return domain.NewFincaFromPersistence(
		m.ID, m.Nombre, m.Ubicacion, m.Descripcion, m.UsuarioID,
		m.TenantID, domain.EstadoFinca(m.Estado), m.CreatedAt, m.UpdatedAt,
	)
}

func FromDomainFinca(f *domain.Finca) *FincaModel {
	return &FincaModel{
		ID:          f.ID(),
		Nombre:      f.Nombre(),
		Ubicacion:   f.Ubicacion(),
		Descripcion: f.Descripcion(),
		UsuarioID:   f.UsuarioID(),
		TenantID:    f.TenantID(),
		Estado:      string(f.Estado()),
		CreatedAt:   f.CreatedAt(),
		UpdatedAt:   f.UpdatedAt(),
	}
}

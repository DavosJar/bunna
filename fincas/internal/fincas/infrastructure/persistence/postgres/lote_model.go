package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
)

type LoteModel struct {
	ID          string    `gorm:"type:varchar(36);primaryKey;column:id"`
	FincaID     string    `gorm:"column:finca_id;index"`
	Nombre      string    `gorm:"column:nombre"`
	Area        float64   `gorm:"column:area"`
	Descripcion string    `gorm:"column:descripcion"`
	Estado      string    `gorm:"column:estado;default:ACTIVO"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (LoteModel) TableName() string { return "lotes" }

func (m *LoteModel) ToDomain() *domain.Lote {
	return domain.NewLoteFromPersistence(
		m.ID, m.FincaID, m.Nombre, m.Area, m.Descripcion,
		domain.EstadoLote(m.Estado), m.CreatedAt, m.UpdatedAt,
	)
}

func FromDomainLote(l *domain.Lote) *LoteModel {
	return &LoteModel{
		ID:          l.ID(),
		FincaID:     l.FincaID(),
		Nombre:      l.Nombre(),
		Area:        l.Area(),
		Descripcion: l.Descripcion(),
		Estado:      string(l.Estado()),
		CreatedAt:   l.CreatedAt(),
		UpdatedAt:   l.UpdatedAt(),
	}
}

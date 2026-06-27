package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/nodos/domain"
)

type NodoModel struct {
	ID            string    `gorm:"type:varchar(36);primaryKey;column:id"`
	FincaID       string    `gorm:"column:finca_id;type:varchar(36);index"`
	LoteID        *string   `gorm:"column:lote_id;type:varchar(36);index"`
	TenantID      string    `gorm:"column:tenant_id;type:varchar(36);index"`
	Nombre        string    `gorm:"column:nombre;type:varchar(100)"`
	NodeKey       string    `gorm:"column:node_key;type:varchar(100);uniqueIndex"`
	Estado        string    `gorm:"column:estado;type:varchar(20);default:ACTIVO"`
	CreadoEn      time.Time `gorm:"column:created_at;autoCreateTime"`
	ActualizadoEn time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (NodoModel) TableName() string { return "nodos" }

func (m *NodoModel) ToDomain() *domain.Nodo {
	return domain.NewNodoFromPersistence(
		m.ID,
		m.TenantID,
		m.FincaID,
		m.NodeKey,
		m.LoteID,
		m.Nombre,
		domain.EstadoNodo(m.Estado),
		m.CreadoEn,
		m.ActualizadoEn,
	)
}

func FromDomainNodo(n *domain.Nodo) *NodoModel {
	return &NodoModel{
		ID:            n.ID(),
		FincaID:       n.FincaID(),
		LoteID:        n.LoteID(),
		TenantID:      n.TenantID(),
		Nombre:        n.Nombre(),
		NodeKey:       n.NodeKey(),
		Estado:        string(n.Estado()),
		CreadoEn:      n.CreadoEn(),
		ActualizadoEn: n.ActualizadoEn(),
	}
}

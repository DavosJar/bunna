package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

type TenantModel struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey;column:id"`
	Nombre             string    `gorm:"column:nombre"`
	Slug               string    `gorm:"column:slug;uniqueIndex"`
	Activo             bool      `gorm:"column:activo;default:true"`
	FechaCreacion      time.Time `gorm:"column:created_at;autoCreateTime"`
	FechaActualizacion time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (TenantModel) TableName() string {
	return "tenants"
}

func (m *TenantModel) ToDomain() *tenant.Tenant {
	return tenant.NuevoTenantDesdeBD(
		m.ID,
		m.Nombre,
		m.Slug,
		m.Activo,
		m.FechaCreacion,
		m.FechaActualizacion,
	)
}

func TenantFromDomain(t *tenant.Tenant) *TenantModel {
	return &TenantModel{
		ID:                 t.ID(),
		Nombre:             t.Nombre(),
		Slug:               t.Slug(),
		Activo:             t.EstaActivo(),
		FechaCreacion:      t.FechaCreacion(),
		FechaActualizacion: t.FechaActualizacion(),
	}
}

// MembresiaModel representa la tabla usuario_tenants
type MembresiaModel struct {
	UsuarioID     string    `gorm:"type:varchar(36);primaryKey;column:usuario_id"`
	TenantID      string    `gorm:"type:varchar(36);primaryKey;column:tenant_id"`
	FechaCreacion time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (MembresiaModel) TableName() string {
	return "usuario_tenants"
}

func (m *MembresiaModel) ToDomain() *tenant.Membresia {
	return tenant.NuevaMembresiaDesdeBD(
		m.UsuarioID,
		m.TenantID,
		m.FechaCreacion,
	)
}

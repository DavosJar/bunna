package postgres

import (
	"time"

	invitaciones "github.com/davosjar/bunna/services/identidad/internal/invitaciones/domain"
)

type InvitacionModel struct {
	ID              string     `gorm:"type:varchar(36);primaryKey;column:id"`
	TenantID        string     `gorm:"type:varchar(36);column:tenant_id;not null"`
	RolID           string     `gorm:"type:varchar(36);column:rol_id;not null"`
	Email           string     `gorm:"type:varchar(254);column:email;not null"`
	Nombre          string     `gorm:"type:varchar(200);column:nombre;default:''"`
	TokenHash       string     `gorm:"type:varchar(64);column:token_hash;uniqueIndex;not null"`
	Expiracion      time.Time  `gorm:"column:expiracion;not null"`
	Aceptada        bool       `gorm:"column:aceptada;default:false"`
	FechaCreacion   time.Time  `gorm:"column:created_at;autoCreateTime"`
	FechaAceptacion *time.Time `gorm:"column:aceptada_at"`
}

func (InvitacionModel) TableName() string {
	return "invitaciones"
}

func (m *InvitacionModel) ToDomain() *invitaciones.Invitacion {
	return invitaciones.NuevaInvitacionDesdeBD(
		m.ID, m.TenantID, m.RolID, m.Email, m.Nombre,
		m.TokenHash, m.Expiracion, m.Aceptada, m.FechaCreacion, m.FechaAceptacion,
	)
}

func InvitacionFromDomain(i *invitaciones.Invitacion) *InvitacionModel {
	return &InvitacionModel{
		ID:              i.ID(),
		TenantID:        i.TenantID(),
		RolID:           i.RolID(),
		Email:           i.Email(),
		Nombre:          i.Nombre(),
		TokenHash:       i.TokenHash(),
		Expiracion:      i.Expiracion(),
		Aceptada:        i.EstaAceptada(),
		FechaCreacion:   i.FechaCreacion(),
		FechaAceptacion: i.FechaAceptacion(),
	}
}

package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
)

// TokenRecuperacionModel representa la tabla tokens_recuperacion
type TokenRecuperacionModel struct {
	ID        string     `gorm:"type:varchar(36);primaryKey;column:id"`
	UsuarioID string     `gorm:"type:varchar(36);column:usuario_id"`
	TokenHash string     `gorm:"column:token_hash"`
	ExpiraEn  time.Time  `gorm:"column:expira_en"`
	Usado     bool       `gorm:"column:usado;default:false"`
	CreadoEn  time.Time  `gorm:"column:creado_en;autoCreateTime"`
	UsadoEn   *time.Time `gorm:"column:usado_en"`
}

func (TokenRecuperacionModel) TableName() string {
	return "tokens_recuperacion"
}

func (m *TokenRecuperacionModel) ToDomain() *recuperacion.TokenRecuperacion {
	return recuperacion.NuevoTokenRecuperacionDesdeBD(
		m.ID,
		m.UsuarioID,
		m.TokenHash,
		m.ExpiraEn,
		m.Usado,
		m.CreadoEn,
		m.UsadoEn,
	)
}

func TokenRecuperacionFromDomain(t *recuperacion.TokenRecuperacion) *TokenRecuperacionModel {
	return &TokenRecuperacionModel{
		ID:        t.ID(),
		UsuarioID: t.UsuarioID(),
		TokenHash: t.TokenHash(),
		ExpiraEn:  t.ExpiraEn(),
		Usado:     t.Usado(),
		CreadoEn:  t.CreadoEn(),
		UsadoEn:   t.UsadoEn(),
	}
}

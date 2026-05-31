package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
)

// CandidatoModel mapea a la tabla "candidatos_reentrenamiento".
// Es un modelo de persistencia puro — no tiene entidad de dominio equivalente
// como aggregate raíz, pero sí se apoya en la entidad CandidatoReentrenamiento
// para la conversión.
type CandidatoModel struct {
	ID                    string    `gorm:"type:varchar(36);primaryKey;column:id"`
	DiagnosticoID         string    `gorm:"column:diagnostico_id;type:varchar(36);uniqueIndex"`
	ImageURL              string    `gorm:"column:image_url;type:text"`
	TieneClorosis         bool      `gorm:"column:tiene_clorosis"`
	Confianza             float64   `gorm:"column:confianza;type:decimal(5,4)"`
	Motivo                *string   `gorm:"column:motivo;type:text"`
	RechazadoPorUsuarioID string    `gorm:"column:rechazado_por_usuario_id;type:varchar(36)"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (CandidatoModel) TableName() string { return "candidatos_reentrenamiento" }

// ToDomain reconstruye la entidad CandidatoReentrenamiento desde el modelo de persistencia.
func (m *CandidatoModel) ToDomain() *domain.CandidatoReentrenamiento {
	motivo := ""
	if m.Motivo != nil {
		motivo = *m.Motivo
	}
	return domain.NewCandidatoReentrenamientoFromStorage(
		m.ID,
		m.DiagnosticoID,
		m.ImageURL,
		m.TieneClorosis,
		m.Confianza,
		motivo,
		m.RechazadoPorUsuarioID,
		m.CreatedAt,
	)
}

// FromDomainCandidato extrae los campos de CandidatoReentrenamiento al modelo plano.
func FromDomainCandidato(c *domain.CandidatoReentrenamiento) *CandidatoModel {
	var motivo *string
	if c.Motivo() != "" {
		m := c.Motivo()
		motivo = &m
	}
	return &CandidatoModel{
		ID:                    c.ID(),
		DiagnosticoID:         c.DiagnosticoID(),
		ImageURL:              c.ImageURL(),
		TieneClorosis:         c.TieneClorosis(),
		Confianza:             c.Confianza(),
		Motivo:                motivo,
		RechazadoPorUsuarioID: c.RechazadoPorUsuarioID(),
		CreatedAt:             c.CreatedAt(),
	}
}

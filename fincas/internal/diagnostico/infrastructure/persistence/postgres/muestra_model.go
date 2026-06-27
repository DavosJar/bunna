package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
)

// MuestraModel mapea la entidad Muestra a la tabla "muestras".
// El VO Ubicacion se persiste como columnas planas (latitud, longitud).
type MuestraModel struct {
	ID        string    `gorm:"type:varchar(36);primaryKey;column:id"`
	FincaID   string    `gorm:"column:finca_id;type:varchar(36);index"`
	LoteID    string    `gorm:"column:lote_id;type:varchar(36);index"`
	TenantID  string    `gorm:"column:tenant_id;type:varchar(36);index"`
	Latitud   float64   `gorm:"column:latitud;type:double precision"`
	Longitud  float64   `gorm:"column:longitud;type:double precision"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (MuestraModel) TableName() string { return "muestras" }

// ToDomain reconstruye la entidad Muestra con su VO Ubicacion embebido.
func (m *MuestraModel) ToDomain() (*domain.Muestra, error) {
	ubicacion, err := domain.NewUbicacion(m.Latitud, m.Longitud)
	if err != nil {
		return nil, err
	}

	return domain.NewMusetraFromStorage(
		m.ID,
		m.FincaID,
		m.LoteID,
		m.TenantID,
		*ubicacion,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// FromDomainMuestra extrae los campos de la entidad Muestra y su Ubicacion
// al modelo plano para persistencia.
func FromDomainMuestra(mu *domain.Muestra) *MuestraModel {
	u := mu.Ubicacion()
	return &MuestraModel{
		ID:       mu.ID(),
		FincaID:  mu.FincaID(),
		LoteID:   mu.LoteID(),
		TenantID: mu.TenantID(),
		Latitud:  u.Latitud(),
		Longitud: u.Longitud(),
	}
}

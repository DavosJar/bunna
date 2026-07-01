package postgres

import (
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
)

// DiagnosticoModel mapea la entidad Diagnostico a la tabla "diagnosticos".
// El VO ResultadoInferencia se persiste como columnas planas para evitar
// joins y facilitar consultas directas sobre confianza o tieneClorosis.
type DiagnosticoModel struct {
	ID             string    `gorm:"type:varchar(36);primaryKey;column:id"`
	Nombre         string    `gorm:"column:nombre;type:varchar(200)"`
	MuestrasID     string    `gorm:"column:muestras_id;type:varchar(36);index"`
	TenantID       string    `gorm:"column:tenant_id;type:varchar(36);index"`
	Estado         string    `gorm:"column:estado;type:varchar(20);default:PENDIENTE"`
	ImageURL       string    `gorm:"column:image_url;type:text"`
	ImageBase64    string    `gorm:"column:image_data;type:text"`
	TieneClorosis  bool      `gorm:"column:tiene_clorosis"`
	Confianza      float64   `gorm:"column:confianza;type:decimal(5,4)"`
	ProcesadoAt    time.Time `gorm:"column:procesado_at"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DiagnosticoModel) TableName() string { return "diagnosticos" }

// ToDomain reconstruye la entidad de dominio Diagnostico a partir del modelo
// de persistencia, incluyendo el VO ResultadoInferencia embebido.
func (m *DiagnosticoModel) ToDomain() (*domain.Diagnostico, error) {
	resultado, err := domain.NewResultadoInferencia(
		m.ImageURL,
		m.ImageBase64,
		m.TieneClorosis,
		m.Confianza,
		m.ProcesadoAt,
	)
	if err != nil {
		return nil, err
	}

	return domain.NewDiagnosticoFromStorage(
		m.ID,
		m.Nombre,
		m.MuestrasID,
		m.TenantID,
		resultado,
		m.CreatedAt,
		m.UpdatedAt,
		domain.EstadoDiagnostico(m.Estado),
	)
}

// FromDomainDiagnostico extrae todos los campos de la entidad Diagnostico
// y su ResultadoInferencia al modelo plano para persistencia.
func FromDomainDiagnostico(d *domain.Diagnostico) *DiagnosticoModel {
	ri := d.ResultadoInferencia()
	return &DiagnosticoModel{
		ID:            d.ID(),
		Nombre:        d.Nombre(),
		MuestrasID:    d.MuestrasId(),
		TenantID:      d.TenantID(),
		Estado:        string(d.Estado()),
		ImageURL:      ri.ImageUrl(),
		ImageBase64:   ri.ImageBase64(),
		TieneClorosis: ri.TieneClorosis(),
		Confianza:     ri.Confianza(),
		ProcesadoAt:   ri.ProcesadoAt(),
		CreatedAt:     d.CreatedAt(),
		UpdatedAt:     d.UpdatedAt(),
	}
}

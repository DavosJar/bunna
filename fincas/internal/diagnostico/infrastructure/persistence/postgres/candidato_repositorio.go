package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	"gorm.io/gorm"
)

type candidatoRepositorio struct {
	db *gorm.DB
}

func NewCandidatoReentrenamientoRepositorio(db *gorm.DB) domain.CandidatoReentrenamientoRepositorio {
	return &candidatoRepositorio{db: db}
}

func (r *candidatoRepositorio) Crear(ctx context.Context, candidato *domain.CandidatoReentrenamiento) error {
	model := FromDomainCandidato(candidato)
	return r.db.WithContext(ctx).Create(model).Error
}

// ObtenerPorDiagnosticoID busca el candidato asociado a un diagnóstico específico.
// diagnostico_id tiene índice único, así que a lo sumo hay un registro.
func (r *candidatoRepositorio) ObtenerPorDiagnosticoID(ctx context.Context, diagnosticoID string) (*domain.CandidatoReentrenamiento, error) {
	var model CandidatoModel
	err := r.db.WithContext(ctx).Where("diagnostico_id = ?", diagnosticoID).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCandidatoNoEncontrado
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// ListarPendientes retorna los N candidatos más antiguos ordenados por created_at ASC.
func (r *candidatoRepositorio) ListarPendientes(ctx context.Context, limite int) ([]domain.CandidatoReentrenamiento, error) {
	if limite < 1 {
		limite = 10
	}

	var models []CandidatoModel
	err := r.db.WithContext(ctx).
		Order("created_at ASC").
		Limit(limite).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	candidatos := make([]domain.CandidatoReentrenamiento, len(models))
	for i, m := range models {
		candidatos[i] = *m.ToDomain()
	}
	return candidatos, nil
}

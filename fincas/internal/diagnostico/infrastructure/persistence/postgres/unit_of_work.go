package postgres

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"gorm.io/gorm"
)

// UnitOfWorkDiagnostico gestiona transacciones GORM que abarcan los repositorios
// del módulo de diagnósticos: Diagnostico, Muestra y CandidatoReentrenamiento.
type UnitOfWorkDiagnostico struct {
	db             *gorm.DB
	diagnosticoRepo domain.DiagnosticoRepositorio
	muestraRepo     domain.MuestraRepositorio
	candidatoRepo   domain.CandidatoReentrenamientoRepositorio
	generadorID     shared.GeneradorID
}

func NewUnitOfWorkDiagnostico(db *gorm.DB, generadorID shared.GeneradorID) *UnitOfWorkDiagnostico {
	return &UnitOfWorkDiagnostico{
		db:              db,
		diagnosticoRepo: NewDiagnosticoRepositorio(db),
		muestraRepo:     NewMuestraRepositorio(db),
		candidatoRepo:   NewCandidatoReentrenamientoRepositorio(db),
		generadorID:     generadorID,
	}
}

// Transaccional ejecuta fn dentro de una transacción GORM.
// Si fn retorna error, se hace rollback. Si retorna nil, se hace commit.
// Dentro de la transacción, cada repositorio opera con el *gorm.DB transaccional.
func (uw *UnitOfWorkDiagnostico) Transaccional(
	ctx context.Context,
	fn func(application.UnitOfWorkDiagnostico) error,
) error {
	return uw.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		txUow := &UnitOfWorkDiagnostico{
			db:              txDB,
			diagnosticoRepo: NewDiagnosticoRepositorio(txDB),
			muestraRepo:     NewMuestraRepositorio(txDB),
			candidatoRepo:   NewCandidatoReentrenamientoRepositorio(txDB),
			generadorID:     uw.generadorID,
		}
		return fn(txUow)
	})
}

func (uw *UnitOfWorkDiagnostico) DiagnosticoRepo() domain.DiagnosticoRepositorio {
	return uw.diagnosticoRepo
}

func (uw *UnitOfWorkDiagnostico) MuestraRepo() domain.MuestraRepositorio {
	return uw.muestraRepo
}

func (uw *UnitOfWorkDiagnostico) CandidatoRepo() domain.CandidatoReentrenamientoRepositorio {
	return uw.candidatoRepo
}

func (uw *UnitOfWorkDiagnostico) GeneradorID() shared.GeneradorID {
	return uw.generadorID
}

package postgres

import (
	"context"

	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	"gorm.io/gorm"
)

type UnitOfWorkPostgres struct {
	db          *gorm.DB
	fincaRepo   domain.FincaRepositorio
	loteRepo    domain.LoteRepositorio
	generadorID shared.GeneradorID
}

func NewUnitOfWorkPostgres(db *gorm.DB, generadorID shared.GeneradorID) *UnitOfWorkPostgres {
	return &UnitOfWorkPostgres{
		db:          db,
		fincaRepo:   NewFincaRepositorio(db),
		loteRepo:    NewLoteRepositorio(db),
		generadorID: generadorID,
	}
}

func (uw *UnitOfWorkPostgres) Transaccional(
	ctx context.Context,
	fn func(tx *UnitOfWorkPostgres) error,
) error {
	return uw.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		txUow := &UnitOfWorkPostgres{
			db:          txDB,
			fincaRepo:   NewFincaRepositorio(txDB),
			loteRepo:    NewLoteRepositorio(txDB),
			generadorID: uw.generadorID,
		}
		return fn(txUow)
	})
}

func (uw *UnitOfWorkPostgres) FincaRepository() domain.FincaRepositorio {
	return uw.fincaRepo
}

func (uw *UnitOfWorkPostgres) LoteRepository() domain.LoteRepositorio {
	return uw.loteRepo
}

func (uw *UnitOfWorkPostgres) GeneradorID() shared.GeneradorID {
	return uw.generadorID
}

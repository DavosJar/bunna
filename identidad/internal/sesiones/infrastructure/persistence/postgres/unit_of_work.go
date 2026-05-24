package postgres

import (
	"context"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	"gorm.io/gorm"
)

// SesionUnitOfWorkPostgres implementa sesiones_domain.UnitOfWork usando GORM.
type SesionUnitOfWorkPostgres struct {
	db                   *gorm.DB
	sesionRepo           sesiones_domain.SesionRepositorio
	credencialesRepo     seguridad_domain.CredencialesRepositorio
	usuarioRepo          usuario_domain.UsuarioRepositorio
	encriptacionServicio seguridad_domain.EncriptacionServicio
	tokenServicio        sesiones_domain.TokenServicio
	generadorID          shared_domain.GeneradorID
}

// NewSesionUnitOfWork crea una nueva instancia de SesionUnitOfWorkPostgres.
func NewSesionUnitOfWork(
	db *gorm.DB,
	sesionRepo sesiones_domain.SesionRepositorio,
	credencialesRepo seguridad_domain.CredencialesRepositorio,
	usuarioRepo usuario_domain.UsuarioRepositorio,
	encriptacionServicio seguridad_domain.EncriptacionServicio,
	tokenServicio sesiones_domain.TokenServicio,
	generadorID shared_domain.GeneradorID,
) sesiones_domain.UnitOfWork {
	return &SesionUnitOfWorkPostgres{
		db:                   db,
		sesionRepo:           sesionRepo,
		credencialesRepo:     credencialesRepo,
		usuarioRepo:          usuarioRepo,
		encriptacionServicio: encriptacionServicio,
		tokenServicio:        tokenServicio,
		generadorID:          generadorID,
	}
}

// Transaccional ejecuta fn dentro de una transacción GORM.
// Si fn retorna error → ROLLBACK automático.
// Si fn retorna nil → COMMIT automático.
func (uw *SesionUnitOfWorkPostgres) Transaccional(ctx context.Context, fn func(tx sesiones_domain.UnitOfWork) error) error {
	return uw.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		txUow := &SesionUnitOfWorkPostgres{
			db:                   txDB,
			sesionRepo:           NewSesionRepositorio(txDB),
			credencialesRepo:     uw.credencialesRepo,
			usuarioRepo:          uw.usuarioRepo,
			encriptacionServicio: uw.encriptacionServicio,
			tokenServicio:        uw.tokenServicio,
			generadorID:          uw.generadorID,
		}
		return fn(txUow)
	})
}

func (uw *SesionUnitOfWorkPostgres) SesionRepositorio() sesiones_domain.SesionRepositorio {
	return uw.sesionRepo
}
func (uw *SesionUnitOfWorkPostgres) CredencialesRepositorio() seguridad_domain.CredencialesRepositorio {
	return uw.credencialesRepo
}
func (uw *SesionUnitOfWorkPostgres) UsuarioRepositorio() usuario_domain.UsuarioRepositorio {
	return uw.usuarioRepo
}
func (uw *SesionUnitOfWorkPostgres) EncriptacionServicio() seguridad_domain.EncriptacionServicio {
	return uw.encriptacionServicio
}
func (uw *SesionUnitOfWorkPostgres) TokenServicio() sesiones_domain.TokenServicio {
	return uw.tokenServicio
}
func (uw *SesionUnitOfWorkPostgres) GeneradorID() shared_domain.GeneradorID { return uw.generadorID }

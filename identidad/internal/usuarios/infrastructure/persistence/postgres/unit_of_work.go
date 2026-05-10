package postgres

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	"gorm.io/gorm"
)

// UnitOfWorkPostgres implementa la interfaz usuario.UnitOfWork usando GORM
type UnitOfWorkPostgres struct {
	db                   *gorm.DB
	usuarioRepo          usuario.UsuarioRepositorio
	credencialesRepo     domain.CredencialesRepositorio
	encriptacionServicio domain.EncriptacionServicio
	generadorID          shared_domain.GeneradorID
}

// NewUnitOfWork crea una nueva instancia de UnitOfWorkPostgres
func NewUnitOfWork(
	db *gorm.DB,
	usuarioRepo usuario.UsuarioRepositorio,
	credencialesRepo domain.CredencialesRepositorio,
	encriptacionServicio domain.EncriptacionServicio,
	generadorID shared_domain.GeneradorID,
) usuario.UnitOfWork {
	return &UnitOfWorkPostgres{
		db:                   db,
		usuarioRepo:          usuarioRepo,
		credencialesRepo:     credencialesRepo,
		encriptacionServicio: encriptacionServicio,
		generadorID:          generadorID,
	}
}

// Transaccional envuelve la función en una transacción GORM
// Si fn retorna error → ROLLBACK automático
// Si fn retorna nil → COMMIT automático
func (uw *UnitOfWorkPostgres) Transaccional(ctx context.Context, fn func(tx usuario.UnitOfWork) error) error {
	return uw.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		// Crear una copia del UnitOfWork pero con la transacción de GORM
		txUow := &UnitOfWorkPostgres{
			db:                   txDB, // ← Usa la transacción, no la conexión normal
			usuarioRepo:          uw.usuarioRepo,
			credencialesRepo:     uw.credencialesRepo,
			encriptacionServicio: uw.encriptacionServicio,
			generadorID:          uw.generadorID,
		}

		// Ejecutar la función con la transacción
		return fn(txUow)
		// Si fn retorna error, GORM hace ROLLBACK automático
		// Si fn retorna nil, GORM hace COMMIT automático
	})
}

// UsuarioRepository retorna el repositorio de usuarios
func (uw *UnitOfWorkPostgres) UsuarioRepository() usuario.UsuarioRepositorio {
	return uw.usuarioRepo
}

// CredencialesRepository retorna el repositorio de credenciales
func (uw *UnitOfWorkPostgres) CredencialesRepository() domain.CredencialesRepositorio {
	return uw.credencialesRepo
}

// EncriptacionServicio retorna el servicio de encriptación
func (uw *UnitOfWorkPostgres) EncriptacionServicio() domain.EncriptacionServicio {
	return uw.encriptacionServicio
}

// GeneradorID retorna el generador de IDs
func (uw *UnitOfWorkPostgres) GeneradorID() shared_domain.GeneradorID {
	return uw.generadorID
}

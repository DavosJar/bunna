package postgres

import (
	"context"

	bootstrap "github.com/davosjar/bunna/services/identidad/internal/bootstrap/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	rbac_postgres "github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/persistence/postgres"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

// BootstrapUnitOfWorkPostgres implementa bootstrap.UnitOfWork usando GORM.
//
// Replica el patrón CORRECTO de `SesionUnitOfWorkPostgres`: en
// `Transaccional`, abre una transacción de GORM y construye un nuevo UoW
// cuyos repositorios operan sobre `txDB` (la conexión transaccional), de
// modo que todas las escrituras participan atómicamente y se revierten
// automáticamente si `fn` retorna error.
//
// `EncriptacionServicio` y `GeneradorID` se reutilizan dentro y fuera de
// la tx (no tocan la BD).
//
// Para el path de lectura (pre-check de idempotencia fuera de la tx), los
// getters exponen los repositorios configurados con la conexión plain-db
// pasada al constructor.
type BootstrapUnitOfWorkPostgres struct {
	db                   *gorm.DB
	usuarioRepo          usuario.UsuarioRepositorio
	credencialesRepo     seguridad.CredencialesRepositorio
	usuarioRolRepo       rbac.UsuarioRolRepositorio
	rolRepo              rbac.RolRepositorio
	encriptacionServicio seguridad.EncriptacionServicio
	generadorID          shareddomain.GeneradorID
}

// NewBootstrapUnitOfWork construye el UoW de bootstrap.
//
// Los repositorios pasados deben estar construidos sobre la misma `db`
// (conexión plain-db fuera de tx). Dentro de `Transaccional` se
// reconstruyen con `txDB`.
func NewBootstrapUnitOfWork(
	db *gorm.DB,
	usuarioRepo usuario.UsuarioRepositorio,
	credencialesRepo seguridad.CredencialesRepositorio,
	rolRepo rbac.RolRepositorio,
	usuarioRolRepo rbac.UsuarioRolRepositorio,
	encriptacionServicio seguridad.EncriptacionServicio,
	generadorID shareddomain.GeneradorID,
) bootstrap.UnitOfWork {
	return &BootstrapUnitOfWorkPostgres{
		db:                   db,
		usuarioRepo:          usuarioRepo,
		credencialesRepo:     credencialesRepo,
		usuarioRolRepo:       usuarioRolRepo,
		rolRepo:              rolRepo,
		encriptacionServicio: encriptacionServicio,
		generadorID:          generadorID,
	}
}

// Transaccional ejecuta fn dentro de una transacción GORM.
// Reconstruye los 4 repos con `txDB` para que sus escrituras participen
// de la tx (ver ADR-001 §4 — patrón `SesionUnitOfWorkPostgres`).
// Si fn retorna error → ROLLBACK automático.
// Si fn retorna nil   → COMMIT automático.
func (uw *BootstrapUnitOfWorkPostgres) Transaccional(
	ctx context.Context,
	fn func(tx bootstrap.UnitOfWork) error,
) error {
	return uw.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		txUow := &BootstrapUnitOfWorkPostgres{
			db:                   txDB,
			usuarioRepo:          usuarios_postgres.NewUsuarioRepositorio(txDB),
			credencialesRepo:     seguridad_postgres.NewCredencialesRepositorio(txDB),
			usuarioRolRepo:       rbac_postgres.NewUsuarioRolRepositorio(txDB),
			rolRepo:              rbac_postgres.NewRolRepositorio(txDB),
			encriptacionServicio: uw.encriptacionServicio,
			generadorID:          uw.generadorID,
		}
		return fn(txUow)
	})
}

// Getters — fuera de la tx devuelven los repos plain-db;
// dentro de la tx (sobre la instancia construida en `Transaccional`)
// devuelven los repos txDB.

func (uw *BootstrapUnitOfWorkPostgres) UsuarioRepositorio() usuario.UsuarioRepositorio {
	return uw.usuarioRepo
}

func (uw *BootstrapUnitOfWorkPostgres) CredencialesRepositorio() seguridad.CredencialesRepositorio {
	return uw.credencialesRepo
}

func (uw *BootstrapUnitOfWorkPostgres) UsuarioRolRepositorio() rbac.UsuarioRolRepositorio {
	return uw.usuarioRolRepo
}

func (uw *BootstrapUnitOfWorkPostgres) RolRepositorio() rbac.RolRepositorio {
	return uw.rolRepo
}

func (uw *BootstrapUnitOfWorkPostgres) EncriptacionServicio() seguridad.EncriptacionServicio {
	return uw.encriptacionServicio
}

func (uw *BootstrapUnitOfWorkPostgres) GeneradorID() shareddomain.GeneradorID {
	return uw.generadorID
}

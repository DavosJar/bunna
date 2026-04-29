package registry

import (
	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/security/bcrypt"
	shared_idgenerator "github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/idgenerator"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

type Registry struct {
	usuarioRepository      usuario.UsuarioRepositorio
	credencialesRepository domain.CredencialesRepositorio
	encriptacionServicio   domain.EncriptacionServicio
	usuarioUnitOfWork      usuario.UnitOfWork
}

func NewRegistry(db *gorm.DB, cfg *config.Config) *Registry {
	// Crear repositorios
	usuarioRepo := usuarios_postgres.NewUsuarioRepositorio(db)
	credencialesRepo := seguridad_postgres.NewCredencialesRepositorio(db)
	encriptacion := bcrypt.NewBcryptEncriptacion(cfg.BcryptCost)

	// Crear generador de IDs (transversal desde shared)
	generadorID := shared_idgenerator.NewUUIDv7Generator()

	// Crear UnitOfWork con los repositorios y servicios
	unitOfWork := usuarios_postgres.NewUnitOfWork(
		db,
		usuarioRepo,
		credencialesRepo,
		encriptacion,
		generadorID,
	)

	return &Registry{
		usuarioRepository:      usuarioRepo,
		credencialesRepository: credencialesRepo,
		encriptacionServicio:   encriptacion,
		usuarioUnitOfWork:      unitOfWork,
	}
}

func (r *Registry) UsuarioRepository() usuario.UsuarioRepositorio {
	return r.usuarioRepository
}

func (r *Registry) CredencialesRepository() domain.CredencialesRepositorio {
	return r.credencialesRepository
}

func (r *Registry) EncriptacionServicio() domain.EncriptacionServicio {
	return r.encriptacionServicio
}

func (r *Registry) UsuarioUnitOfWork() usuario.UnitOfWork {
	return r.usuarioUnitOfWork
}

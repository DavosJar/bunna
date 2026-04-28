package registry

import (
	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

type Registry struct {
	usuarioRepository      usuario.UsuarioRepositorio
	credencialesRepository domain.CredencialesRepositorio
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		usuarioRepository:      postgres.NewUsuarioRepositorio(db),
		credencialesRepository: seguridad_postgres.NewCredencialesRepositorio(db),
	}
}

func (r *Registry) UsuarioRepository() usuario.UsuarioRepositorio {
	return r.usuarioRepository
}

func (r *Registry) CredencialesRepository() domain.CredencialesRepositorio {
	return r.credencialesRepository
}

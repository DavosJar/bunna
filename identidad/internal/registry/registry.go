package registry

import (
	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

type Registry struct {
	usuarioRepository usuario.UsuarioRepositorio
}

func NewRegistry(db *gorm.DB) *Registry {
	return &Registry{
		usuarioRepository: postgres.NewUsuarioRepositorio(db),
	}
}

func (r *Registry) UsuarioRepository() usuario.UsuarioRepositorio {
	return r.usuarioRepository
}

package domain

import (
	"context"

	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type UnitOfWork interface {
	Transaccional(ctx context.Context, fn func(tx UnitOfWork) error) error

	SesionRepositorio() SesionRepositorio
	CredencialesRepositorio() seguridad_domain.CredencialesRepositorio
	UsuarioRepositorio() usuario_domain.UsuarioRepositorio
	EncriptacionServicio() seguridad_domain.EncriptacionServicio
	TokenServicio() TokenServicio
	GeneradorID() shared_domain.GeneradorID
}

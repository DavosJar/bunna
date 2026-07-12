package register

import (
	"context"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type UnitOfWork interface {
	Transaccional(ctx context.Context, fn func(tx UnitOfWork) error) error
	UsuarioRepository() usuario.UsuarioRepositorio
	CredencialesRepository() seguridad.CredencialesRepositorio
	TenantRepository() tenant.TenantRepositorio
	MembresiaRepository() tenant.MembresiaRepositorio
	RolRepository() rbac.RolRepositorio
	UsuarioTenantRolRepository() rbac.UsuarioTenantRolRepositorio
	RolPermisoRepository() rbac.RolPermisoRepositorio
	PermisoRepository() rbac.PermisoRepositorio
	RolPublisher() rbac.RolPublisher
}

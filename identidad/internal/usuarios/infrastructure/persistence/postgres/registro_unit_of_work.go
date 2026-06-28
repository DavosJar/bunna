package postgres

import (
	"context"

	rbac_domain "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	rbac_postgres "github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/persistence/postgres"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	tenant_domain "github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	tenant_postgres "github.com/davosjar/bunna/services/identidad/internal/tenants/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/register"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	"gorm.io/gorm"
)

type RegistroUnitOfWorkPostgres struct {
	db                   *gorm.DB
	usuarioRepo          usuario_domain.UsuarioRepositorio
	credencialesRepo     seguridad_domain.CredencialesRepositorio
	tenantRepo           tenant_domain.TenantRepositorio
	membresiaRepo        tenant_domain.MembresiaRepositorio
	rolRepo              rbac_domain.RolRepositorio
	usuarioTenantRolRepo rbac_domain.UsuarioTenantRolRepositorio
	rolPermisoRepo       rbac_domain.RolPermisoRepositorio
}

func NewRegistroUnitOfWork(db *gorm.DB) register.UnitOfWork {
	return &RegistroUnitOfWorkPostgres{
		db:                   db,
		usuarioRepo:          NewUsuarioRepositorio(db),
		credencialesRepo:     seguridad_postgres.NewCredencialesRepositorio(db),
		tenantRepo:           tenant_postgres.NewTenantRepositorio(db),
		membresiaRepo:        tenant_postgres.NewMembresiaRepositorio(db),
		rolRepo:              rbac_postgres.NewRolRepositorio(db),
		usuarioTenantRolRepo: rbac_postgres.NewUsuarioTenantRolRepositorio(db),
		rolPermisoRepo:       rbac_postgres.NewRolPermisoRepositorio(db),
	}
}

func (u *RegistroUnitOfWorkPostgres) Transaccional(ctx context.Context, fn func(tx register.UnitOfWork) error) error {
	return u.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		txUow := &RegistroUnitOfWorkPostgres{
			db:                   txDB,
			usuarioRepo:          NewUsuarioRepositorio(txDB),
			credencialesRepo:     seguridad_postgres.NewCredencialesRepositorio(txDB),
			tenantRepo:           tenant_postgres.NewTenantRepositorio(txDB),
			membresiaRepo:        tenant_postgres.NewMembresiaRepositorio(txDB),
			rolRepo:              rbac_postgres.NewRolRepositorio(txDB),
			usuarioTenantRolRepo: rbac_postgres.NewUsuarioTenantRolRepositorio(txDB),
			rolPermisoRepo:       rbac_postgres.NewRolPermisoRepositorio(txDB),
		}
		return fn(txUow)
	})
}

func (u *RegistroUnitOfWorkPostgres) UsuarioRepository() usuario_domain.UsuarioRepositorio {
	return u.usuarioRepo
}

func (u *RegistroUnitOfWorkPostgres) CredencialesRepository() seguridad_domain.CredencialesRepositorio {
	return u.credencialesRepo
}

func (u *RegistroUnitOfWorkPostgres) TenantRepository() tenant_domain.TenantRepositorio {
	return u.tenantRepo
}

func (u *RegistroUnitOfWorkPostgres) MembresiaRepository() tenant_domain.MembresiaRepositorio {
	return u.membresiaRepo
}

func (u *RegistroUnitOfWorkPostgres) RolRepository() rbac_domain.RolRepositorio {
	return u.rolRepo
}

func (u *RegistroUnitOfWorkPostgres) UsuarioTenantRolRepository() rbac_domain.UsuarioTenantRolRepositorio {
	return u.usuarioTenantRolRepo
}

func (u *RegistroUnitOfWorkPostgres) RolPermisoRepository() rbac_domain.RolPermisoRepositorio {
	return u.rolPermisoRepo
}

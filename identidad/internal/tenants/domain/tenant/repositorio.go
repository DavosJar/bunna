package tenant

import "context"

// TenantRepositorio define las operaciones de persistencia para tenants.
type TenantRepositorio interface {
	Crear(ctx context.Context, t *Tenant) (*Tenant, error)
	ObtenerPorID(ctx context.Context, id string) (*Tenant, error)
	ObtenerPorSlug(ctx context.Context, slug string) (*Tenant, error)
	Actualizar(ctx context.Context, t *Tenant) (*Tenant, error)
	Listar(ctx context.Context) ([]*Tenant, error)
	ListarPorUsuario(ctx context.Context, usuarioID string) ([]*Tenant, error)
}

// MembresiaRepositorio define las operaciones de persistencia para membresías.
type MembresiaRepositorio interface {
	Crear(ctx context.Context, m *Membresia) error
	Eliminar(ctx context.Context, usuarioID, tenantID string) error
	ExisteMiembro(ctx context.Context, usuarioID, tenantID string) (bool, error)
	ListarUsuariosPorTenant(ctx context.Context, tenantID string) ([]string, error)
}

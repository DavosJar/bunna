package rbac

import (
	"context"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

// RolRepositorio define las operaciones de persistencia para roles
type RolRepositorio interface {
	ObtenerPorNombre(ctx context.Context, nombre string) (*RolDB, error)
	ObtenerPorID(ctx context.Context, id string) (*RolDB, error)
	Listar(ctx context.Context, especificacion EspecificacionRol, paginacion shareddomain.Paginacion) ([]*RolDB, error)
	Crear(ctx context.Context, rol *RolDB) error
	ActualizarDescripcion(ctx context.Context, id, descripcion string) error
}

// PermisoRepositorio define las operaciones de persistencia para permisos
type PermisoRepositorio interface {
	ObtenerPorCodigo(ctx context.Context, codigo string) (*PermisoDB, error)
	Listar(ctx context.Context) ([]*PermisoDB, error)
	Crear(ctx context.Context, permiso *PermisoDB) error
	ActualizarNombreDescripcion(ctx context.Context, id, nombre, descripcion string) error
	ListarPorRol(ctx context.Context, rolID string) ([]*PermisoDB, error)
}

// RolPermisoRepositorio maneja la relación rol ↔ permiso
type RolPermisoRepositorio interface {
	LimpiarPermisosDeRol(ctx context.Context, rolID string) error
	AsignarPermiso(ctx context.Context, rolID, permisoID string) error
}

// UsuarioRolRepositorio maneja roles globales (solo SYS_ADMIN)
type UsuarioRolRepositorio interface {
	Crear(ctx context.Context, usuarioID, rolID string) error
	Eliminar(ctx context.Context, usuarioID, rolID string) error
	ListarRolesPorUsuario(ctx context.Context, usuarioID string) ([]*RolDB, error)
	TieneRol(ctx context.Context, usuarioID, rolNombre string) (bool, error)
}

// UsuarioTenantRolRepositorio maneja roles de usuario dentro de un tenant
type UsuarioTenantRolRepositorio interface {
	Crear(ctx context.Context, usuarioID, tenantID, rolID string) error
	Eliminar(ctx context.Context, usuarioID, tenantID, rolID string) error
	ListarRolesPorUsuarioEnTenant(ctx context.Context, usuarioID, tenantID string) ([]*RolDB, error)
	TieneRolEnTenant(ctx context.Context, usuarioID, tenantID, rolNombre string) (bool, error)
}

// RolDB es la representación de un rol en la capa de dominio/repositorio
type RolDB struct {
	ID          string
	Nombre      string
	Descripcion string
	EsSistema   bool
}

// PermisoDB es la representación de un permiso en la capa de dominio/repositorio
type PermisoDB struct {
	ID          string
	Codigo      string
	Nombre      string
	Descripcion string
	Modulo      string
}
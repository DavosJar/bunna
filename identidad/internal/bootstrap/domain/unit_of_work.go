package bootstrap

import (
	"context"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

// UnitOfWork es la fachada transaccional para el caso de uso de bootstrap
// del primer sys_admin. Cruza 3 BCs (usuarios, seguridad, rbac) y por eso
// vive en `internal/bootstrap/domain` (preocupación operacional cross-cutting).
//
// Las implementaciones deben abrir una transacción de BD en `Transaccional`
// y reconstruir los repositorios con la conexión transaccional (txDB) antes
// de invocar `fn`, de modo que todas las escrituras participen atómicamente.
// El patrón correcto se replica de `SesionUnitOfWorkPostgres`.
//
// Para el path de lectura (pre-check de idempotencia fuera de la tx), los
// getters exponen los repositorios configurados con la conexión plain-db.
type UnitOfWork interface {
	// Transaccional ejecuta fn dentro de una transacción de BD.
	//.fn recibe un `tx UnitOfWork` cuyos repos operan sobre la tx abierta.
	// Si fn retorna error → ROLLBACK automático.
	// Si fn retorna nil   → COMMIT automático.
	Transaccional(ctx context.Context, fn func(tx UnitOfWork) error) error

	// Repositorios (sobre la conexión del UoW: plain-db fuera de tx, txDB dentro).
	UsuarioRepositorio() usuario.UsuarioRepositorio
	CredencialesRepositorio() seguridad.CredencialesRepositorio
	UsuarioRolRepositorio() rbac.UsuarioRolRepositorio
	RolRepositorio() rbac.RolRepositorio

	// Servicios (no tocan la BD; se reutilizan dentro y fuera de tx).
	EncriptacionServicio() seguridad.EncriptacionServicio
	GeneradorID() shareddomain.GeneradorID
}

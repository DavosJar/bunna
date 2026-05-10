package usuario

import (
	"context"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

// UnitOfWork es la fachada para operaciones transaccionales
// Oculta los detalles de GORM y proporciona una interfaz limpia para manejar transacciones atómicas
type UnitOfWork interface {
	// Transaccional ejecuta la función dentro de una transacción
	// Si fn retorna error, hace rollback automático
	// Si fn retorna nil, hace commit automático
	Transaccional(ctx context.Context, fn func(tx UnitOfWork) error) error

	// UsuarioRepository retorna el repositorio de usuarios dentro de la transacción
	UsuarioRepository() UsuarioRepositorio

	// CredencialesRepository retorna el repositorio de credenciales dentro de la transacción
	CredencialesRepository() domain.CredencialesRepositorio

	// EncriptacionServicio retorna el servicio de encriptación
	EncriptacionServicio() domain.EncriptacionServicio

	// GeneradorID retorna el generador de IDs (desde shared)
	GeneradorID() shared_domain.GeneradorID
}

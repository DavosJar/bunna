package recuperacion

import "context"

// TokenRecuperacionRepositorio define operaciones de persistencia
type TokenRecuperacionRepositorio interface {
	Crear(ctx context.Context, token *TokenRecuperacion) error
	ObtenerPorHash(ctx context.Context, hash string) (*TokenRecuperacion, error)
	Actualizar(ctx context.Context, token *TokenRecuperacion) error
}

// UsuarioRecuperacion proyección del usuario para recuperación
type UsuarioRecuperacion struct {
	ID      string
	Nombre  string
	Correo  string
}

// UsuarioRecuperacionRepositorio operaciones sobre usuarios para recuperación
type UsuarioRecuperacionRepositorio interface {
	ObtenerPorCorreo(ctx context.Context, correo string) (*UsuarioRecuperacion, error)
	ActualizarPassword(ctx context.Context, usuarioID, nuevoHash string) error
}

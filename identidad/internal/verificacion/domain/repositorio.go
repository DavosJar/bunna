package verificacion

import "context"

// VerificacionRepositorio define operaciones para buscar usuarios por token
type VerificacionRepositorio interface {
	ObtenerPorHashToken(ctx context.Context, hash string) (*UsuarioVerificacion, error)
	ActualizarPrueba(ctx context.Context, usuarioID string, prueba PruebaVerificacion) error
	MarcarVerificado(ctx context.Context, usuarioID string) error
	ObtenerPorID(ctx context.Context, usuarioID string) (*UsuarioVerificacion, error)
}

// UsuarioVerificacion es una proyección del usuario con solo los datos necesarios
type UsuarioVerificacion struct {
	ID                 string
	Nombre             string
	Correo             string
	EstadoVerificacion string
	PruebaVerificacion PruebaVerificacion
	ContadorReenvios   int
	UltimoReenvio      interface{}
}

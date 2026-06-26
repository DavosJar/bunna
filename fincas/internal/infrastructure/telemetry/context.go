package telemetry

import "context"

type usuarioCtxKey struct{}

// WithUsuarioID almacena el ID de usuario autenticado en el contexto de la request.
func WithUsuarioID(ctx context.Context, usuarioID string) context.Context {
	return context.WithValue(ctx, usuarioCtxKey{}, usuarioID)
}

// GetUsuarioIDFromCtx extrae el ID de usuario del contexto. Retorna vacío si no existe.
func GetUsuarioIDFromCtx(ctx context.Context) string {
	if v := ctx.Value(usuarioCtxKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

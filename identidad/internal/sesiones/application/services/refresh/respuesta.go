package refresh

import "time"

// RespuestaRefresh contiene el nuevo par de tokens generado tras un refresh exitoso.
type RespuestaRefresh struct {
	// AccessToken es el nuevo token de acceso en texto plano.
	AccessToken string

	// RefreshToken es el nuevo token de refresco en texto plano.
	// El token anterior queda invalidado tras este refresh.
	RefreshToken string

	// ExpiracionAccess es la fecha y hora de expiración del nuevo access token.
	ExpiracionAccess time.Time

	// ExpiracionRefresh es la fecha y hora de expiración del nuevo refresh token.
	ExpiracionRefresh time.Time

	// SesionID es el identificador de la sesión renovada.
	SesionID string

	// UsuarioID es el identificador del usuario dueño de la sesión.
	UsuarioID string
}

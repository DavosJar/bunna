// Package refresh implementa el caso de uso de renovación de sesión mediante refresh token.
package refresh

// ComandoRefresh contiene los datos necesarios para renovar una sesión activa.
type ComandoRefresh struct {
	// RefreshToken es el token de refresco emitido en el login o en un refresh anterior.
	RefreshToken string
}
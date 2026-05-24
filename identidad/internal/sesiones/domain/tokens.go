package domain

import "time"

// TokenPair es un value object inmutable.
// Agrupa access token + refresh token + sus fechas de expiración.
// Los tokens aquí son en PLANO (para retornar al cliente).
// Lo que se persiste en Sesion son los hashes, no este struct.
type TokenPair struct {
	accessToken       string
	refreshToken      string
	expiracionAccess  time.Time
	expiracionRefresh time.Time
}

func NuevoTokenPair(
	accessToken string,
	refreshToken string,
	expiracionAccess time.Time,
	expiracionRefresh time.Time,
) (TokenPair, error) {
	if accessToken == "" {
		return TokenPair{}, ErrAccessTokenRequerido
	}
	if refreshToken == "" {
		return TokenPair{}, ErrRefreshTokenRequerido
	}
	return TokenPair{
		accessToken:       accessToken,
		refreshToken:      refreshToken,
		expiracionAccess:  expiracionAccess,
		expiracionRefresh: expiracionRefresh,
	}, nil
}

func (t TokenPair) AccessToken() string          { return t.accessToken }
func (t TokenPair) RefreshToken() string         { return t.refreshToken }
func (t TokenPair) ExpiracionAccess() time.Time  { return t.expiracionAccess }
func (t TokenPair) ExpiracionRefresh() time.Time { return t.expiracionRefresh }

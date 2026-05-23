package recuperacion

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// TokenRecuperacion representa el token de restablecimiento de contraseña
type TokenRecuperacion struct {
	id        string
	usuarioID string
	tokenHash string
	expiraEn  time.Time
	usado     bool
	creadoEn  time.Time
	usadoEn   *time.Time
}

// NuevoTokenRecuperacion crea un token con el hash del token en plano
func NuevoTokenRecuperacion(id, usuarioID, tokenEnPlano string, expiraEn time.Time) *TokenRecuperacion {
	return &TokenRecuperacion{
		id:        id,
		usuarioID: usuarioID,
		tokenHash: HashearToken(tokenEnPlano),
		expiraEn:  expiraEn,
		usado:     false,
		creadoEn:  time.Now(),
	}
}

// NuevoTokenRecuperacionDesdeBD reconstruye desde persistencia
func NuevoTokenRecuperacionDesdeBD(id, usuarioID, tokenHash string, expiraEn time.Time, usado bool, creadoEn time.Time, usadoEn *time.Time) *TokenRecuperacion {
	return &TokenRecuperacion{
		id:        id,
		usuarioID: usuarioID,
		tokenHash: tokenHash,
		expiraEn:  expiraEn,
		usado:     usado,
		creadoEn:  creadoEn,
		usadoEn:   usadoEn,
	}
}

// EsValido verifica que el token sea válido (no expirado, no usado)
func (t *TokenRecuperacion) EsValido(ahora time.Time) error {
	if t.usado {
		return ErrEnlaceYaUtilizado
	}
	if ahora.After(t.expiraEn) {
		return ErrEnlaceExpirado
	}
	return nil
}

// Usar marca el token como usado
func (t *TokenRecuperacion) Usar(ahora time.Time) {
	t.usado = true
	t.usadoEn = &ahora
}

// HashearToken genera hash SHA-256 de un token
func HashearToken(tokenEnPlano string) string {
	hash := sha256.Sum256([]byte(tokenEnPlano))
	return fmt.Sprintf("%x", hash)
}

// Getters
func (t *TokenRecuperacion) ID() string        { return t.id }
func (t *TokenRecuperacion) UsuarioID() string { return t.usuarioID }
func (t *TokenRecuperacion) TokenHash() string { return t.tokenHash }
func (t *TokenRecuperacion) ExpiraEn() time.Time { return t.expiraEn }
func (t *TokenRecuperacion) Usado() bool       { return t.usado }
func (t *TokenRecuperacion) CreadoEn() time.Time { return t.creadoEn }
func (t *TokenRecuperacion) UsadoEn() *time.Time { return t.usadoEn }

package verificacion

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// PruebaVerificacion es un Value Object que encapsula el hash del token
// y su fecha de expiración. Nunca almacena el token en plano.
type PruebaVerificacion struct {
	secretoHash string
	expiraEn    time.Time
}

// NuevaPruebaVerificacion crea una prueba con el hash del token dado.
func NuevaPruebaVerificacion(tokenEnPlano string, expiraEn time.Time) PruebaVerificacion {
	hash := hashearToken(tokenEnPlano)
	return PruebaVerificacion{
		secretoHash: hash,
		expiraEn:    expiraEn,
	}
}

// NuevaPruebaVerificacionDesdeBD reconstruye desde persistencia con hash ya calculado.
func NuevaPruebaVerificacionDesdeBD(secretoHash string, expiraEn time.Time) PruebaVerificacion {
	return PruebaVerificacion{
		secretoHash: secretoHash,
		expiraEn:    expiraEn,
	}
}

// PruebaVerificacionVacia retorna una prueba vacía (sin token asignado).
func PruebaVerificacionVacia() PruebaVerificacion {
	return PruebaVerificacion{}
}

// Expiro retorna true si el token ya expiró.
func (p PruebaVerificacion) Expiro(ahora time.Time) bool {
	if p.expiraEn.IsZero() {
		return true
	}
	return ahora.After(p.expiraEn)
}

// EstaPendiente retorna true si tiene un secreto asignado.
func (p PruebaVerificacion) EstaPendiente() bool {
	return p.secretoHash != ""
}

// CoincideCon verifica si el token en plano coincide con el hash almacenado.
func (p PruebaVerificacion) CoincideCon(tokenEnPlano string) bool {
	return p.secretoHash == hashearToken(tokenEnPlano)
}

// Getters
func (p PruebaVerificacion) SecretoHash() string { return p.secretoHash }
func (p PruebaVerificacion) ExpiraEn() time.Time  { return p.expiraEn }

// HashearTokenPublico permite hashear un token desde fuera del dominio.
func HashearToken(tokenEnPlano string) string {
	return hashearToken(tokenEnPlano)
}

func hashearToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptEncriptacion implementa la interfaz EncriptacionServicio usando bcrypt
type BcryptEncriptacion struct {
	cost int
}

// NewBcryptEncriptacion crea una nueva instancia de BcryptEncriptacion
// cost es el factor de costo de bcrypt (recomendado: 12, rango: 10-14)
func NewBcryptEncriptacion(cost int) *BcryptEncriptacion {
	return &BcryptEncriptacion{cost: cost}
}

// Hashear toma un password en plano y retorna su hash usando bcrypt
func (b *BcryptEncriptacion) Hashear(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(hash), err
}

// Verificar compara un password en plano contra su hash usando bcrypt
func (b *BcryptEncriptacion) Verificar(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

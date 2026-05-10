package domain

// EncriptacionServicio define la interfaz para operaciones de encriptación y verificación de passwords
type EncriptacionServicio interface {
	// Hashear toma un password en plano y retorna su hash
	Hashear(password string) (string, error)

	// Verificar compara un password en plano contra un hash y retorna true si coinciden
	Verificar(password, hash string) bool
}
